package configs

type Client struct {
	*Grpc `yaml:"grpc"`
}

type Grpc struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
