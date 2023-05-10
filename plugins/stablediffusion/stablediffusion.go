package stablediffusion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pluginName         = "StableDiffusion"
	pluginInputExample = "A beautiful girl"
	pluginDesc         = `Stable diffusion is text-to-image model capable of generating images given any text input`
)

type StableDiffusion struct {
	sdAddr string

	client *http.Client
}

func NewStableDiffusion(sdAddr string) *StableDiffusion {
	if len(sdAddr) == 0 {
		panic("stable diffusion address is empty")
	}

	return &StableDiffusion{
		sdAddr: sdAddr,
		client: &http.Client{},
	}
}

func (s *StableDiffusion) Do(ctx context.Context, query string) (answer string, err error) {

	url := fmt.Sprintf("http://%v/sd", s.sdAddr)

	resp, err := s.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var sdResp struct {
		Result bool     `json:"result"`
		Images []string `json:"images"` // base64
	}

	if err := json.NewDecoder(resp.Body).Decode(&sdResp); err != nil {
		return "", err
	}

	if len(sdResp.Images) == 0 {
		return "", nil
	}

	return sdResp.Images[0], nil
}

func (StableDiffusion) GetName() string {
	return pluginName
}

func (StableDiffusion) GetInputExample() string {
	return pluginInputExample
}

func (StableDiffusion) GetDesc() string {
	return pluginDesc
}
