package config

import "time"

type Rule struct {
	Limit  int
	Window time.Duration
}

type Config struct {
	Routes map[string]Rule
	Default Rule
}
