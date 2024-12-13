package lastfm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Структура для ответа от Last.fm API
type LastFmResponse struct {
	Track struct {
		Name   string `json:"name"`
		Artist struct {
			Name string `json:"name"`
		} `json:"artist"`
		Album struct {
			Title string `json:"title"`
		} `json:"album"`
		Wiki struct {
			Published string `json:"published"`
		} `json:"wiki"`
	} `json:"track"`
}

// Функция для обращения к Last.fm API
func GetLastFmData(artist, track, apiKey, apiURL string) (LastFmResponse, error) {
	fullURL := fmt.Sprintf("%s&api_key=%s&artist=%s&track=%s&format=json", apiURL, apiKey, url.QueryEscape(artist), url.QueryEscape(track))

	resp, err := http.Get(fullURL)
	if err != nil {
		return LastFmResponse{}, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LastFmResponse{}, fmt.Errorf("API returned non-OK status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LastFmResponse{}, fmt.Errorf("error reading API response: %v", err)
	}

	var lastFmResponse LastFmResponse
	err = json.Unmarshal(body, &lastFmResponse)
	if err != nil {
		return LastFmResponse{}, fmt.Errorf("error parsing JSON from API: %v", err)
	}

	return lastFmResponse, nil
}
