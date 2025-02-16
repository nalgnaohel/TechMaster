package business

import (
	"03/internal/dialog"
	"03/internal/groq_client"
	"03/internal/models"
	"03/internal/word"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type groqBusiness struct {
	dialogRepo     dialog.DialogRepository
	wordRepo       word.WordRepository
	contextTimeout time.Duration
}

func NewGroqBusiness(dialogRepo dialog.DialogRepository, wordRepo word.WordRepository, timeout time.Duration) groq_client.GroqBusiness {
	return &groqBusiness{
		dialogRepo:     dialogRepo,
		wordRepo:       wordRepo,
		contextTimeout: timeout,
	}
}

func (gb *groqBusiness) ChatCompletion(c context.Context, groqClient *models.GroqClient, prompt string) (*string, *string, error) {
	ctx, cancel := context.WithTimeout(c, gb.contextTimeout)
	defer cancel()
	fmt.Println("Called!: ")
	groqMessage := make([]models.GroqMessage, 0)
	if prompt != "" {
		prompt += "Sau đó, từ đoạn hội thoại vừa tạo, lọc ra các từ quan trọng, bỏ qua danh từ riêng cần học. Danh sách các từ này được trả về dưới dạng JSON trong thẻ 'words'." +
			" Tiếp đó, dịch từng từ trong danh sách vừa tạo sang tiếng Anh, rồi trả về JSON gồm mảng trong đó mỗi phần tử gồm từ tiếng Việt (thẻ 'vi') và từ tiếng Anh tương đương (thẻ 'en')." +
			" Chỉ cần xuất ra hội thoại, 2 JSON danh sách từ và không cần giải thích. "
		prompt += "Không đánh số các câu trong hội thoại."
		userMessage := models.GroqMessage{
			Role:    "user",
			Content: prompt,
		}
		groqMessage = append(groqMessage, userMessage)
		//fmt.Println(groqMessage)
	} else {
		return nil, nil, errors.New("prompt is empty")
	}
	groqRequest := &models.GroqRequest{
		Messages:    groqMessage,
		LLMModel:    "deepseek-r1-distill-llama-70b",
		MaxTokens:   3500,
		Temperature: 0.2,
		TopP:        1,
		Stream:      false,
		Stop:        nil,
	}
	groqRequestJson, err := json.Marshal(groqRequest)
	if err != nil {
		return nil, nil, err
	}

	//send request to Groq API
	chatCompletionUrl := "https://api.groq.com/openai/v1/chat/completions"

	//fmt.Println("jsonbody: ", string(groqRequestJson))
	request, err := http.NewRequest(http.MethodPost, chatCompletionUrl, bytes.NewBuffer(groqRequestJson))

	//fmt.Println("buffer: ", bytes.NewBuffer(groqRequestJson), "\n")
	if err != nil {
		return nil, nil, err
	}

	//add Headers to post request
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+groqClient.ApiKey)

	//fmt.Printf("request: ", request, "\n")
	client := &http.Client{}
	resp, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, nil, err
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
		return nil, nil, errors.New("unexpected status code: " + resp.Status + " - " + bodyString)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	groqResponse := &models.GroqResponse{}
	err = json.Unmarshal(respBody, &groqResponse)
	if err != nil {
		return nil, nil, err
	}
	//save dialog to database
	var responseString string
	var dialogContent string
	if groqResponse.Choices != nil && len(groqResponse.Choices) > 0 {
		responseContent := groqResponse.Choices[0].Message.Content
		responseString = responseContent
		pos := strings.Index(responseString, "/think")
		if pos != -1 {
			st := pos + 7
			for responseString[st] == ' ' || responseString[st] == '\n' {
				st++
			}
			responseString = responseString[st:]
		}
		//fmt.Println("Dialog: ", responseString, " ", len(resCompos))
		poss := strings.Index(responseString, "\n\n")

		dialog := &models.Dialog{
			Lang:    "vi",
			Content: responseString[:poss],
		}
		createdDialog, err := gb.dialogRepo.Create(dialog)
		fmt.Println("Created dialog: ", createdDialog)

		if err != nil {
			return nil, nil, errors.New("error saving dialog to database")
		}
		dialogContent = responseString[:poss]
		//save words to database
		stID := strings.Index(responseString, "\"vi\":")
		if stID == -1 {
			return nil, &responseString, errors.New("no words")
		}
		responseString = responseString[stID:]
		cont1, cont2 := "", ""
		for i := 0; i < len(responseString)-6; i++ {
			//fmt.Println("i: ", responseString[i:i+5])
			if responseString[i:i+5] == "\"vi\":" {
				ed := strings.Index(responseString[i:], ",")
				cont1 = responseString[i+7 : i+ed-1]
				i += ed
			} else if responseString[i:i+5] == "\"en\":" {
				st := strings.Index(responseString[i+5:], "\"")
				if st == -1 {
					continue
				}
				ed := strings.Index(responseString[i+5+st+1:], "\"")
				cont2 = responseString[i+5+st+1 : i+6+st+ed]
				i += 6 + st + ed

				word := &models.Word{
					Lang:      "vi",
					Content:   cont1,
					Translate: cont2,
				}
				createdWord, err := gb.wordRepo.Create(word)
				if err != nil {
					if err.Error() != "word already exists" {
						return nil, &createdDialog.Content, errors.New("error saving word to database")
					}
				}
				fmt.Println("Created word: ", createdWord)
				er := gb.wordRepo.AddDialogWord(createdDialog.ID, createdWord.ID)
				if er != nil {
					if er.Error() != "dialog word already exists" {
						return nil, &createdWord.Content, errors.New("error saving dialog word to database")
					}
				}
			}
		}
	} else {
		return nil, nil, fmt.Errorf("no choices")
	}
	dialogContent = strings.TrimLeft(dialogContent, "Hội thoại:")
	responseString = "{\n  {" + responseString
	responseString = strings.TrimSuffix(responseString, "\n")
	responseString = strings.TrimSuffix(responseString, "```")
	return &responseString, &dialogContent, nil
}
