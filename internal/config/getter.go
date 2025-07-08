package config

// GetServerPort returns the server port
func (c *Config) GetServerPort() string {
	return c.Server.Port
}

// GetGinMode returns the Gin mode
func (c *Config) GetGinMode() string {
	return c.Server.GinMode
}

// GetShutdownTimeout returns the server shutdown timeout in seconds
func (c *Config) GetShutdownTimeout() int {
	return c.Server.ShutdownTimeout
}

// GetDBHost returns the database host
func (c *Config) GetDBHost() string {
	return c.Database.Host
}

// GetDBPort returns the database port
func (c *Config) GetDBPort() string {
	return c.Database.Port
}

// GetDBUser returns the database user
func (c *Config) GetDBUser() string {
	return c.Database.User
}

// GetDBPassword returns the database password
func (c *Config) GetDBPassword() string {
	return c.Database.Password
}

// GetDBName returns the database name
func (c *Config) GetDBName() string {
	return c.Database.Name
}

// GetDBSSLMode returns the database SSL mode
func (c *Config) GetDBSSLMode() string {
	return c.Database.SSLMode
}

// GetRedisHost returns the Redis host
func (c *Config) GetRedisHost() string {
	return c.Redis.Host
}

// GetRedisPort returns the Redis port
func (c *Config) GetRedisPort() string {
	return c.Redis.Port
}

// GetRedisPassword returns the Redis password
func (c *Config) GetRedisPassword() string {
	return c.Redis.Password
}

// GetRedisDB returns the Redis database number
func (c *Config) GetRedisDB() int {
	return c.Redis.DB
}

// GetDBConnectionString returns the database connection string
func (c *Config) GetDBConnectionString() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.Name +
		" sslmode=" + c.Database.SSLMode
}

// GetRedisAddress returns the Redis address
func (c *Config) GetRedisAddress() string {
	return c.Redis.Host + ":" + c.Redis.Port
}
