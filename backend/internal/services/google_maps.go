package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/fooder/backend/internal/models"
)

type GoogleMapsService struct {
	apiKey string
}

func NewGoogleMapsService() *GoogleMapsService {
	return &GoogleMapsService{
		apiKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
	}
}

type PlacesSearchResponse struct {
	Results []struct {
		PlaceID          string `json:"place_id"`
		Name             string `json:"name"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

type PlaceDetailsResponse struct {
	Result struct {
		PlaceID                string `json:"place_id"`
		Name                   string `json:"name"`
		FormattedAddress       string `json:"formatted_address"`
		FormattedPhoneNumber   string `json:"formatted_phone_number"`
		InternationalPhoneNumber string `json:"international_phone_number"`
		Website                string `json:"website"`
		Geometry               struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"result"`
	Status string `json:"status"`
}

func (s *GoogleMapsService) SearchPlaces(query string) ([]models.GooglePlaceResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("Google Maps API key not configured")
	}

	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	params := url.Values{}
	params.Set("query", query+" restaurant")
	params.Set("type", "restaurant")
	params.Set("key", s.apiKey)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search places: %w", err)
	}
	defer resp.Body.Close()

	var searchResp PlacesSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if searchResp.Status != "OK" && searchResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Google Maps API error: %s", searchResp.Status)
	}

	results := make([]models.GooglePlaceResult, 0, len(searchResp.Results))
	for _, r := range searchResp.Results {
		results = append(results, models.GooglePlaceResult{
			PlaceID:   r.PlaceID,
			Name:      r.Name,
			Address:   r.FormattedAddress,
			Latitude:  r.Geometry.Location.Lat,
			Longitude: r.Geometry.Location.Lng,
		})
	}

	return results, nil
}

func (s *GoogleMapsService) GeocodeCities(query string) ([]models.GooglePlaceResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("Google Maps API key not configured")
	}

	baseURL := "https://maps.googleapis.com/maps/api/place/autocomplete/json"
	params := url.Values{}
	params.Set("input", query)
	params.Set("types", "(cities)")
	params.Set("key", s.apiKey)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to geocode: %w", err)
	}
	defer resp.Body.Close()

	var autoResp struct {
		Predictions []struct {
			PlaceID     string `json:"place_id"`
			Description string `json:"description"`
			StructuredFormatting struct {
				MainText      string `json:"main_text"`
				SecondaryText string `json:"secondary_text"`
			} `json:"structured_formatting"`
		} `json:"predictions"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&autoResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if autoResp.Status != "OK" && autoResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Google Maps API error: %s", autoResp.Status)
	}

	// For each prediction, get the coordinates
	results := make([]models.GooglePlaceResult, 0, len(autoResp.Predictions))
	for _, p := range autoResp.Predictions {
		// Get place details to retrieve coordinates
		detailsURL := "https://maps.googleapis.com/maps/api/place/details/json"
		detailParams := url.Values{}
		detailParams.Set("place_id", p.PlaceID)
		detailParams.Set("fields", "geometry")
		detailParams.Set("key", s.apiKey)

		detailResp, err := http.Get(detailsURL + "?" + detailParams.Encode())
		if err != nil {
			continue
		}

		var details struct {
			Result struct {
				Geometry struct {
					Location struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"location"`
				} `json:"geometry"`
			} `json:"result"`
			Status string `json:"status"`
		}

		if err := json.NewDecoder(detailResp.Body).Decode(&details); err != nil {
			detailResp.Body.Close()
			continue
		}
		detailResp.Body.Close()

		if details.Status == "OK" {
			results = append(results, models.GooglePlaceResult{
				PlaceID:   p.PlaceID,
				Name:      p.StructuredFormatting.MainText,
				Address:   p.Description,
				Latitude:  details.Result.Geometry.Location.Lat,
				Longitude: details.Result.Geometry.Location.Lng,
			})
		}
	}

	return results, nil
}

func (s *GoogleMapsService) GetPlaceDetails(placeID string) (*models.GooglePlaceResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("Google Maps API key not configured")
	}

	baseURL := "https://maps.googleapis.com/maps/api/place/details/json"
	params := url.Values{}
	params.Set("place_id", placeID)
	params.Set("fields", "place_id,name,formatted_address,geometry,formatted_phone_number,international_phone_number,website")
	params.Set("key", s.apiKey)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get place details: %w", err)
	}
	defer resp.Body.Close()

	var detailsResp PlaceDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&detailsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if detailsResp.Status != "OK" {
		return nil, fmt.Errorf("Google Maps API error: %s", detailsResp.Status)
	}

	phone := detailsResp.Result.InternationalPhoneNumber
	if phone == "" {
		phone = detailsResp.Result.FormattedPhoneNumber
	}

	return &models.GooglePlaceResult{
		PlaceID:   detailsResp.Result.PlaceID,
		Name:      detailsResp.Result.Name,
		Address:   detailsResp.Result.FormattedAddress,
		Latitude:  detailsResp.Result.Geometry.Location.Lat,
		Longitude: detailsResp.Result.Geometry.Location.Lng,
		Phone:     phone,
		Website:   detailsResp.Result.Website,
	}, nil
}
