package manifests

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kocierik/k8s-to-diagram/pkg/types"
	"gopkg.in/yaml.v3"
)

func ReadManifests(dir string) ([]types.K8sResource, error) {
	var manifests []types.K8sResource
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var resource types.K8sResource
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
