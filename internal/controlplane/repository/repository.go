package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

type DatabaseInterface interface {
	InitTables(ctx context.Context) error
	InsertNode(ctx context.Context, node Node) (*Node, error)
}

func NewRepository(db *sql.DB) DatabaseInterface {
	return &Repository{db: db}
}

func (r *Repository) InitTables(ctx context.Context) error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS cluster (
			id TEXT PRIMARY KEY,
			pod_cidr TEXT,
			access_token TEXT,
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS nodes (
			id TEXT PRIMARY KEY,
			cluster_id TEXT FOREIGN KEY (cluster_id) REFERENCES cluster (id),
			name TEXT NOT NULL,
			ip TEXT NOT NULL,
			pod_cidr TEXT,
			status TEXT NOT NULL,
			last_heartbeat DATETIME NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS podes (
			id TEXT PRIMARY KEY,
			node_id TEXT NOT NULL FOREIGN KEY (node_id) REFERENCES nodes (id),
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

// pod_cidr - 10.100.1.0/24 -> 10.100.1.0 - 10.100.1.255, где 10.100.1 - подсеть ноды, а 10.100.0.0/16 -> подсеть всего кластера
func (r *Repository) InsertNode(ctx context.Context, node Node) (*Node, error) {
	_, err := r.db.Exec(`INSTERT INTO nodes (id, name, ip, pod_cidr, status, last_heartbeat) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		node.Name,
		node.Ip,
		node.PodCIDR,
		node.Status,
		node.LastHeartbeat,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert node: %w", err)
	}

	return &node, nil
}

//func (r *Repository) InsertPod(pod Pod)
