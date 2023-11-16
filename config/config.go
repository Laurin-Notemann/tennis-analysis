package config

type Config struct {
	DB   DBConfig
	JWT  JwtConfig
	ECHO EchoConfig
}

type DBConfig struct {
	Url     string `required:"true" split_words:"true"`
	TestUrl string `required:"false" split_words:"true"`
}

type JwtConfig struct {
	AccessToken  string `required:"true" split_words:"true"`
	RefreshToken string `required:"true" split_words:"true"`
}

type EchoConfig struct {
	Port int `required:"true" split_words:"true"`
  Host string `required:"true" split_words:"true"`
}
