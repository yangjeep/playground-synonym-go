package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type Completion struct {
	Text string `json:"text"`
}

type Choice struct {
	Text string `json:"text"`
}

type OpenAIResponse struct {
	Completions []Completion `json:"completions"`
}

func generateSynonyms(apiKey string, word string, numSynonyms int) ([]string, error) {
	prompt := fmt.Sprintf("Generate %d synonyms for %s:", numSynonyms, word)
	data := url.Values{
		"prompt":      {prompt},
		"max_tokens":  {"100"},
		"temperature": {"0.5"},
		"n":           {"1"},
		"stop":        {""},
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/engines/text-davinci-002/completions", strings.NewReader(data.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	message := response.Completions[0].Text

	re := regexp.MustCompile(`"(.+)"`)
	synonyms := re.FindAllString(message, numSynonyms)
	return synonyms, nil
}

func main() {
	godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	c := gogpt.NewClient(apiKey)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:       gogpt.GPT3Davinci,
		MaxTokens:   100,
		Temperature: 0.5,
		Prompt:      "List 10 synonyms for swimsuit",
		//Stream:    true,
	}

	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return
	}

	message := resp.Choices[0].Text

	fmt.Println(message)

	/*
		re := regexp.MustCompile(`"(.+)"`)
		synonyms := re.FindAllString(message, 10)

		for _, synonym := range synonyms {
			fmt.Println(synonym)
		}
	*/
}
