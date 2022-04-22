package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/googleapis/go-type-adapters/adapters"
	petv1 "github.com/mfridman/buf-generate-twirp-go/go/pet/v1"
	"github.com/twitchtv/twirp"
)

func main() {
	r := chi.NewRouter()
	// Twirp Handler.
	petStoreServiceHandler := petv1.NewPetStoreServiceServer(&petStoreService{})
	r.Mount(petStoreServiceHandler.PathPrefix(), petStoreServiceHandler)
	// Twirp Client.
	petStoreClient := petv1.NewPetStoreServiceProtobufClient("http://localhost:8080", http.DefaultClient)

	// Client hits server evert 1 seconds.
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			resp, err := petStoreClient.GetPet(context.TODO(), &petv1.GetPetRequest{})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp.GetPet().GetName(), resp.GetPet().GetPetType())
		}
	}()
	// Run server. Blocking. ctrl+c to stop.
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

type petStoreService struct{}

func (s *petStoreService) GetPet(context.Context, *petv1.GetPetRequest) (*petv1.GetPetResponse, error) {
	ts, err := adapters.TimeToProtoDateTime(time.Now())
	if err != nil {
		return nil, err
	}
	return &petv1.GetPetResponse{
		Pet: &petv1.Pet{
			PetType:   petv1.PetType_PET_TYPE_DOG,
			PetId:     "lab-123",
			Name:      "dante",
			CreatedAt: ts,
		},
	}, nil
}

func (s *petStoreService) PutPet(context.Context, *petv1.PutPetRequest) (*petv1.PutPetResponse, error) {
	return nil, twirp.NewError(twirp.Unimplemented, "unimplemented")
}

func (s *petStoreService) DeletePet(context.Context, *petv1.DeletePetRequest) (*petv1.DeletePetResponse, error) {
	return nil, twirp.NewError(twirp.Unimplemented, "unimplemented")
}

func (s *petStoreService) PurchasePet(context.Context, *petv1.PurchasePetRequest) (*petv1.PurchasePetResponse, error) {
	return nil, twirp.NewError(twirp.Unimplemented, "unimplemented")
}
