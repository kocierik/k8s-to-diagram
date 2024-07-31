package api

import (
	"fmt"
	"net/http"

	"github.com/kocierik/k8s-to-diagram/pkg/api/routes"
)

func StartServer(port int) error {
	mux := routes.SetupRoutes()

	fmt.Printf("Server starting on port %d...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
