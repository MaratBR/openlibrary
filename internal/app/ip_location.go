package app

import (
	"context"
	"math/rand/v2"
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

var placeholderIPLocations = []IPLocation{
	{Country: "United States", Region: "California", City: "San Francisco"},
	{Country: "Germany", Region: "Berlin", City: "Berlin"},
	{Country: "Japan", Region: "Tokyo", City: "Tokyo"},
	{Country: "Brazil", Region: "São Paulo", City: "São Paulo"},
	{Country: "Australia", Region: "New South Wales", City: "Sydney"},
	{Country: "Canada", Region: "Ontario", City: "Toronto"},
}

func (placeholderIPLocationService) Locate(context.Context, string) (IPLocation, error) {
	// TODO: Replace placeholder data with a real IP geolocation provider.
	return placeholderIPLocations[rand.IntN(len(placeholderIPLocations))], nil
}

func NewIPLocationService() IPLocationService { return placeholderIPLocationService{} }
