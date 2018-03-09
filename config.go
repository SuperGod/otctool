package main

type Config struct {
	Influx struct {
		Addr string
		User string
		Pwd  string
		DB   string `json:"db"`
	}
	Sleep int64 // sleep seconds
}
