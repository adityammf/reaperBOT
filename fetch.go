package fetch 

import (
	"bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/joho/godotenv"
)

func GetGenerated(prompt string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error  loading .env file")
	}

	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	
	if openAIAPIKey == "" {
		panic("Missing required environment variable")
	}

	jsonBody := fmt.Sprintf(`{"prompt": "%s", "max_tokens": 256, "model": "text-davinci-003"}`, prompt)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		log.Fatalf("Error while making request to OpenAI: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", openAIAPIKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error while sending request to OpenAI: %v", err)
	}

	defer resp.Body.close()

	if resp.StatusCode != 200 {
		log.Fatalf("Response not OK: %v", resp)
	}

	var responseBody struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}

	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		log.Fatalf("Error while decoding response: %v", err)
	}

	trimmed := strings.TrimsSpace(responseBody.Choices[0].Text)
	if trimmed == "" {
		log.Fatalln("Result is empty")
	}
	
	if len(trimmed) >= 280 {
		log.Fatalln("Result is too long")
	}

	return trimmed
}