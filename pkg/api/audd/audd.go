package audd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// Структура для получения ссылки из ответа AudD API
type Media struct {
	Provider string `json:"provider"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

// Функция для обращения к AudD API
func GetAudDData(artist, track, apiKey, apiURL string) (AudDResponse, error) {
	fullURL := fmt.Sprintf("%s?q=%s+%s&api_token=%s", apiURL, url.QueryEscape(artist), url.QueryEscape(track), apiKey)
	fmt.Println(fullURL)
	resp, err := http.Get(fullURL)
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
