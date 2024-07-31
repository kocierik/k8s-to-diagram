package render

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"oss.terrastruct.com/d2/d2compiler"
	"oss.terrastruct.com/d2/d2exporter"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

func RenderD2Graph(graph1 string) ([]byte, error) {
	fmt.Println("grafico --> ", graph1)
	graph, config, err := d2compiler.Compile("", strings.NewReader(graph1), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to compile graph: %w", err)
	}

	// Cambiamo il tema in uno scuro
	darkTheme := d2themescatalog.DarkCatalog[0].ID
	graph.ApplyTheme(darkTheme)

	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return nil, fmt.Errorf("failed to create text ruler: %w", err)
	}
	if err := graph.SetDimensions(nil, ruler, nil); err != nil {
		return nil, fmt.Errorf("failed to set graph dimensions: %w", err)
	}
	if err := d2dagrelayout.Layout(context.Background(), graph, nil); err != nil {
		return nil, fmt.Errorf("failed to layout graph: %w", err)
	}
	diagram, err := d2exporter.Export(context.Background(), graph, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to export diagram: %w", err)
	}
	diagram.Config = config
	out, err := d2svg.Render(diagram, &d2svg.RenderOpts{
		ThemeID: &darkTheme,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to render SVG: %w", err)
	}
	if err := os.WriteFile(filepath.Join("images/k8s_infrastructure.svg"), out, 0600); err != nil {
		return nil, fmt.Errorf("failed to write SVG file: %w", err)
	}
	return out, nil
}
