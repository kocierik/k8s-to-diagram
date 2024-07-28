package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kocierik/k8s-to-diagram/pkg/ai"
	"github.com/kocierik/k8s-to-diagram/pkg/graph"
	"github.com/kocierik/k8s-to-diagram/pkg/manifests"
	"github.com/kocierik/k8s-to-diagram/pkg/render"
	"github.com/kocierik/k8s-to-diagram/pkg/utils"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	manifestDir := "./manifests"
	manifests, err := manifests.ReadManifests(manifestDir)
	if err != nil {
		fmt.Printf("Error reading manifests: %v\n", err)
		return
	}

	graphData := graph.GenerateD2Graph(manifests)
	render.RenderD2Graph(graphData)

	// Step 3: Send the Diagram to Gemini for Analysis
	imagePath := "images/k8s_infrastructure.svg"
	err = utils.ConvertSVGToPNG(imagePath)
	if err != nil {
		fmt.Printf("Error converting SVG to PNG: %v\n", err)
		return
	}
	imagePath = "images/k8s_infrastructure.png"
	ai.AnalyzeWithGemini(imagePath)
	fmt.Println("\n\n Diagram generated and analyzed successfully.")
}
