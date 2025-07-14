package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/mesameen/micro-app/src/pkg/discovery"
	"github.com/mesameen/micro-app/src/pkg/logger"
)

// Registry defines in memory registry
type Registry struct {
	sync.RWMutex
	serviceAddrs map[string]map[string]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service
// registry instance
func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: make(map[string]map[string]*serviceInstance),
	}
}

// Register creates a service instance record in the registry.
func (r *Registry) Register(
	ctx context.Context,
	instanceID string,
	serviceName string,
	hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[string]*serviceInstance{}
	}
	r.serviceAddrs[serviceName][instanceID] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
	return nil
}

// Deregister removes record from the registry
func (r *Registry) Deregister(
	ctx context.Context,
	instanceID string,
	serviceName string,
) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName], instanceID)
	return nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry
func (r *Registry) ReportHealthyState(
	ctx context.Context,
	instanceID string,
	serviceName string,
) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("service " + serviceName + " is not registered yet")
	}
	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("instance " + instanceID + " of service " + serviceName + " is not registered yet")
	}
	r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances
// of the requested service
func (r *Registry) ServiceAddresses(
	ctx context.Context,
	serviceName string,
) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var addresses []string
	for i, v := range r.serviceAddrs[serviceName] {
		// if last active 5 seconds earlier consider it isn't healthy
		if v.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			logger.Infof("Instance %s of service %s is not active, skipping", i, serviceName)
			continue
		}
		addresses = append(addresses, v.hostPort)
	}
	return addresses, nil
}
