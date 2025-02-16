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
		prompt += "Sau đó, từ đoạn hội thoại vừa tạo, lọc ra các từ quan trọng, bỏ qua danh từ riêng cần học. Danh sách các từ này được trả về dưới dạng JSON trong thẻ 'words'." +
			" Tiếp đó, dịch từng từ trong danh sách vừa tạo sang tiếng Anh, rồi trả về JSON gồm mảng trong đó mỗi phần tử gồm từ tiếng Việt (thẻ 'vi') và từ tiếng Anh tương đương (thẻ 'en')." +
			" Chỉ cần xuất ra hội thoại, 2 JSON danh sách từ và không cần giải thích. Mảng và 2 danh sách từ được phân cách nhau bởi câu \"---------\"."
		prompt += "Không đánh số các câu trong hội thoại."
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
	//save dialog to database
	var responseString string
	if groqResponse.Choices != nil && len(groqResponse.Choices) > 0 {
		responseContent := groqResponse.Choices[0].Message.Content
		resCompos := strings.Split(responseContent, "json")
		responseString = resCompos[0][0:]
		fmt.Println("Dialog: ", responseString)
		pos := strings.Index(responseString, "/think")
		if pos != -1 {
			responseString = responseString[pos+9:]
		}
		poss := strings.Index(responseString, "---------")
		responseString = responseString[0 : poss-2]

		dialog := &models.Dialog{
			Lang:    "vi",
			Content: responseString,
		}
		createdDialog, err := gb.dialogRepo.Create(dialog)
		fmt.Println("Created dialog: ", createdDialog)
		if err != nil {
			return &responseString, errors.New("error saving dialog to database")
		}
		fmt.Println("List vi-en:\n", resCompos[2])
		//save words to database
		transPairs := strings.Split(resCompos[2], "},")
		for _, pair := range transPairs {
			startViID := strings.Index(pair, "\"vi\":")
			startEnID := strings.Index(pair, "\"en\":")
			st1, st2 := min(startViID, startEnID), max(startViID, startEnID)
			ed1, ed2 := st2-1, len(pair)-1
			for pair[ed1] != '"' {
				ed1--
			}
			for pair[ed2] != '"' {
				ed2--
			}
			content1 := pair[st1+7 : ed1]
			content2 := pair[st2+7 : ed2]
			fmt.Println("content1: ", content1)
			fmt.Println("content2: ", content2)

			word := &models.Word{
				Lang:      "vi",
				Content:   content1,
				Translate: content2,
			}
			createdWord, err := gb.wordRepo.Create(word)
			if err != nil {
				if err.Error() != "word already exists" {
					return &responseString, errors.New("error saving word to database")
				}
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
