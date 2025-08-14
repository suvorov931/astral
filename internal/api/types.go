package api

type HttpServer struct {
	Host string `env:"HTTP_HOST" env-required:"true"`
	Port int    `env:"HTTP_PORT" env-required:"true"`
}

type User struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}

type Meta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}
