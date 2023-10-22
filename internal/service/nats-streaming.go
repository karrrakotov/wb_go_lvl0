package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/stan.go"

	"karrrakotov/wb_go_lvl0/internal/models"
	"karrrakotov/wb_go_lvl0/internal/storage"
)

type natsStreaming struct {
	storageOrder storage.StorageOrder
}

var in_memory = make(map[string]models.OrderDTO) // Инициализация карты

func (s *natsStreaming) Subscribe(sc stan.Conn) (err error) {
	// Подписка на канал
	_, err = sc.Subscribe("test", func(m *stan.Msg) {
		fmt.Printf("Полученные данные: %s\n", string(m.Data))

		var data models.OrderDTO
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			log.Println("не удалось прочитать данные json:", err)
			return
		}

		// Проверка данных на валидность
		if s.validateData(data) {
			// Получение структур данных
			order := models.Order{
				OrderUID:          data.OrderUID,
				TrackNumber:       data.TrackNumber,
				Entry:             data.Entry,
				Locale:            data.Locale,
				InternalSignature: data.InternalSignature,
				CustomerID:        data.CustomerID,
				DeliveryService:   data.DeliveryService,
				ShardKey:          data.ShardKey,
				SMID:              data.SMID,
				DateCreated:       data.DateCreated,
				OOFShard:          data.OOFShard,
			}

			delivery := models.Delivery{
				OrderUID: data.OrderUID,
				Name:     data.Delivery.Name,
				Phone:    data.Delivery.Phone,
				Zip:      data.Delivery.Zip,
				City:     data.Delivery.City,
				Address:  data.Delivery.Address,
				Region:   data.Delivery.Region,
				Email:    data.Delivery.Email,
			}

			payment := models.Payment{
				OrderUID:     data.OrderUID,
				Transaction:  data.Payment.Transaction,
				RequestID:    data.Payment.RequestID,
				Currency:     data.Payment.Currency,
				Provider:     data.Payment.Provider,
				Amount:       data.Payment.Amount,
				PaymentDT:    data.Payment.PaymentDT,
				Bank:         data.Payment.Bank,
				DeliveryCost: data.Payment.DeliveryCost,
				GoodsTotal:   data.Payment.GoodsTotal,
				CustomFee:    data.Payment.CustomFee,
			}

			var items []models.Item
			for _, val := range data.Items {
				item := models.Item{
					ChrtID:      val.ChrtID,
					OrderUID:    data.OrderUID,
					TrackNumber: val.TrackNumber,
					Price:       val.Price,
					RID:         val.RID,
					Name:        val.Name,
					Sale:        val.Sale,
					Size:        val.Size,
					TotalPrice:  val.TotalPrice,
					NmID:        val.NmID,
					Brand:       val.Brand,
					Status:      val.Status,
				}

				items = append(items, item)
			}

			// Сохранение структур в базу данных
			if err := s.storageOrder.InsertInto(order, delivery, payment, items); err != nil {
				log.Fatalln("ошибка при вставке данных в таблицы БД")
			}

			// Сохранение данных в память сервера
			in_memory[order.OrderUID] = data
		} else {
			log.Println("входные данные не соответствуют спецификации")
		}
	})

	if err != nil {
		log.Fatalf("ошибка при подписке на канал: %v", err)
	}

	return nil
}

func (s *natsStreaming) validateData(data models.OrderDTO) (isOk bool) {
	// Проверка соответствия каждого поля
	if data.OrderUID != "" &&
		data.TrackNumber != "" &&
		data.Entry != "" &&
		data.Delivery.Name != "" &&
		data.Delivery.Phone != "" &&
		data.Delivery.Zip != "" &&
		data.Delivery.City != "" &&
		data.Delivery.Address != "" &&
		data.Delivery.Region != "" &&
		data.Delivery.Email != "" &&
		data.Payment.Transaction != "" &&
		data.Payment.Currency != "" &&
		data.Payment.Provider != "" &&
		data.Payment.PaymentDT >= 0 &&
		data.Payment.Bank != "" &&
		len(data.Items) > 0 &&
		data.Locale != "" &&
		data.CustomerID != "" &&
		data.SMID >= 0 &&
		data.DeliveryService != "" &&
		data.ShardKey != "" &&
		data.DateCreated != "" &&
		data.OOFShard != "" {

		// Проверка входного массива Items
		for _, val := range data.Items {
			if val.TrackNumber != "" &&
				val.ChrtID >= 0 &&
				val.Price >= 0 &&
				val.TotalPrice >= 0 &&
				val.NmID >= 0 &&
				val.Status >= 0 &&
				val.RID != "" &&
				val.Name != "" &&
				val.Size != "" &&
				val.Brand != "" {
			} else {
				return
			}
		}

		return true
	}

	return
}

func NewServiceNatsStreaming(storageOrder storage.StorageOrder) ServiceNatsStreaming {
	return &natsStreaming{
		storageOrder: storageOrder,
	}
}
