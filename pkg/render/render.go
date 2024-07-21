package render

import (
	"context"
	"fmt"
	"os"

	"github.com/dreampuf/mermaid.go"
)

func RenderMermaidGraph(graph string) {
	ctx := context.Background()
	re, err := mermaid_go.NewRenderEngine(ctx)
	if err != nil {
		fmt.Printf("Error creating render engine: %v\n", err)
		return
	}
	defer re.Cancel()

	svgContent, err := re.Render(graph)
	if err != nil {
		fmt.Printf("Error rendering SVG: %v\n", err)
		return
	}
	err = os.WriteFile("images/k8s_infrastructure.svg", []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("Error writing SVG file: %v\n", err)
	}

	pngContent, _, err := re.RenderAsPng(graph)
	if err != nil {
		fmt.Printf("Error rendering PNG: %v\n", err)
		return
	}
	err = os.WriteFile("images/k8s_infrastructure.png", pngContent, 0644)
	if err != nil {
		fmt.Printf("Error writing PNG file: %v\n", err)
	}

	fmt.Println("Infrastructure schema generated successfully.")
}
