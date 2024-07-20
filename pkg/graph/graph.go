package graph

import (
	"encoding/json"
	"fmt"

	"github.com/kocierik/k8s-to-diagram/pkg/types"
)

func GenerateMermaidGraph(manifests []types.K8sResource) string {
	graph := "graph TD;\n"
	resourceMap := make(map[string]types.K8sResource)

	// First pass: create nodes
	for _, resource := range manifests {
		var comm types.Communication
		if commAnnotation, ok := resource.Metadata.Annotations["communication"]; ok {
			err := json.Unmarshal([]byte(commAnnotation), &comm)
			if err != nil {
				fmt.Printf("Error unmarshalling communication annotation for %s: %v\n", resource.Metadata.Name, err)
				continue
			}
			if comm.Name != "" {
				// Use comm.Name as the key in resourceMap
				id := comm.Name
				resourceMap[id] = resource
				graph += fmt.Sprintf("    %s[%s: %s];\n", id, resource.Kind, resource.Metadata.Name)
			}
		}
	}

	// Second pass: create connections
	for _, resource := range manifests {
		var comm types.Communication
		if commAnnotation, ok := resource.Metadata.Annotations["communication"]; ok {
			err := json.Unmarshal([]byte(commAnnotation), &comm)
			if err != nil {
				fmt.Printf("Error unmarshalling communication annotation for %s: %v\n", resource.Metadata.Name, err)
				continue
			}
			sourceID := comm.Name

			// Handle outbound communications
			for _, outbound := range comm.Outbound {
				targetID := outbound.Service
				if _, exists := resourceMap[targetID]; exists {
					graph += fmt.Sprintf("    %s --> |port %d| %s;\n", sourceID, outbound.Port, targetID)
				} else {
					fmt.Printf("Warning: Outbound service %s not found for %s\n", outbound.Service, sourceID)
				}
			}

			// Handle inbound communications
			for _, inbound := range comm.Inbound {
				sourceServiceID := inbound.Service
				if _, exists := resourceMap[sourceServiceID]; exists {
					graph += fmt.Sprintf("    %s --> |port %d| %s;\n", sourceServiceID, inbound.Port, sourceID)
				} else {
					fmt.Printf("Warning: Inbound service %s not found for %s\n", inbound.Service, sourceID)
				}
			}
		}
	}

	return graph
}
