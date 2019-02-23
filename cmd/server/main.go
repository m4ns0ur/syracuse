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
	c, err := cs.Citizens.GetByID(gr.GetUserId())
	if err != nil {
		return nil, err
	}

	return &citizens.GetResponse{
		Data: c.ToProto(),
	}, nil
}

// Select return a collection of users.
func (cs *CitizensService) Select(ctx context.Context, gr *citizens.SelectRequest) (*citizens.SelectResponse, error) {
	cc, err := cs.Citizens.Select()
	if err != nil {
		return nil, err
	}

	data := make([]*citizens.Citizen, 0)
	for _, c := range cc {
		data = append(data, c.ToProto())
	}

	return &citizens.SelectResponse{
		Data: data,
	}, nil
}

// Create creates a new user into database.
func (cs *CitizensService) Create(ctx context.Context, gr *citizens.CreateRequest) (*citizens.CreateResponse, error) {
	email := gr.GetData().GetEmail()
	u, err := cs.Citizens.GetByEmail(email)
	if err != nil {
		c := &syracuse.Citizen{
			Email:    gr.Data.Email,
			FullName: gr.Data.FullName,
		}

		if err := cs.Citizens.Create(c); err != nil {
			return nil, err
		}

		return &citizens.CreateResponse{
			Data: c.ToProto(),
		}, nil
	}

	return &citizens.CreateResponse{
		Data: u.ToProto(),
	}, nil
}

// Update updates a user.
func (cs *CitizensService) Update(ctx context.Context, gr *citizens.UpdateRequest) (*citizens.UpdateResponse, error) {
	u, err := cs.Citizens.GetByID(gr.GetUserId())
	if err != nil {
		return nil, err
	}

	u.Email = gr.Data.Email
	u.FullName = gr.Data.FullName
	if err := cs.Citizens.Update(u); err != nil {
		return nil, err
	}

	return &citizens.UpdateResponse{
		Data: u.ToProto(),
	}, nil
}

// Delete delete a user.
func (cs *CitizensService) Delete(ctx context.Context, gr *citizens.DeleteRequest) (*citizens.DeleteResponse, error) {
	u, err := cs.Citizens.GetByID(gr.GetUserId())
	if err != nil {
		return nil, err
	}

	if err := cs.Citizens.Delete(u); err != nil {
		return nil, err
	}

	return &citizens.DeleteResponse{
		Data: u.ToProto(),
	}, nil
}
