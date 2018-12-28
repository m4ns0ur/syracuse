package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-toschool/syracuse"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/go-toschool/syracuse/postgres"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

func main() {
	port := flag.Int64("port", 8001, "listening port")
	postgresDSN := flag.String("postgres-dsn", "postgres://localhost:5432/syracuse?sslmode=disable", "Postgres DSN")

	flag.Parse()
	fmt.Println(*postgresDSN)
	db, err := sqlx.Connect("postgres", *postgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	srv := grpc.NewServer()

	citizens.RegisterCitizenshipServer(srv, &CitizensService{
		Citizens: &postgres.CitizensService{
			Store: db,
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
	c, err := cs.Citizens.Get(&syracuse.CitizensQuery{
		ID: gr.GetUserId(),
	})
	if err != nil {
		return nil, err
	}

	return &citizens.GetResponse{
		Data: &citizens.Citizen{
			Id:        c.ID,
			Email:     c.Email,
			FullName:  c.FullName,
			CreatedAt: c.CreatedAt.Unix(),
			UpdatedAt: c.UpdatedAt.Unix(),
		},
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
		data = append(data, &citizens.Citizen{
			Id:        c.ID,
			Email:     c.Email,
			FullName:  c.FullName,
			CreatedAt: c.CreatedAt.Unix(),
			UpdatedAt: c.UpdatedAt.Unix(),
		})
	}

	return &citizens.SelectResponse{
		Data: data,
	}, nil
}

// Create creates a new user into database.
func (cs *CitizensService) Create(ctx context.Context, gr *citizens.CreateRequest) (*citizens.CreateResponse, error) {
	c := &syracuse.Citizen{
		Email:    gr.Data.Email,
		FullName: gr.Data.FullName,
	}

	u, err := cs.Citizens.Get(&syracuse.CitizensQuery{
		Email: gr.Data.Email,
	})
	if err != nil {
		if err := cs.Citizens.Create(c); err != nil {
			return nil, err
		}

		return &citizens.CreateResponse{
			Data: &citizens.Citizen{
				Id:        c.ID,
				Email:     c.Email,
				FullName:  c.FullName,
				CreatedAt: c.CreatedAt.Unix(),
				UpdatedAt: c.UpdatedAt.Unix(),
			},
		}, nil
	}

	return &citizens.CreateResponse{
		Data: &citizens.Citizen{
			Id:        u.ID,
			Email:     u.Email,
			FullName:  u.FullName,
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
		},
	}, nil
}

// Update updates a user.
func (cs *CitizensService) Update(ctx context.Context, gr *citizens.UpdateRequest) (*citizens.UpdateResponse, error) {
	u, err := cs.Citizens.Get(&syracuse.CitizensQuery{
		ID: gr.UserId,
	})
	if err != nil {
		return nil, err
	}

	u.Email = gr.Data.Email
	u.FullName = gr.Data.FullName
	if err := cs.Citizens.Update(u); err != nil {
		return nil, err
	}

	return &citizens.UpdateResponse{
		Data: &citizens.Citizen{
			Id:        u.ID,
			Email:     u.Email,
			FullName:  u.FullName,
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
		},
	}, nil
}

// Delete delete a user.
func (cs *CitizensService) Delete(ctx context.Context, gr *citizens.DeleteRequest) (*citizens.DeleteResponse, error) {
	u, err := cs.Citizens.Get(&syracuse.CitizensQuery{
		ID: gr.UserId,
	})
	if err != nil {
		return nil, err
	}

	if err := cs.Citizens.Delete(u); err != nil {
		return nil, err
	}

	return &citizens.DeleteResponse{
		Data: &citizens.Citizen{
			Id:        u.ID,
			Email:     u.Email,
			FullName:  u.FullName,
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
		},
	}, nil
}
