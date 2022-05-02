package main

type Config struct {
	Entries map[string]string `json:"entries"`
}

type Group struct {
	Id     string   `json:"id"`
	Config []Config `json:"config"`

	//	Configs Config `json:"configs"`
}
