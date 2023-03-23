package app

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

// var messages []Message

func getOpenAIResponse(apiKey string, question string) OpenaiResponse {
	// Messegesにはこれまでのやり取りを入れるっぽい？
	message := Message{
		Role:    "user",
		Content: question,
	}
	requestBody := OpenaiRequest{
		Model:    "gpt-3.5-turbo",
		Messages: []Message{message},
	}

	requestJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response OpenaiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		println("Error: ", err.Error())
		return OpenaiResponse{}
	}

	// messagesの中にやり取りの記録を追加してる？
	// messages = append(messages, Message{
	// 	Role:    "assistant",
	// 	Content: response.Choices[0].Messages.Content,
	// })

	return response
}
