package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreampuf/mermaid.go"
	"gopkg.in/yaml.v3"
)

type K8sResource struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string            `yaml:"name"`
		Annotations map[string]string `yaml:"annotations"`
	} `yaml:"metadata"`
}

type Communication struct {
	Name     string                `json:"name"`
	Inbound  []CommunicationDetail `json:"inbound"`
	Outbound []CommunicationDetail `json:"outbound"`
}

type CommunicationDetail struct {
	Service string `json:"service"`
	Port    int    `json:"port"`
}

func readManifests(dir string) ([]K8sResource, error) {
	var manifests []K8sResource
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var resource K8sResource
			err = yaml.Unmarshal(content, &resource)
			if err != nil {
				return err
			}
			manifests = append(manifests, resource)
		}
		return nil
	})
	return manifests, err
}

func renderMermaidGraph(graph string) {
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
	err = os.WriteFile("k8s_infrastructure.svg", []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("Error writing SVG file: %v\n", err)
	}
	pngContent, _, err := re.RenderAsPng(graph)
	if err != nil {
		fmt.Printf("Error rendering PNG: %v\n", err)
		return
	}
	err = os.WriteFile("k8s_infrastructure.png", pngContent, 0644)
	if err != nil {
		fmt.Printf("Error writing PNG file: %v\n", err)
	}
	fmt.Println("Infrastructure schema generated successfully.")
}

func generateMermaidGraph(manifests []K8sResource) string {
	graph := "graph TD;\n"
	resourceMap := make(map[string]K8sResource)

	// First pass: create nodes
	for _, resource := range manifests {
		var comm Communication
		if commAnnotation, ok := resource.Metadata.Annotations["communication"]; ok {
			err := json.Unmarshal([]byte(commAnnotation), &comm)
			if err != nil {
				fmt.Printf("Error unmarshalling communication annotation for %s: %v\n", resource.Metadata.Name, err)
				continue
			}
		}

		if comm.Name != "" {
			id := fmt.Sprintf("%s_%s", resource.Kind, comm.Name)
			resourceMap[id] = resource
			graph += fmt.Sprintf("    %s[%s:%s];\n", id, resource.Kind, resource.Metadata.Name)
		}
	}

	return graph
}

func main() {
	manifestDir := "./manifests"
	manifests, err := readManifests(manifestDir)
	if err != nil {
		fmt.Printf("Error reading manifests: %v\n", err)
		return
	}

	graph := generateMermaidGraph(manifests)
	renderMermaidGraph(graph)
}

