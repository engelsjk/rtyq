package conf

var versionNumber string

type AppConfiguration struct {
	Name    string
	Help    string
	Version string
}

var AppConfig = AppConfiguration{
	Name:    "rtyq",
	Help:    "generate and query spatial rtrees on disk",
	Version: versionNumber,
}
