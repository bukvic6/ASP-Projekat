package configstore

type Config struct {
	Id      string            `json:"id"`
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}

type Group struct {
	Version string   `json:"version"`
	Id      string   `json:"id"`
	Config  []Config `json:"config"`
}
