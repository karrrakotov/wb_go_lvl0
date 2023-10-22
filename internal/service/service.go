package service

import (
	"github.com/nats-io/stan.go"

	"karrrakotov/wb_go_lvl0/internal/models"
)

type ServiceOrder interface {
	FindOne(order_uid string) (responseDTO models.OrderDTO, err error)
	RecoveryInMemory() error
}

type ServiceNatsStreaming interface {
	Subscribe(stan stan.Conn) (err error)
}
