package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kocierik/k8s-to-diagram/pkg/ai"
	"github.com/kocierik/k8s-to-diagram/pkg/graph"
	"github.com/kocierik/k8s-to-diagram/pkg/manifests"
	"github.com/kocierik/k8s-to-diagram/pkg/render"
)

func CreateDiagramHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create a temporary directory to store the uploaded files
	tempDir, err := ioutil.TempDir("", "manifests")
	if err != nil {
		http.Error(w, "Failed to create temp directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Process uploaded files
	files := r.MultipartForm.File["manifests"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create a new file in the temp directory
		tempFile, err := os.Create(filepath.Join(tempDir, fileHeader.Filename))
		if err != nil {
			http.Error(w, "Failed to create temp file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		// Copy the uploaded file to the temp file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Parse manifests from the temp directory
	k8sResources, err := manifests.ReadManifests(tempDir)
	if err != nil {
		http.Error(w, "Failed to parse manifests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate graph
	graphString := graph.GenerateD2Graph(k8sResources)

	// Render diagram
	image, err := render.RenderD2Graph(graphString)
	if err != nil {
		http.Error(w, "Failed to render diagram: "+err.Error(), http.StatusInternalServerError)
		return
	}

	analysis := ai.AnalyzeWithGemini(image)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Diagram generated and analyzed successfully",
		"analysis": analysis,
	})
}
