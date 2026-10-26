package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NodeRepository struct {
	pool *pgxpool.Pool
}

func NewNodeRepository(pool *pgxpool.Pool) ports.NodeRepository {
	return &NodeRepository{pool: pool}
}

func (r *NodeRepository) Create(ctx context.Context, node *domain.Node) error {
	query := `
	INSERT INTO nodes (id, user_id, name, node_token_hash, frp_auth_token, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query,
		node.ID,
		node.UserID,
		node.Name,
		node.NodeTokenHash,
		node.FrpAuthToken,
		node.CreatedAt,
		node.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert node: %w", err)
	}
	return nil
}

func (r *NodeRepository) FindByID(ctx context.Context, id string) (*domain.Node, error) {
	id = strings.TrimSpace(id)
	query := `
	SELECT id, user_id, name, node_token_hash, frp_auth_token, created_at, updated_at
	FROM nodes
	WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var n domain.Node
	err := row.Scan(&n.ID, &n.UserID, &n.Name, &n.NodeTokenHash, &n.FrpAuthToken, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: node not found", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to query node by id: %w", err)
	}
	return &n, nil
}

func (r *NodeRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Node, error) {
	tokenHash = strings.TrimSpace(tokenHash)
	query := `
	SELECT id, user_id, name, node_token_hash, frp_auth_token, created_at, updated_at
	FROM nodes
	WHERE node_token_hash = $1
	`
	row := r.pool.QueryRow(ctx, query, tokenHash)

	var n domain.Node
	err := row.Scan(&n.ID, &n.UserID, &n.Name, &n.NodeTokenHash, &n.FrpAuthToken, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: node not found", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to query node by token hash: %w", err)
	}
	return &n, nil
}

func (r *NodeRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Node, error) {
	userID = strings.TrimSpace(userID)
	query := `
	SELECT id, user_id, name, node_token_hash, frp_auth_token, created_at, updated_at
	FROM nodes
	WHERE user_id = $1
	ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes by user: %w", err)
	}
	defer rows.Close()

	var nodes []*domain.Node
	for rows.Next() {
		var n domain.Node
		if err := rows.Scan(&n.ID, &n.UserID, &n.Name, &n.NodeTokenHash, &n.FrpAuthToken, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, &n)
	}
	return nodes, rows.Err()
}

func (r *NodeRepository) UpdateOwner(ctx context.Context, nodeID, userID string) error {
	nodeID = strings.TrimSpace(nodeID)
	userID = strings.TrimSpace(userID)
	query := `
	UPDATE nodes
	SET user_id = $1, updated_at = NOW()
	WHERE id = $2
	`
	cmdTag, err := r.pool.Exec(ctx, query, userID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to update node owner: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *NodeRepository) UpdateName(ctx context.Context, nodeID, name string) error {
	nodeID = strings.TrimSpace(nodeID)
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.ErrValidation
	}
	query := `
	UPDATE nodes
	SET name = $1, updated_at = NOW()
	WHERE id = $2
	`
	cmdTag, err := r.pool.Exec(ctx, query, name, nodeID)
	if err != nil {
		return fmt.Errorf("failed to update node name: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

