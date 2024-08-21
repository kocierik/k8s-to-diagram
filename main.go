package main

import (
	"fmt"
	"github.com/kocierik/k8s-to-diagram/pkg/graph"
	"github.com/kocierik/k8s-to-diagram/pkg/manifests"
	"github.com/kocierik/k8s-to-diagram/pkg/render"
)

func main() {
	manifestDir := "./manifests"
	manifests, err := manifests.ReadManifests(manifestDir)
	if err != nil {
		fmt.Printf("Error reading manifests: %v\n", err)
		return
	}

	graphData := graph.GenerateD2Graph(manifests)
	render.RenderD2Graph(graphData)

	if err != nil {
		fmt.Printf("Error converting SVG to PNG: %v\n", err)
		return
	}
	fmt.Println("\n\n Diagram generated and analyzed successfully.")
}
