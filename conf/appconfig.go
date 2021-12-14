package conf

type AppConfiguration struct {
	Name    string
	Help    string
	Version string
}

var AppConfig = AppConfiguration{
	Name:    "rtyq",
	Help:    "generate and query in-memory spatial rtree",
	Version: "0.2.0",
}
