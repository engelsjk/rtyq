package config

// Lookup ...
type Lookup struct {
	Data struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
	} `json:"data"`
	Database struct {
		Path      string `json:"path"`
		Extension string `json:"index"`
	} `json:"database"`
}

// Config ...
type Config struct {
	Port    int `json:"port"`
	Lookups []Lookup
}

func Load(path string) {

}
