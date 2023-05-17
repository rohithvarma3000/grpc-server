package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"bytes"
	"encoding/json"
	

	"google.golang.org/grpc"
	pb "github.com/rohithvarma3000/grpc-server/comms"
)

const (
	port = 5000
	openaiURL = "https://api.openai.com/v1/chat/completions"
)

type OpenAIResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []Choices `json:"choices"`
	Usage   Usage     `json:"usage"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Choices struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) ChatReply(ctx context.Context, in *pb.Chat) (*pb.Reply, error) {
	bodyString := fmt.Sprintf(`{
		"model":"gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "%s"}]
	}`, in.Input)

	body := []byte(bodyString)

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer <YOUR_API_KEY>")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	response := &OpenAIResponse{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		panic(err)
	}

	return &pb.Reply{Output: response.Choices[0].Message.Content}, nil
}


func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Print("failed to serve:")
	}

}