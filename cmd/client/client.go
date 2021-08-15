package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cjosti/fc3-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUsersStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Cleiton",
		Email: "osti@teste.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	log.Println(res)

}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Osti",
		Email: "o@o.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "O1",
			Name:  "Osti",
			Email: "o@o.com",
		},
		{
			Id:    "O2",
			Name:  "Osti",
			Email: "o@o.com",
		},
		{
			Id:    "O3",
			Name:  "Osti",
			Email: "o@o.com",
		},
		{
			Id:    "O4",
			Name:  "Osti",
			Email: "o@o.com",
		},
		{
			Id:    "O5",
			Name:  "Osti",
			Email: "o@o.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving request: %v", err)
	}

	fmt.Println(res)
}

func AddUsersStreamBoth(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "O1",
			Name:  "Osti1",
			Email: "o@o.com",
		},
		{
			Id:    "O2",
			Name:  "Osti2",
			Email: "o@o.com",
		},
		{
			Id:    "O3",
			Name:  "Osti3",
			Email: "o@o.com",
		},
		{
			Id:    "O4",
			Name:  "Osti4",
			Email: "o@o.com",
		},
		{
			Id:    "O5",
			Name:  "Osti5",
			Email: "o@o.com",
		},
	}

	stream, err := client.AddUsersStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending User: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
