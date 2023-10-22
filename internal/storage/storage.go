package storage

import "karrrakotov/wb_go_lvl0/internal/models"

type StorageOrder interface {
	InsertInto(order models.Order, delivery models.Delivery, payment models.Payment, items []models.Item) error
	FindAll() (ordersDTO map[string]models.OrderDTO, err error)
}
