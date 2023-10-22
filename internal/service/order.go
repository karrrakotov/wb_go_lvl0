package service

import (
	"fmt"
	"reflect"

	"karrrakotov/wb_go_lvl0/internal/models"
	"karrrakotov/wb_go_lvl0/internal/storage"
)

type serviceOrder struct {
	storageOrder storage.StorageOrder
}

func (s *serviceOrder) FindOne(order_uid string) (responseDTO models.OrderDTO, err error) {
	// Обращаемся к im_memory и возвращаем нужный заказ по order_uid
	responseDTO = in_memory[order_uid]

	if reflect.DeepEqual(responseDTO, models.OrderDTO{}) {
		return responseDTO, fmt.Errorf("нет заказа с order_uid = %v", order_uid)
	}
	return responseDTO, nil
}

func (s *serviceOrder) RecoveryInMemory() error {
	// Находим все заказы
	ordersDTO, err := s.storageOrder.FindAll()
	if err != nil {
		return err
	}

	// Сохраняем их в оперативную память
	in_memory = ordersDTO
	return nil
}

func NewServiceOrder(storageOrder storage.StorageOrder) ServiceOrder {
	return &serviceOrder{
		storageOrder: storageOrder,
	}
}
