package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	pb "github.com/rohithvarma3000/grpc-server/comms"
)

const (
	port = 5000
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) ChatReply(ctx context.Context, in *pb.Chat) (*pb.Reply, error) {
	return &pb.Reply{Output: "Hello " + in.Input}, nil
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