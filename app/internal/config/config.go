package config

const DefaultConfig string = "./config/params-local.json"

type ParamsLocal struct {
	Port      string `json:"port"`
	Signature string `json:"signature"`
	CountConn int    `json:"count_conn"`
	Db        struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Serv     string `json:"serv"`
		Table    string `json:"table"`
	}
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		Db       string `json:"db"`
	}
	Pusher struct {
		Url    string `json:"url"`
		ApiKey string `json:"apiKey"`
	}
}
