package controlplane

import (
	"context"
	"k8s-pet-project/internal/controlplane/application"
	"k8s-pet-project/internal/controlplane/delivery"
	"k8s-pet-project/internal/controlplane/repository"
	"k8s-pet-project/internal/database"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.NewStorage("/opt/controlplane.db")
	if err != nil {
		panic(err)
	}

	repo := repository.NewRepository(db)
	if err := repo.InitTables(ctx); err != nil {
		panic(err)
	}

	useCase := application.NewUseCase(&repo)
	_ = delivery.New(useCase)
}
