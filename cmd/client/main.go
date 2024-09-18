package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/leetcode-golang-classroom/golang-grpc-sample/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPersonServiceClient(conn)
	// timeout for ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Example 1: Create new person
	fmt.Println("Creating a new person...")
	createReq := &pb.CreatePersonRequest{
		Name:        "John Wick",
		Email:       "john.wick@test.com",
		PhoneNumber: "123-456-789",
	}
	createRes, err := client.Create(ctx, createReq)
	if err != nil {
		log.Fatalf("Error during Create: %v", err)
	}
	fmt.Printf("Person created: %+v\n", createReq)

	// Example2: Read the created person by ID
	fmt.Println("Reading the person by ID...")
	readReq := &pb.SinglePersonRequest{
		Id: createRes.GetId(),
	}
	readRes, err := client.Read(ctx, readReq)
	if err != nil {
		log.Fatalf("Error during Read: %v", err)
	}
	fmt.Printf("Person details: %+v\n", readRes)

	// Example3: Update the person's details
	fmt.Println("Updating the person's details...")
	updateReq := &pb.UpdatePersonRequest{
		Id:          createRes.GetId(),
		Name:        "Luke Skywalker",
		Email:       "luke.skywaler@test.com",
		PhoneNumber: "083-111-000",
	}
	updateRes, err := client.Update(ctx, updateReq)
	if err != nil {
		log.Fatalf("Error during Update: %v", err)
	}
	fmt.Printf("Update response: %s\n", updateRes.GetResponse())
	// Example 4: Delete the person by ID
	fmt.Println("Deleting the person by ID...")
	deleteReq := &pb.SinglePersonRequest{
		Id: createRes.GetId(),
	}
	deleteRes, err := client.Delete(ctx, deleteReq)
	if err != nil {
		log.Fatalf("Error during delete: %v", err)
	}
	fmt.Printf("Delete response: %s\n", deleteRes.GetResponse())
}
