package config

type Config struct {
	Interval       int64
	PingHosts      []string
	PingPrivileged bool
}
