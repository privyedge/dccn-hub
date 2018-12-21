package  config

// Config broker & mongo
type Config struct {
	Broker Address
	Mongo Address
	Register Address
}

type Address struct {
	Host string `json:"host,omitempty"`
	Port uint16 `json:"port,omitempty"`
}