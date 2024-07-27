package graph

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/kocierik/k8s-to-diagram/pkg/types"
)

func sanitizeName(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return reg.ReplaceAllString(name, "_")
}

func getShapeForKind(kind string) string {
	switch kind {
	case "Deployment":
		return "rectangle"
	case "Service":
		return "square"
	case "Pod":
		return "parallelogram"
	case "ConfigMap":
		return "diamond"
	case "Secret":
		return "hexagon"
	case "Ingress":
		return "stored_data"
	case "StatefulSet":
		return "cylinder"
	case "PersistentVolume":
		return "cloud"
	case "PersistentVolumeClaim":
		return "cylinder"
	default:
		return "rectangle" // Default shape
	}
}

func GenerateD2Graph(manifests []types.K8sResource) string {
	graph := "Diagram Architecture {\n"
	resourceMap := make(map[string]types.K8sResource)

	// First pass: create nodes and clusters
	for _, resource := range manifests {
		var comm types.Communication
		if commAnnotation, ok := resource.Metadata.Annotations["communication"]; ok {
			err := json.Unmarshal([]byte(commAnnotation), &comm)
			if err != nil {
				fmt.Printf("Error unmarshalling communication annotation for %s: %v\n", resource.Metadata.Name, err)
				continue
			}
			if comm.Name != "" {
				id := sanitizeName(comm.Name)
				resourceMap[id] = resource
				replicas := resource.Spec.Replicas
				if replicas == 0 {
					replicas = 1
				}
				graph += fmt.Sprintf(" %s {\n", id)
				graph += fmt.Sprintf(" label: \"%s\"\n", resource.Metadata.Name)

				// Get shape based on Kind
				shape := getShapeForKind(resource.Kind)

				for i := 0; i < replicas; i++ {
					replicaID := sanitizeName(fmt.Sprintf("%s_%d", comm.Name, i))
					graph += fmt.Sprintf(" %s: {\n", replicaID)
					graph += fmt.Sprintf("  label: \"%s[%d]\"\n", resource.Metadata.Name, i+1)
					graph += fmt.Sprintf("  shape: %s\n", shape)
					graph += " }\n"
				}
				graph += " }\n"
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

			sourceID := sanitizeName(comm.Name)

			// Handle outbound communications
			for _, outbound := range comm.Outbound {
				targetID := sanitizeName(outbound.Service)
				if _, exists := resourceMap[targetID]; exists {
					graph += fmt.Sprintf("  %s -> %s: \"port %d\"\n", sourceID, targetID, outbound.Port)
				} else {
					fmt.Printf("Warning: Outbound service %s not found for %s\n", outbound.Service, sourceID)
				}
			}

			// Handle inbound communications
			for _, inbound := range comm.Inbound {
				sourceServiceID := sanitizeName(inbound.Service)
				if _, exists := resourceMap[sourceServiceID]; exists {
					graph += fmt.Sprintf("  %s -> %s: \"port %d\"\n", sourceServiceID, sourceID, inbound.Port)
				} else {
					fmt.Printf("Warning: Inbound service %s not found for %s\n", inbound.Service, sourceID)
				}
			}
		}
	}

	graph += "}\n"
	return graph
}
