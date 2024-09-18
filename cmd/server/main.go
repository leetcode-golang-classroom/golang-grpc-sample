package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/leetcode-golang-classroom/golang-grpc-sample/proto"
	"google.golang.org/grpc"
)

type Person struct {
	ID          int32
	Name        string
	Email       string
	PhoneNumber string
}

var nextID int32 = 1
var persons = make(map[int32]Person)

type server struct {
	pb.UnimplementedPersonServiceServer
}

func (s *server) Create(ctx context.Context, in *pb.CreatePersonRequest) (*pb.PersonProfileResponse, error) {
	person := Person{Name: in.GetName(), Email: in.GetEmail(), PhoneNumber: in.GetPhoneNumber()}
	if person.Email == "" || person.Name == "" || person.PhoneNumber == "" {
		return &pb.PersonProfileResponse{}, errors.New("fields missing")
	}
	person.ID = nextID
	persons[person.ID] = person
	nextID = nextID + 1
	return &pb.PersonProfileResponse{
		Id:          person.ID,
		Email:       person.Email,
		Name:        person.Name,
		PhoneNumber: person.PhoneNumber,
	}, nil
}

func (s *server) Read(ctx context.Context, in *pb.SinglePersonRequest) (*pb.PersonProfileResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.PersonProfileResponse{}, errors.New("not found")
	}
	return &pb.PersonProfileResponse{
		Id: person.ID, Name: person.Name, Email: person.Email,
		PhoneNumber: person.PhoneNumber,
	}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdatePersonRequest) (*pb.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.SuccessResponse{Response: "Not found!"}, errors.New("not found")
	}
	person.Name = in.GetName()
	person.Email = in.GetEmail()
	person.PhoneNumber = in.GetPhoneNumber()
	if person.Email == "" || person.Name == "" || person.PhoneNumber == "" {
		return &pb.SuccessResponse{Response: "fields missing!"},
			errors.New("fields missing")
	}
	persons[person.ID] = person
	return &pb.SuccessResponse{Response: "Done!"}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.SinglePersonRequest) (*pb.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.SuccessResponse{Response: "Not found!"}, errors.New("not found")
	}
	delete(persons, id)
	return &pb.SuccessResponse{Response: "Delete!"}, nil
}
func Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":8080")
	errCh := make(chan error, 1)
	if err != nil {
		errCh <- fmt.Errorf("failed to listen: %v", err)
	}
	go func() {
		s := grpc.NewServer()
		pb.RegisterPersonServiceServer(s, &server{})
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			errCh <- fmt.Errorf("failed to serve: %v", err)
		}
		defer s.GracefulStop()
	}()
	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		fmt.Println("server is stopping")
		return lis.Close()
	}
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	err := Start(ctx)
	if err != nil {
		log.Println("failed to start app:", err)
	}
}
