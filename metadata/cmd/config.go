package main

type config struct {
	API              apiConfig        `yaml:"api"`
	ServiceDiscovery serviceDiscovery `yaml:"serviceDiscovery"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type serviceDiscovery struct {
	Consul consulConfig `yaml:"consul"`
}

type consulConfig struct {
	Address string `yaml:"address"`
}
