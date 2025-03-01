package APICalls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/logfire-sh/cli/pkg/cmd/sql/models"
)

func GetRecommendations(token string, endpoint string, teamId string, role string) (models.RecommendResponse, error) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	req, err := http.NewRequest("GET", endpoint+"/ai/teams/"+teamId+"/sql-recommend?role="+strings.ReplaceAll(role, " ", "-"), nil)
	if err != nil {
		return models.RecommendResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return models.RecommendResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.RecommendResponse{}, err
	}

	var RecommendationsResp models.RecommendResponse
	err = json.Unmarshal(body, &RecommendationsResp)
	if err != nil {
		return models.RecommendResponse{}, err
	}

	return RecommendationsResp, err
}

func GetFilterRecommendations(token string, endpoint string, teamId string, role string) (models.RecommendFilterResponse,
	error) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	req, err := http.NewRequest("GET", endpoint+"/ai/teams/"+teamId+"/filter-recommend?role="+strings.ReplaceAll(role,
		" ", "-"), nil)
	if err != nil {
		return models.RecommendFilterResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return models.RecommendFilterResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.RecommendFilterResponse{}, err
	}

	var RecommendationsResp models.RecommendFilterResponse
	err = json.Unmarshal(body, &RecommendationsResp)
	if err != nil {
		return models.RecommendFilterResponse{}, err
	}

	return RecommendationsResp, err
}
