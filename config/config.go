package config

type Config struct {
  DB DBConfig
}

type DBConfig struct {
  Url string `required:"true" split_words:"true"`
  TestUrl string `required:"true" split_words:"true"` 
}
