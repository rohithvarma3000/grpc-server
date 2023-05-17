package main

import (
	"context"
	"fmt"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/rohithvarma3000/grpc-server/comms"
)

const (
	input = "Hello"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewChatServiceClient(conn)
	reply, err := client.ChatReply(context.Background(), &pb.Chat{Input: input})
	if err != nil {
		panic(err)
	}
	fmt.Println(reply.Output)


}