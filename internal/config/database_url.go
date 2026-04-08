package config

import (
	"net/url"
)

// PostgresURL returns a URL suitable for github.com/golang-migrate/migrate and lib/pq style drivers.
func (c *DatabaseConfig) PostgresURL() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Host + ":" + c.Port,
		Path:   "/" + c.Name,
	}
	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
