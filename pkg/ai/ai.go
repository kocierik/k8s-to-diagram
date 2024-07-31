package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func AnalyzeWithGemini(imagePath string) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-1.5-flash")

	imgData1, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	prompt := []genai.Part{
		genai.ImageData("png", imgData1),
		genai.Text("What do you think about this architecture, can you improve it?"),
	}
	resp, err := model.GenerateContent(ctx, prompt...)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Candidates[0].Content)

}
