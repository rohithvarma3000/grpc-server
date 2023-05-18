package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"encoding/json"
	

	"google.golang.org/grpc"
	pb "github.com/rohithvarma3000/grpc-server/comms"
)

const (
	port = 5000
	apiURL = "https://api.api-ninjas.com/v1/nutrition?query="
	apiKey = "<YOUR_API_KEY>"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

type APIResponse struct {
	Name                string  `json:"name"`
	Calories            float64 `json:"calories"`
	ServingSizeG        float64 `json:"serving_size_g"`
	FatTotalG           float64 `json:"fat_total_g"`
	FatSaturatedG       float64 `json:"fat_saturated_g"`
	ProteinG            float64 `json:"protein_g"`
	SodiumMg            int     `json:"sodium_mg"`
	PotassiumMg         int     `json:"potassium_mg"`
	CholesterolMg       int     `json:"cholesterol_mg"`
	CarbohydratesTotalG float64 `json:"carbohydrates_total_g"`
	FiberG              float64 `json:"fiber_g"`
	SugarG              float64 `json:"sugar_g"`
}


func (s *server) ChatReply(ctx context.Context, in *pb.Chat) (*pb.Reply, error) {
	req, err := http.NewRequest("GET", apiURL + in.Input, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	response := []APIResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	return &pb.Reply{Output: fmt.Sprintf("%v",response[0].Calories)}, nil
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