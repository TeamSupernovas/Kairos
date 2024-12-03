package service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/location"
)
type GeoPoint struct {
	Latitude	float64
	Longitude	float64
}

type GeoLocationService struct {
	client *location.Client
}

func NewGeoLocationService(awsConfig aws.Config) *GeoLocationService {
	client :=  location.NewFromConfig(awsConfig)
	return &GeoLocationService{client: client}
}

func (geoLocationService *GeoLocationService)CalculateGeoPoint(address, placeIndexName string) (GeoPoint, error) {
	output, err := geoLocationService.client.SearchPlaceIndexForText(context.Background(), &location.SearchPlaceIndexForTextInput{
		Text:      aws.String(address),
		IndexName: aws.String(placeIndexName),
		MaxResults: aws.Int32(1),
	})
	if err != nil {
		return GeoPoint{}, fmt.Errorf("failed to search place index: %w", err)
	}

	// Extract coordinates from the response
	if len(output.Results) == 0 {
		return GeoPoint{}, fmt.Errorf("no results found for the given address")
	}
	coords := output.Results[0].Place.Geometry.Point
	return GeoPoint{Latitude: coords[1], Longitude: coords[0]}, nil
}