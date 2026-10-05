package agent

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

type AgentInterface interface {
	InitTables(ctx context.Context) error
	InsertNodeData(ctx context.Context, id string, podCidr string) error
}

func NewRepository(db *sql.DB) AgentInterface {
	return &Repository{db: db}
}

func (r *Repository) InitTables(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS node (
			id TEXT PRIMARY KEY,
			pod_cidr TEXT,
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS pods (
			id TEXT PRIMARY KEY,
			node_id TEXT NOT NULL FOREIGN KEY (node_id) REFERENCES node (id),
		    service TEXT NOT NULL,
		    tag TEXT NOT NULL,
		    is_healthy BOOLEAN NOT NULL,
		    ip  TEXT NOT NULL,
		    port NUMERIC NOT NULL,
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func (r *Repository) InsertNodeData(ctx context.Context, id string, podCidr string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO node (id, pod_cidr) VALUES (?, ?)`, id, podCidr)
	if err != nil {
		return fmt.Errorf("failed to insert node data: %w", err)
	}
	return nil
}
