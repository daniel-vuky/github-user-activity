package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <GITHUB_USER_TOKEN>")
		return
	}
	username := os.Args[1]
	if username == "" {
		fmt.Println("Error: GITHUB_USER_TOKEN is required")
		return
	}
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s/events", username), nil)
	if err != nil {
		fmt.Println("Error fetching events:", err)
		return
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-Github-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GITHUB_USER_TOKEN")))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error fetching events:", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error fetching events:", response.Status)
		return
	}
	var events []struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Actor struct {
			Id           int    `json:"id"`
			Login        string `json:"login"`
			DisplayLogin string `json:"display_login"`
			GravatarId   string `json:"gravatar_id"`
			Url          string `json:"url"`
			AvatarUrl    string `json:"avatar_url"`
		}
		Repo struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Url  string `json:"url"`
		}
		Payload struct {
			Action string `json:"action"`
		}
		Public    bool   `json:"public"`
		CreatedAt string `json:"created_at"`
	}
	err = json.NewDecoder(response.Body).Decode(&events)
	if err != nil {
		fmt.Println("Error decoding events:", err)
		return
	}
	fmt.Printf("%-20s %-20s %-30s %-30s %-10s %-30s\n", "Event ID", "Event Type", "Actor", "Repo", "Public", "Created At")
	for _, event := range events {
		fmt.Printf("%-20s %-20s %-30s %-30s %-10t %-30s\n", event.Id, event.Type, event.Actor.Login, event.Repo.Name, event.Public, event.CreatedAt)
	}
	fmt.Println("Events fetched successfully.")
	fmt.Println("Total events:", len(events))
}
