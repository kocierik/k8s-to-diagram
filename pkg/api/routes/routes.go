package routes

import (
	"net/http"

	"github.com/kocierik/k8s-to-diagram/pkg/api/handlers"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/create-diagram", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB max
		handlers.CreateDiagramHandler(w, r)
	}))
	return mux
}
