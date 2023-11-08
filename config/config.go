package config

type Config struct {
  DB DBConfig
  JWT JwtConfig
}

type DBConfig struct {
  Url string `required:"true" split_words:"true"`
  TestUrl string `required:"true" split_words:"true"` 
}

type JwtConfig struct {
  AccessToken string `required:"true" split_words:"true"`
  RefreshToken string `required:"true" split_words:"true"`
}


