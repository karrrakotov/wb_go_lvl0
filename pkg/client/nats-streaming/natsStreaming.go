package natsstreaming

import (
	"log"
	"time"

	"github.com/nats-io/stan.go"

	"karrrakotov/wb_go_lvl0/internal/config"
)

const (
	connectWait        = time.Second * 30
	pubAckWait         = time.Second * 30
	interval           = 10
	maxOut             = 5
	maxPubAcksInflight = 25
)

func NewNatsConnect(cfg *config.Config) (stan.Conn, error) {
	return stan.Connect(
		cfg.Nats.ClusterID,
		cfg.Nats.ClientID,
		stan.ConnectWait(connectWait),
		stan.PubAckWait(pubAckWait),
		stan.NatsURL(cfg.Nats.URL),
		stan.Pings(interval, maxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("ошибка соединения, причина: %v", reason)
		}),
		stan.MaxPubAcksInflight(maxPubAcksInflight),
	)
}
