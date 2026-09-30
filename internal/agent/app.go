package agent

import (
	"context"
	"k8s-pet-project/internal/agent/repository"
	"k8s-pet-project/internal/database"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.NewStorage("/opt/agent.db")
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewRepository(db)
	err = repo.InitTables(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":6432", nil))
}
