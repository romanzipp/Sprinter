package models

import (
	"github.com/go-ping/ping"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Check struct {
	gorm.Model
	Addr       string
	IP         string
	PacketLoss float64
	Latency    time.Duration
	Success    bool
}

func ExecPingCheck(host string, db *gorm.DB) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		log.Error().Msgf("error creating pinger for host %s", host)
		return
	}

	pinger.Count = 1
	pinger.Timeout = 5 * time.Second
	pinger.Debug = true
	err = pinger.Run()
	if err != nil {
		log.Error().Msgf("error creating pinger for host %s", host)
		return
	}

	stats := pinger.Statistics()

	log.Debug().Msgf("ping {host: %s, ip: %s, packet-loss: %.2f, rtt: %d}", stats.Addr, stats.IPAddr, stats.PacketLoss, stats.AvgRtt)

	check := Check{
		Addr:       stats.Addr,
		IP:         stats.IPAddr.String(),
		PacketLoss: stats.PacketLoss,
		Latency:    stats.AvgRtt,
		Success:    stats.PacketLoss == float64(0),
	}

	db.Save(&check)
}
