package services

type Config interface {
	HttpHost() string
	HttpPort() string
	PostgresHost() string
	PostgresPort() string
	PostgresUser() string
	PostgresPassword() string
	PostgresName() string
	EnvLevel() string
	SslMode() string
}
