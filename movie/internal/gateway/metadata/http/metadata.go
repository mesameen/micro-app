package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/movie/internal/gateway"
	"github.com/mesameen/micro-app/pkg/discovery"
	"github.com/mesameen/micro-app/pkg/logger"
)

// Gateway defines a movie metadata HTTP gateway
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

// Get gets movie metadata by movie id
func (g *Gateway) GetMovieDetails(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"
	logger.Infof("URL to get movie metadata: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", err)
	}
	var metadata *model.Metadata
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	logger.Infof("Repsonse got from metadata service: %s", string(resBytes))
	err = json.Unmarshal(resBytes, metadata)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
