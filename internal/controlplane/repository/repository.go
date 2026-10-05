package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

type ClusterInterface interface {
	InitTables(ctx context.Context) error
	InsertNode(ctx context.Context, node Node) (*Node, error)
	GetClusterInfo(ctx context.Context) (*Cluster, error)
	GetNodes(ctx context.Context) ([]Node, error)
}

func NewRepository(db *sql.DB) ClusterInterface {
	return &Repository{db: db}
}

func (r *Repository) InitTables(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS cluster (
			id TEXT PRIMARY KEY,
			pod_cidr TEXT,
			access_token TEXT,
			registry_username  TEXT,
			registry_password TEXT,
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS nodes (
			id TEXT PRIMARY KEY,
			cluster_id TEXT FOREIGN KEY (cluster_id) REFERENCES cluster (id),
			ip TEXT NOT NULL,
			pod_cidr TEXT,
			status TEXT NOT NULL,
			last_heartbeat DATETIME NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS pods (
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
	_, err := r.db.ExecContext(ctx, `INSERT INTO nodes (id, ip, pod_cidr, status, last_heartbeat) VALUES (?, ?, ?, ?, ?, ?, ?)`,
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

func (r *Repository) GetNodes(ctx context.Context) ([]Node, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM nodes`)
	if err != nil {
		return nil, fmt.Errorf("failed to select nodes: %w", err)
	}
	defer rows.Close()
	var nodes []Node
	for rows.Next() {
		var node Node
		err = rows.Scan(&node.Id, &node.ClusterId, &node.Ip, &node.PodCIDR, &node.Status, &node.LastHeartbeat)
		if err != nil {
			return nil, fmt.Errorf("failed to scan nodes: %w", err)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (r *Repository) GetClusterInfo(ctx context.Context) (*Cluster, error) {
	var cluster Cluster
	err := r.db.QueryRowContext(ctx, `SELECT * FROM cluster`).Scan(&cluster.Id, &cluster.PodCidr, &cluster.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster token: %w", err)
	}

	return &cluster, nil
}

//func (r *Repository) InsertPod(pod Pod)
