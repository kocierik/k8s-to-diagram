package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/h2non/bimg"
)

func ConvertByteToSvg(svgBytes []byte) ([]byte, error) {
	bimg.VipsCacheSetMax(0)
	bimg.VipsCacheSetMaxMem(0)

	// Converti l'SVG in PNG
	image := bimg.NewImage(svgBytes)
	pngData, err := image.Convert(bimg.PNG)
	if err != nil {
		return nil, fmt.Errorf("could not convert SVG to PNG: %v", err)
	}
	return pngData, nil

	// err = bimg.Write("output.png", pngData)
	// if err != nil {
	//     fmt.Errorf(err)
	// }

}

// ConvertSVGToPNG takes an input SVG file path and converts it to a PNG file.
// It returns the path of the generated PNG file.
func ConvertSVGToPNG(svgPath string) error {
	// Ensure the file has an .svg extension
	if !strings.HasSuffix(svgPath, ".svg") {
		return fmt.Errorf("input file is not an SVG: %s", svgPath)
	}

	// Generate the PNG file path by changing the extension to .png
	pngPath := strings.TrimSuffix(svgPath, ".svg") + ".png"

	// Execute the rsvg-convert command to convert SVG to PNG
	cmd := exec.Command("rsvg-convert", "-o", pngPath, svgPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not convert SVG to PNG: %v", err)
	}

	return nil
}
