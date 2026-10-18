package contracts

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrSignatureExpired = errors.New("url signature has expired")
	ErrInvalidSignature = errors.New("invalid url signature")
	ErrMalformedURL     = errors.New("malformed signed url or token")
)

// SignChunkURL generates an HMAC-SHA256 signature binding node, media, chunk, and expiration time.
func SignChunkURL(secret []byte, nodeID, mediaID, chunkPath string, expiresAt time.Time) string {
	exp := strconv.FormatInt(expiresAt.Unix(), 10)
	payload := fmt.Sprintf("%s:%s:%s:%s", nodeID, mediaID, chunkPath, exp)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("sig=%s&exp=%s", sig, exp)
}

// VerifyChunkSignature verifies an HMAC-SHA256 signature against node, media, chunk, and expiration.
func VerifyChunkSignature(secret []byte, nodeID, mediaID, chunkPath string, sigParam, expParam string, now time.Time) error {
	expUnix, err := strconv.ParseInt(expParam, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: invalid exp value", ErrMalformedURL)
	}

	expTime := time.Unix(expUnix, 0)
	if now.After(expTime) {
		return ErrSignatureExpired
	}

	payload := fmt.Sprintf("%s:%s:%s:%s", nodeID, mediaID, chunkPath, expParam)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(strings.ToLower(sigParam)), []byte(strings.ToLower(expectedSig))) {
		return ErrInvalidSignature
	}

	return nil
}
