package main

type config struct {
	API              apiConfig        `yaml:"api"`
	ServiceDiscovery serviceDiscovery `yaml:"serviceDiscovery"`
	Jaeger           jaegerConfig     `yaml:"jaeger"`
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

type jaegerConfig struct {
	URL string `yaml:"url"`
}
