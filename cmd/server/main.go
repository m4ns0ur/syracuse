package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-toschool/syracuse/database"
	"github.com/go-toschool/syracuse/service"

	"github.com/go-toschool/syracuse"
	"github.com/go-toschool/syracuse/citizens"

	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

func main() {
	port := flag.Int64("port", 8001, "listening port")
	postgresDSN := flag.String("postgres-dsn", "postgres://gotoschool:goto1234@localhost:5432/drachma?sslmode=disable", "Postgres DSN")

	flag.Parse()
	pgSvc, err := database.NewPostgres(*postgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	srv := grpc.NewServer()

	citizens.RegisterCitizenshipServer(srv, &CitizensService{
		Citizens: &service.Citizens{
			Store: pgSvc,
		},
	})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting Syracuse service...")
	log.Println(fmt.Sprintf("Syracuse service, Listening on: %d", *port))
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// CitizensService grpc service implementation.
type CitizensService struct {
	Citizens syracuse.Citizens
}

// Get Gets a user by ID.
func (cs *CitizensService) Get(ctx context.Context, gr *citizens.GetRequest) (*citizens.GetResponse, error) {
	user, err := cs.Citizens.GetByID(gr.GetUserId())
	if err != nil {
		return nil, err
	}

	return &citizens.GetResponse{
		Data: user.ToProto(),
	}, nil
}

// Select return a collection of users.
func (cs *CitizensService) Select(ctx context.Context, gr *citizens.SelectRequest) (*citizens.SelectResponse, error) {
	users, err := cs.Citizens.Select()
	if err != nil {
		return nil, err
	}

	data := make([]*citizens.Citizen, 0)
	for _, user := range users {
		data = append(data, user.ToProto())
	}

	return &citizens.SelectResponse{
		Data: data,
	}, nil
}

// Create creates a new user into database.
func (cs *CitizensService) Create(ctx context.Context, gr *citizens.CreateRequest) (*citizens.CreateResponse, error) {
	email := gr.GetData().GetEmail()
	user, err := cs.Citizens.GetByEmail(email)
	if err != nil {
		citizen := &syracuse.Citizen{
			Email:    gr.GetData().GetEmail(),
			FullName: gr.GetData().GetFullName(),
		}

		if err := cs.Citizens.Create(citizen); err != nil {
			return nil, err
		}

		return &citizens.CreateResponse{
			Data: citizen.ToProto(),
		}, nil
	}

	return &citizens.CreateResponse{
		Data: user.ToProto(),
	}, nil
}

// Update updates a user.
func (cs *CitizensService) Update(ctx context.Context, gr *citizens.UpdateRequest) (*citizens.UpdateResponse, error) {
	user, err := cs.Citizens.GetByID(gr.GetUserId())
	if err != nil {
		return nil, err
	}

	user.Email = gr.GetData().GetEmail()
	user.FullName = gr.GetData().GetFullName()
	if err := cs.Citizens.Update(user); err != nil {
		return nil, err
	}

	return &citizens.UpdateResponse{
		Data: user.ToProto(),
	}, nil
}

// Delete delete a user.
func (cs *CitizensService) Delete(ctx context.Context, gr *citizens.DeleteRequest) (*citizens.DeleteResponse, error) {
	user, err := cs.Citizens.GetByID(gr.GetUserId())
	if err != nil {
		return nil, err
	}

	if err := cs.Citizens.Delete(user); err != nil {
		return nil, err
	}

	return &citizens.DeleteResponse{
		Data: user.ToProto(),
	}, nil
}
