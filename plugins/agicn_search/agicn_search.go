package agicn_search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	pluginName         = "AgicnSearch"
	pluginInputExample = "Who is Google boss?"
	pluginDesc         = `Search something by query input.`

	baseURL = "https://agicn-ducksearch.vercel.app/search"
)

type searchResponse struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Body  string `json:"body"`
}

type AgicnSearch struct {
	client *http.Client
}

func NewAgicnSearch() *AgicnSearch {
	c := &http.Client{}
	return &AgicnSearch{c}
}

func (s AgicnSearch) Do(ctx context.Context, query string) (answer string, err error) {
	searchResults, err := s.doHTTPRequest(ctx, query)
	if err != nil {
		return "", err
	}

	answer = s.makeAnswer(searchResults)
	return answer, nil
}

func (s AgicnSearch) doHTTPRequest(ctx context.Context, query string) ([]searchResponse, error) {
	params := url.Values{}
	params.Add("q", query)

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResults []searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return nil, err
	}

	return searchResults, nil
}

func (s AgicnSearch) makeAnswer(searchResults []searchResponse) string {

	if len(searchResults) == 0 {
		return ""
	}

	result := searchResults[0]

	return fmt.Sprintf("%s\n%s\n%s", result.Title, result.Body, result.Href)
}

func (s AgicnSearch) GetName() string {

	return pluginName
}

func (s AgicnSearch) GetInputExample() string {
	return pluginInputExample
}

func (s AgicnSearch) GetDesc() string {
	return pluginDesc
}
