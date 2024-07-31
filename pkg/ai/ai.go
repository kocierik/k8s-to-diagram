package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/kocierik/k8s-to-diagram/pkg/utils"
	"google.golang.org/api/option"
)

func AnalyzeWithGemini(image []byte) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	pngImage, err := utils.ConvertByteToSvg(image)

	if err != nil {
		log.Fatal(err)
	}

	prompt := []genai.Part{
		genai.ImageData("png", pngImage),
		genai.Text("What do you think about this architecture, can you improve it?"),
	}
	resp, err := model.GenerateContent(ctx, prompt...)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content)

}
