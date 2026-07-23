package config

type Credentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	IMEI     string `json:"imei"`
	Password string `json:"password"`
}
