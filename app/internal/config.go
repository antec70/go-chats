package internal

const DefaultConfig string = "./config/params-local.json"

type ParamsLocal struct {
	Db struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Serv     string `json:"serv"`
		Table    string `json:"table"`
	}
}
