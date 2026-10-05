package main

import (
	"k8s-pet-project/shared/cloudapi/timeweb"
	"net/http"
	"os"
	"time"
)

func main() {
	//ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	_ = timeweb.NewTimewebClient(os.Getenv("TIMEWEB_TOKEN"), httpClient)

}
