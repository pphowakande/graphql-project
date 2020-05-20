package configuration

// Configuration used for the main app
type Configuration struct {
	DbConnections map[string]DBConfig
	GrpcURL       string
	BMSEndpoint   string
}


// DBConfig is the config for database connections
type DBConfig struct {
	DBUser string
	DBPassword string
	DBHost string
	DBName string
}