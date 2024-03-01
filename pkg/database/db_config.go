package database

// DbConfig Модель данных для описания соединения с БД
type DbConfig struct {
	Driver               string
	Host                 string
	User                 string
	Password             string
	Db                   string
	Port                 string
	SslMode              bool
	MaxOpenConnections   int
	MaxIdleConnections   int
	Logging              bool
	DisableAutomaticPing bool
}
