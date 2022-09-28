package config

type Config struct {
	Interval       int64
	PingHosts      []string
	PingTimeout    int64
	PingPrivileged bool
}
