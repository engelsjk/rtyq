package conf

var versionNumber string

type AppConfiguration struct {
	Name    string
	Help    string
	Version string
}

var AppConfig = AppConfiguration{
	Name:    "rtyq",
	Help:    "generate and query an in-memory spatial rtree",
	Version: versionNumber,
}
