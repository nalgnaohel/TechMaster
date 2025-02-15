package business

import (
	"03/internal/dialog"
	"03/internal/groq_client"
	"03/internal/models"
	"03/internal/word"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type groqBusiness struct {
	dialogRepo dialog.DialogRepository
	wordRepo   word.WordRepository
}

func NewGroqBusiness(dialogRepo dialog.DialogRepository, wordRepo word.WordRepository) groq_client.GroqBusiness {
	return &groqBusiness{
		dialogRepo: dialogRepo,
		wordRepo:   wordRepo,
	}
}

func (gb *groqBusiness) ChatCompletion(groqClient *models.GroqClient, prompt string) (*string, error) {
	groqMessage := make([]models.GroqMessage, 0)
	if prompt != "" {
		prompt += "Sau đó, từ đoạn hội thoại vừa tạo, lọc ra các từ quan trọng, bỏ qua danh từ riêng cần học. Danh sách các từ này được trả về dưới dạng JSON trong thẻ 'words'. Tiếp đó, dịch từng từ trong danh sách vừa tạo sang tiếng Anh, rồi trả về JSON gồm mảng trong đó mỗi phần tử gồm từ tiếng Việt và từ tiếng Anh tương đương. Chỉ cần xuất ra hội thoại, 2 JSON danh sách từ và không cần giải thích."
		userMessage := models.GroqMessage{
			Role:    "user",
			Content: prompt,
		}
		groqMessage = append(groqMessage, userMessage)
		fmt.Println(groqMessage)
	} else {
		return nil, errors.New("prompt is empty")
	}
	groqRequest := &models.GroqRequest{
		Messages:    groqMessage,
		LLMModel:    "deepseek-r1-distill-llama-70b",
		MaxTokens:   1024,
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

	response, err := http.NewRequest(http.MethodPost, chatCompletionUrl, bytes.NewBuffer(groqRequestJson))
	if err != nil {
		return nil, err
	}

	//add Headers to post request
	response.Header.Set("Content-Type", "application/json")
	response.Header.Set("Authorization", "Bearer "+groqClient.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(response)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("unexpected status code: reason: " + resp.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	groqResponse := &models.GroqResponse{}
	err = json.Unmarshal(respBody, &groqResponse)
	if err != nil {
		return nil, err
	}

	//save dialog to database
	var responseString string
	if groqResponse.Choices != nil && len(groqResponse.Choices) > 0 {
		responseContent := groqResponse.Choices[0].Message.Content
		resCompos := strings.Split(responseContent, "json")
		responseString = resCompos[0][0 : len(resCompos[0])-5]
		dialog := &models.Dialog{
			Lang:    "vi",
			Content: responseString,
		}
		createdDialog, err := gb.dialogRepo.Create(dialog)
		if err != nil {
			return &responseString, errors.New("error saving dialog to database")
		}

		//save words to database
		translated := resCompos[2][5 : len(resCompos[2])-5]
		translatedWords := strings.Split(translated, "{")[1:]
		for i := 0; i < len(translatedWords); i++ {
			//fmt.Println(translatedWords[i])
			translatedWords[i] = strings.TrimLeft(translatedWords[i], " ")
			if i == len(translatedWords)-1 {
				translatedWords[i] = translatedWords[i][0 : len(translatedWords[i])-2]
			} else {
				translatedWords[i] = translatedWords[i][0 : len(translatedWords[i])-5]
			}
			fmt.Println("-----\n" + translatedWords[i])
			translatedWordsParts := strings.Split(translatedWords[i], ",")
			fmt.Println("Parts: ", translatedWordsParts)
			content := ""
			translate := ""
			for j := 0; j < len(translatedWordsParts); j++ {
				translatedWordsParts[j] = strings.TrimLeft(translatedWordsParts[j], " ")
				toSave := strings.Split(translatedWordsParts[j], ":")
				toSave[1] = strings.TrimLeft(toSave[1], " ")
				fmt.Println(toSave)
				lang := "vi"
				if j%2 == 1 {
					lang = "en"
					translate = toSave[1]
				} else {
					content = toSave[1]
				}
				fmt.Println(lang + ": " + toSave[1])
			}
			word := &models.Word{
				Lang:      "vi",
				Content:   content,
				Translate: translate,
			}
			createdWord, err := gb.wordRepo.Create(word)
			if err != nil {
				return &responseString, errors.New("error saving word to database")
			}
			fmt.Println("Created word: ", createdWord)
			er := gb.wordRepo.AddDialogWord(createdDialog.ID, createdWord.ID)
			if er != nil {
				return &responseString, errors.New("error saving dialog word to database")
			}
		}
	} else {
		return nil, fmt.Errorf("no choices")
	}
	return &responseString, nil
}
