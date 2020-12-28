package config

const DefaultConfig string = "./config/params-local.json"

type ParamsLocal struct {
	Port string `json:"port"`
	Db   struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Serv     string `json:"serv"`
		Table    string `json:"table"`
	}
}
