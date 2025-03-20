package db

type AppConfig struct {
	Dbtype        string
	DbUri         string
	DbMaxPoolSize int
	DbMinPoolSize int
	DbName        string
}
