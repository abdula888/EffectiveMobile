package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Структура для ответа от AudD API
type AudDResponse struct {
	Result []struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Lyrics string `json:"lyrics"`
		Media  string `json:"media"` // Временно используем строку для media
	} `json:"result"`
}

// Функция для обращения к AudD API
func GetAudDData(artist string, track string) (AudDResponse, error) {
	apiKey := os.Getenv("AUDD_API_KEY")
	apiURL := fmt.Sprintf("https://api.audd.io/findLyrics/?q=%s+%s&api_token=%s", url.QueryEscape(artist), url.QueryEscape(track), apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return AudDResponse{}, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AudDResponse{}, fmt.Errorf("API returned non-OK status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AudDResponse{}, fmt.Errorf("error reading API response: %v", err)
	}

	var audDResponse AudDResponse
	err = json.Unmarshal(body, &audDResponse)
	if err != nil {
		return AudDResponse{}, fmt.Errorf("error parsing JSON from API: %v", err)
	}

	return audDResponse, nil
}
