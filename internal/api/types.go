package api

type HttpServer struct {
	Host string `env:"HTTP_HOST"`
	Port int    `env:"HTTP_PORT"`
}

type User struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}
