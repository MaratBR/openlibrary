package app

import (
	"context"
)

type IPLocation struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

type IPLocationService interface {
	Locate(context.Context, string) (IPLocation, error)
}

type placeholderIPLocationService struct{}

func (placeholderIPLocationService) Locate(context.Context, string) (IPLocation, error) {
	// No verified provider is configured. Never manufacture location evidence.
	return IPLocation{}, nil
}

func NewIPLocationService() IPLocationService { return placeholderIPLocationService{} }
