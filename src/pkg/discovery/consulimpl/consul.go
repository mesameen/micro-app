package consulimpl

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/mesameen/micro-app/src/pkg/discovery"
)

// Registry defines a Consul based service registry
type Registry struct {
	client *consul.Client
}

// NewRegistry creates a new consul based registry instance
func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{
		client: client,
	}, nil
}

// Register creates a service instance record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return fmt.Errorf("hostPort must be in a form of <host>:<port>, example localhost:8091")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	// registering a service with ttl 5 seconds
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	})
}

// Deregister removes a service instance record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddresses returns the list of  addresses of active instances
// of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}
	var addresses []string
	for _, e := range entries {
		addresses = append(addresses, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return addresses, nil
}

// ReportHealthyState is a push mechanism for reporting the healthy state
// to the registry.
func (r *Registry) ReportHealthyState(ctx context.Context, instanceID string, serviceName string) error {
	// r.client.Agent().PassTTL()
	return r.client.Agent().UpdateTTL(instanceID, "output", "pass")
}
