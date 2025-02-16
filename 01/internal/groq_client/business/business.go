package business

import (
	"app01/internal/groq_client"
	"app01/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type groqBusiness struct {
}

func NewGroqBusiness() groq_client.GroqBusiness {
	return &groqBusiness{}
}

func (gb *groqBusiness) ChatCompletion(groqClient *models.GroqClient, prompt string) (*string, error) {
	groqMessage := make([]models.GroqMessage, 0)
	if prompt != "" {
		userMessage := models.GroqMessage{
			Role:    "user",
			Content: prompt,
		}
		groqMessage = append(groqMessage, userMessage)
		//fmt.Println(groqMessage)
	} else {
		return nil, errors.New("prompt is empty")
	}
	groqRequest := &models.GroqRequest{
		Messages:    groqMessage,
		LLMModel:    "deepseek-r1-distill-llama-70b",
		MaxTokens:   2048,
		Temperature: 0.2,
		TopP:        1,
		Stream:      false,
		Stop:        nil,
	}
	groqRequestJson, err := json.Marshal(groqRequest)
	if err != nil {
		return nil, err
	}

	//send request to Groq API
	chatCompletionUrl := "https://api.groq.com/openai/v1/chat/completions"

	//fmt.Println("jsonbody: ", string(groqRequestJson))
	request, err := http.NewRequest(http.MethodPost, chatCompletionUrl, bytes.NewBuffer(groqRequestJson))

	//fmt.Println("buffer: ", bytes.NewBuffer(groqRequestJson), "\n")
	if err != nil {
		return nil, err
	}

	//add Headers to post request
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+groqClient.ApiKey)

	//fmt.Printf("request: ", request, "\n")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Response Body:", bodyString)
		return nil, errors.New("unexpected status code: " + resp.Status + " - " + bodyString)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	groqResponse := &models.GroqResponse{}
	err = json.Unmarshal(respBody, &groqResponse)
	if err != nil {
		return nil, err
	}
	return &groqResponse.Choices[0].Message.Content, nil
}
