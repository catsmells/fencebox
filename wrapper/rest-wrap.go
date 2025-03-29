package esri
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)
// temporary scheme until server is figured out
type EsriClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}
type Feature struct {
	Attributes map[string]interface{} `json:"attributes"`
	Geometry   map[string]interface{} `json:"geometry"`
}
type FeatureResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
func NewEsriClient(baseURL, token string) *EsriClient {
	return &EsriClient{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}
func (c *EsriClient) sendRequest(endpoint string, method string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?f=json&token=%s", c.BaseURL, endpoint, c.Token)
	var req *http.Request
	var err error
	if payload != nil {
		jsonData, _ := json.Marshal(payload)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}
	return body, nil
}
func (c *EsriClient) AddFeature(featureService string, feature Feature) (*FeatureResponse, error) {
	payload := map[string]interface{}{
		"features": []Feature{feature},
	}
	body, err := c.sendRequest(fmt.Sprintf("%s/addFeatures", featureService), "POST", payload)
	if err != nil {
		return nil, err
	}
	var response FeatureResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
func (c *EsriClient) UpdateFeature(featureService string, feature Feature) (*FeatureResponse, error) {
	payload := map[string]interface{}{
		"features": []Feature{feature},
	}
	body, err := c.sendRequest(fmt.Sprintf("%s/updateFeatures", featureService), "POST", payload)
	if err != nil {
		return nil, err
	}
	var response FeatureResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
func (c *EsriClient) QueryFeature(featureService, where string) ([]Feature, error) {
	query := fmt.Sprintf("%s/query?where=%s&outFields=*&f=json&token=%s", featureService, where, c.Token)
	body, err := c.sendRequest(query, "GET", nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Features []Feature `json:"features"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Features, nil
}
