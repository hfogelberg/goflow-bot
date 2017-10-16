package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Response struct {
	ID        string    `db:"id" json:"id"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Lang      string    `db:"lang" json:"lang"`
	Result    struct {
		Source           string `db:"source" json:"source"`
		ResolvedQuery    string `db:"resolved_query" json:"resolvedQuery"`
		Action           string `db:"action" json:"action"`
		ActionIncomplete bool   `db:"action_incomplete" json:"actionIncomplete"`
		Parameters       struct {
		} `db:"parameters" json:"parameters"`
		Contexts []interface{} `db:"contexts" json:"contexts"`
		Metadata struct {
		} `db:"metadata" json:"metadata"`
		Fulfillment struct {
			Speech   string `db:"speech" json:"speech"`
			Messages []struct {
				Type   int    `db:"type" json:"type"`
				ID     string `db:"id" json:"id"`
				Speech string `db:"speech" json:"speech"`
			} `db:"messages" json:"messages"`
		} `db:"fulfillment" json:"fulfillment"`
		Score float64 `db:"score" json:"score"`
	} `db:"result" json:"result"`
	Status struct {
		Code      int    `db:"code" json:"code"`
		ErrorType string `db:"error_type" json:"errorType"`
	} `db:"status" json:"status"`
	SessionID string `db:"session_id" json:"sessionId"`
}

func main() {
	msg := "Where do you live"
	if err := chat(msg); err != nil {
		return
	}

	fmt.Println("OK!")
}

func chat(message string) error {
	var response Response

	accesstoken := os.Getenv("DF_CLIENT_ACCESS_TOKEN")
	client := &http.Client{}

	q, err := buildQueryString(message)
	req, err := http.NewRequest("GET", q, nil)
	if err != nil {
		log.Printf("Error sending request %s\n", err.Error())
		return err
	}
	key := fmt.Sprintf("Bearer %s", accesstoken)
	req.Header.Set("Authorization", key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error posting to Dialogflow %s\n", err.Error())
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println(err)
	}

	fmt.Printf("Response\n%s\n", response.Result.Fulfillment.Speech)

	return nil
}

func buildQueryString(message string) (string, error) {
	var queryURL *url.URL
	baseURL := "https://api.api.ai/v1/query"

	queryURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("Error parsing query %s\n", err.Error())
		return "", err
	}
	parameters := url.Values{}
	parameters.Add("v", "20150910")
	parameters.Add("lang", "en")
	parameters.Add("query", message)
	parameters.Add("sessionId", fmt.Sprint(time.Now().Unix()))
	queryURL.RawQuery = parameters.Encode()
	q := queryURL.String()

	return q, nil
}
