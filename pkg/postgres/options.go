package postgres

import "time"

type Option func(*Postgres)

// MaxPoolSize максимальний розмір пулу з'єднань для Postgres.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts кількість спроб підключення до Postgres.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout таймаут для підключення до Postgres.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}
