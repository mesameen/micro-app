package http

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/mesameen/micro-app/movie/internal/gateway"
	"github.com/mesameen/micro-app/pkg/discovery"
	"github.com/mesameen/micro-app/pkg/logger"
	"github.com/mesameen/micro-app/rating/pkg/model"
)

// Gateway defines an HTTP gateway to rating service
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a rating service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	query := req.URL.Query()
	query.Add("id", string(recordID))
	query.Add("type", string(recordType))
	req.URL.RawQuery = query.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", err)
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	logger.Infof("Response got from rating service. %s", string(resBytes))
	v, err := strconv.ParseFloat(string(resBytes), 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// PutRating writes a rating
func (g *Gateway) PutRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", rating.RecordType)
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", res)
	}
	return nil
}
