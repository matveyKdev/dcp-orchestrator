package repository

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

type DatabaseInterface interface {
	InitTables(ctx context.Context) error
	GetAllPods(ctx context.Context) ([]Pod, error)
	InsertPod(ctx context.Context, pod Pod) error
	UpdatePod(ctx context.Context, pod Pod) error
	DeletePod(ctx context.Context, pod Pod) error
}

func NewRepository(db *sql.DB) DatabaseInterface {
	return &Repository{db: db}
}

func (r *Repository) InitTables(ctx context.Context) error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS node_info (
			id TEXT PRIMARY KEY,
			pod_cidr TEXT,
			ip TEXT,
		)
	`)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS pods (
			id TEXT PRIMARY KEY,
			service_name TEXT,
			ip TEXT,
			port TEXT,
			status TEXT,
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllPods(ctx context.Context) ([]Pod, error) {
	rows, err := r.db.Query(`SELECT * FROM pods`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pods []Pod
	for rows.Next() {
		var pod Pod
		err = rows.Scan(pod.ServiceName, pod.Ip, pod.Port, pod.Status)
		if err != nil {
			return nil, err
		}
		pods = append(pods, pod)
	}
	return pods, nil
}

func (r *Repository) InsertPod(ctx context.Context, pod Pod) error {
	_, err := r.db.Exec(`INSERT INTO pods (service_name, ip, port, status) VALUES (?, ?, ?, ?)`, pod.ServiceName, pod.Ip, pod.Port, pod.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdatePod(ctx context.Context, pod Pod) error {
	_, err := r.db.Exec(`UPDATE pods SET service_name = ?, ip = ?, port = ?, status = ?`,
		pod.ServiceName, pod.Ip, pod.Port, pod.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeletePod(ctx context.Context, pod Pod) error {
	_, err := r.db.Exec(`DELETE * FROM pods WHERE service_name = ?, ip = ? AND port = ?`,
		pod.ServiceName, pod.Ip, pod.Port)
	if err != nil {
		return err
	}
	return nil
}
