package http

import (
	"net/http"
	"testing"
)

func TestGetSongs(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/songs/")
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, но получен %d", resp.StatusCode)
	}
}

func TestGetSongText(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/groups/Test/songs/Test")
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, но получен %d", resp.StatusCode)
	}

	resp, err = http.Get("http://localhost:8080/groups/Test/songs/Test1")
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Ожидался статус 500, но получен %d", resp.StatusCode)
	}
}
