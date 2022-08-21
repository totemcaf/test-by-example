package model

type Response struct {
	StatusCode int `yaml:"statusCode"`
	Body       *Json
}
