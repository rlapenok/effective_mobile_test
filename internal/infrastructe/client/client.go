package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
	"github.com/rlapenok/effective_mobile_test/internal/domain/swagger_client"
	"go.uber.org/zap"
)

type Client struct {
	client *http.Client
	url    *string
}

func New(url *string) *Client {
	return &Client{client: &http.Client{}, url: url}
}
func (c *Client) createUrl(group string, song string) string {
	return fmt.Sprintf("%s/info?group=%s&song=%s", *c.url, url.PathEscape(group), url.PathEscape(song))
}

func (c *Client) GetSongDetail(group string, song string) (*models.SongDetails, error) {

	var songDetails models.SongDetails

	url := c.createUrl(group, song)
	resp, err := c.client.Get(url)
	if err != nil {
		zap.L().Error("Error while send request to SongDetailsApi", zap.Error(err))
		return nil, swagger_client.SwaggerClientError{Code: 500, Err: err}
	}
	if resp.StatusCode == http.StatusNotFound {
		zap.L().Error("SongDetailsApi - Not found", zap.Error(err))
		return nil, swagger_client.SwaggerClientError{Code: 404, Err: errors.New("SongDetailsApi - Not found")}
	}
	if err := json.NewDecoder(resp.Body).Decode(&songDetails); err != nil {
		zap.L().Error("Error decode response from SongDetailApi", zap.Error(err))
		return nil, &swagger_client.SwaggerClientError{Err: err}
	}
	defer resp.Body.Close()
	return &songDetails, nil

}
