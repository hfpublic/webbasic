package configs

type Server struct {
	TokenTime       int  `yaml:"tokentimes"`
	DemoMode        bool `yaml:"demomode"`
	AuthorizeSwitch bool `yaml:"authorize"`
	*HTTPServer     `yaml:"httpserver"`
}

type HTTPServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
