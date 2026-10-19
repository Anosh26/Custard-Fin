package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
)

// S3API abstracts the AWS S3 client methods used by S3ChunkCache for testability.
type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// Config defines connection and bucket settings for S3-compatible chunk storage.
type Config struct {
	Endpoint        string // Custom endpoint for MinIO/Garage/Floci (empty for standard AWS S3)
	Bucket          string // S3 bucket name
	Region          string // AWS region (default: "us-east-1")
	AccessKeyID     string // S3 Access Key ID
	SecretAccessKey string // S3 Secret Access Key
	UsePathStyle    bool   // Force path-style addressing (required for MinIO/Garage/Floci)
}

// Validate verifies required configuration.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Bucket) == "" {
		return errors.New("s3 bucket name is required")
	}
	if strings.TrimSpace(c.Region) == "" {
		c.Region = "us-east-1"
	}
	return nil
}

// S3ChunkCache implements ports.ChunkCache using an S3-compatible object store.
// Hot chunks are uploaded once from the edge (Invariant 7) and tagged for 24-hour retention (Invariant 1).
type S3ChunkCache struct {
	client S3API
	bucket string
}

// NewS3ChunkCache initializes an S3ChunkCache from configuration.
func NewS3ChunkCache(cfg Config) (*S3ChunkCache, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid s3 config: %w", err)
	}

	s3Options := func(o *s3.Options) {
		o.Region = cfg.Region
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
		if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
			o.Credentials = credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")
		}
	}

	client := s3.New(s3.Options{}, s3Options)
	return NewS3ChunkCacheWithClient(client, cfg.Bucket), nil
}

// NewS3ChunkCacheWithClient creates an S3ChunkCache using an injected S3API client.
func NewS3ChunkCacheWithClient(client S3API, bucket string) *S3ChunkCache {
	return &S3ChunkCache{
		client: client,
		bucket: bucket,
	}
}

// Get retrieves an object from the S3 cache.
func (s *S3ChunkCache) Get(ctx context.Context, key string) (*ports.CacheItem, error) {
	if strings.TrimSpace(key) == "" {
		return nil, errors.New("empty chunk cache key")
	}

	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *types.NotFound
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &notFound) || errors.As(err, &noSuchKey) || strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("s3 get failed for key %q: %w", key, err)
	}

	contentType := "application/octet-stream"
	if out.ContentType != nil && *out.ContentType != "" {
		contentType = *out.ContentType
	}

	var length int64
	if out.ContentLength != nil {
		length = *out.ContentLength
	}

	return &ports.CacheItem{
		Body:          out.Body,
		ContentType:   contentType,
		ContentLength: length,
	}, nil
}

// Put stores an object in the S3 cache, setting retention tagging (Invariant 1).
func (s *S3ChunkCache) Put(ctx context.Context, key string, r io.Reader, contentLength int64, contentType string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("empty chunk cache key")
	}
	if r == nil {
		return errors.New("cannot put nil reader into s3 cache")
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
		Tagging:     aws.String("retention=24h&managed-by=custardfin"),
	}
	if contentLength > 0 {
		input.ContentLength = aws.Int64(contentLength)
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("s3 put failed for key %q: %w", key, err)
	}
	return nil
}

// Exists checks if an object exists in the S3 bucket using HeadObject.
func (s *S3ChunkCache) Exists(ctx context.Context, key string) (bool, error) {
	if strings.TrimSpace(key) == "" {
		return false, errors.New("empty chunk cache key")
	}

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *types.NotFound
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &notFound) || errors.As(err, &noSuchKey) || strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, fmt.Errorf("s3 head failed for key %q: %w", key, err)
	}
	return true, nil
}

// Delete removes an object from the S3 bucket.
func (s *S3ChunkCache) Delete(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("empty chunk cache key")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 delete failed for key %q: %w", key, err)
	}
	return nil
}
