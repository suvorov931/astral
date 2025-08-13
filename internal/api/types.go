package api

type HttpServer struct {
	Host string `env:"HTTP_HOST" env-required:"true"`
	Port int    `env:"HTTP_PORT" env-required:"true"`
}

type User struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}
