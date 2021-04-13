package configs

type Config struct {
	*Server   `yaml:"server"`
	*Client   `yaml:"client"`
	*Database `yaml:"database"`
}

type Database struct {
	*Mysql         `yaml:"mysql"`
	*ElasticSearch `yaml:"elasticsearch"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
	URL      string `yaml:"url"`
}

type ElasticSearch struct {
	Addresses []string `yaml:"addresses"`
}
