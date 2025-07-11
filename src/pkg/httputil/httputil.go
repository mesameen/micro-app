package httputil

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/mesameen/micro-app/src/pkg/discovery"
)

// ServiceConnection attempts to select a random service
// instance and returns a httpClient connection to it
func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (string, *http.Client, error) {
	addrs, err := registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return "", nil, err
	}
	addr := addrs[rand.Intn(len(addrs))]
	return addr, &http.Client{}, nil
}
