package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
)

func UpdateSource(client *http.Client, token, endpoint string, teamid, sourceid, sourcename string) (models.Source, error) {
	data := models.SourceCreate{
		Name: sourcename,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return models.Source{}, err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/team/"+teamid+"/source/"+sourceid, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.Source{}, err
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

		return models.Source{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Source{}, err
	}

	var sourceUpdateResp models.SourceCreateResponse
	err = json.Unmarshal(body, &sourceUpdateResp)
	if err != nil {
		return models.Source{}, err
	}

	if !sourceUpdateResp.IsSuccessful {
		return sourceUpdateResp.Data, errors.New("source update failed")
	}

	return sourceUpdateResp.Data, nil
}

func GetAllSources(client *http.Client, token, endpoint string, teamId string) ([]models.Source, error) {
	url := fmt.Sprintf(endpoint+"api/team/%s/source", teamId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {

		return []models.Source{}, err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return []models.Source{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []models.Source{}, err
	}

	var sourceResp models.SourcesResponse

	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		return []models.Source{}, err
	}

	if !sourceResp.IsSuccessful {
		return []models.Source{}, errors.New(sourceResp.Message[0])
	}

	return sourceResp.Data, nil
}

func GetSource(token, endpoint string, teamId, sourceId string) (models.Source, error) {
	client := http.Client{}

	url := fmt.Sprintf(endpoint+"api/team/%s/source/%s", teamId, sourceId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Source{}, err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return models.Source{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Source{}, err
	}

	var sourceResp models.SourceResponse

	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		return models.Source{}, err
	}

	if !sourceResp.IsSuccessful {
		return models.Source{}, errors.New(sourceResp.Message[0])
	}

	return sourceResp.Data, nil
}

func CreateSource(token, endpoint string, teamId, sourceName, platform string) (models.Source, error) {
	client := http.Client{}

	// platform should be mapped to its respective int as sourceType, for kubernetes its 1
	sourceType, exists := models.PlatformMap[strings.ToLower(platform)]
	if !exists {
		return models.Source{}, errors.New("invalid platform")
	}

	data := models.SourceCreate{
		Name:       sourceName,
		SourceType: sourceType,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return models.Source{}, err
	}

	url := endpoint + "api/team/" + teamId + "/source"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.Source{}, err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return models.Source{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Source{}, err
	}

	var sourceResp models.SourceCreateResponse
	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		return models.Source{}, err
	}

	if !sourceResp.IsSuccessful {
		fmt.Print(sourceResp)
		return models.Source{}, errors.New("failed to create source")
	}

	return sourceResp.Data, nil
}

func GetSchema(token, endpoint, teamId string, sourceids []string) ([]map[string]string, error) {
	client := http.Client{}

	idsParam := strings.Join(sourceids, "&")

	url := fmt.Sprintf("%s/api/team/%s/schema?%s", endpoint, teamId, idsParam)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Non-200 response code")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var data []map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	return data, nil
}

func GetConfiguration(token, endpoint string, teamId, sourceId string) (interface{}, error) {
	client := http.Client{}

	url := fmt.Sprintf(endpoint+"api/team/%s/source/%s/configuration", teamId, sourceId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.ConfigurationResponse{}, err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return models.ConfigurationResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ConfigurationResponse{}, err
	}

	var configResp models.ConfigurationResponse

	err = json.Unmarshal(body, &configResp)
	if err != nil {
		return models.ConfigurationResponse{}, err
	}

	if !configResp.IsSuccessful {
		return models.ConfigurationResponse{}, errors.New(configResp.Message[0])
	}

	return configResp.Data, nil
}
