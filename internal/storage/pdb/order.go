package pdb

import (
	"database/sql"
	"fmt"

	"karrrakotov/wb_go_lvl0/internal/models"
	"karrrakotov/wb_go_lvl0/internal/storage"
)

type storageOrder struct {
	database *sql.DB
}

func (st *storageOrder) FindAll() (ordersDTO map[string]models.OrderDTO, err error) {
	ordersDTO = make(map[string]models.OrderDTO)

	// Поиск всех Order
	orders, err := st.FindAllOrder()
	if err != nil {
		return nil, err
	}

	// Поиск всех Delivery
	deliverys, err := st.FindAllDelivery()
	if err != nil {
		return nil, err
	}

	// Поиск всех Payment
	payments, err := st.FindAllPayment()
	if err != nil {
		return nil, err
	}

	// Поиск всех Item
	items, err := st.FindAllItem()
	if err != nil {
		return nil, err
	}

	// Заполнение мапы
	for _, val := range orders {
		ordersDTO[val.OrderUID] = models.OrderDTO{
			OrderUID:          val.OrderUID,
			TrackNumber:       val.TrackNumber,
			Entry:             val.Entry,
			Locale:            val.Locale,
			InternalSignature: val.InternalSignature,
			CustomerID:        val.CustomerID,
			DeliveryService:   val.DeliveryService,
			ShardKey:          val.ShardKey,
			SMID:              val.SMID,
			DateCreated:       val.DateCreated,
			OOFShard:          val.OOFShard,
		}
	}

	// Сохранение Delivery
	for _, val := range deliverys {
		dto := ordersDTO[val.OrderUID]
		dto.Delivery = models.DeliveryDTO{
			Name:    val.Name,
			Phone:   val.Phone,
			Zip:     val.Zip,
			City:    val.City,
			Address: val.Address,
			Region:  val.Region,
			Email:   val.Email,
		}
		ordersDTO[val.OrderUID] = dto
	}

	// Сохранение Payment
	for _, val := range payments {
		dto := ordersDTO[val.OrderUID]
		dto.Payment = models.PaymentDTO{
			Transaction:  val.Transaction,
			RequestID:    val.RequestID,
			Currency:     val.Currency,
			Provider:     val.Provider,
			Amount:       val.Amount,
			PaymentDT:    val.PaymentDT,
			Bank:         val.Bank,
			DeliveryCost: val.DeliveryCost,
			GoodsTotal:   val.GoodsTotal,
			CustomFee:    val.CustomFee,
		}
		ordersDTO[val.OrderUID] = dto
	}

	// Сохранение Item
	for _, val := range items {
		dto := ordersDTO[val.OrderUID]
		item := models.ItemDTO{
			ChrtID:      val.ChrtID,
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

		dto.Items = append(dto.Items, item)
		ordersDTO[val.OrderUID] = dto
	}

	return ordersDTO, nil
}

func (st *storageOrder) InsertInto(order models.Order, delivery models.Delivery, payment models.Payment, items []models.Item) error {
	// Начало транзакции
	tx, err := st.database.Begin()
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %v", err)
	}

	// Вставка данных в таблицу order
	if err := st.InsertIntoOrder(order); err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при сохранении данных в таблицу order: %v", err)
	}

	// Вставка данных в таблицу delivery
	if err := st.InsertIntoDelivery(delivery); err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при сохранении данных в таблицу delivery: %v", err)
	}

	// Вставка данных в таблицу payment
	if err := st.InsertIntoPayment(payment); err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при сохранении данных в таблицу payment: %v", err)
	}

	// Вставка данных в таблицу item
	for _, val := range items {
		if err := st.InsertIntoItem(val); err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка при сохранении данных в таблицу item: %v", err)
		}
	}

	// Делаем commmit
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при коммите транзакции: %v", err)
	}

	return nil
}

func (st *storageOrder) InsertIntoOrder(order models.Order) error {
	_, err := st.database.Exec("INSERT INTO Orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OOFShard)

	if err != nil {
		return err
	}
	return nil
}

func (st *storageOrder) InsertIntoDelivery(delivery models.Delivery) error {
	_, err := st.database.Exec("INSERT INTO Delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		delivery.OrderUID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)

	if err != nil {
		return err
	}
	return nil
}

func (st *storageOrder) InsertIntoPayment(payment models.Payment) error {
	_, err := st.database.Exec("INSERT INTO Payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		payment.OrderUID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDT, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)

	if err != nil {
		return err
	}
	return nil
}

func (st *storageOrder) InsertIntoItem(item models.Item) error {
	_, err := st.database.Exec("INSERT INTO Item (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		item.ChrtID, item.OrderUID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)

	if err != nil {
		return err
	}
	return nil
}

func (st *storageOrder) FindAllOrder() (orders []models.Order, err error) {
	// Составляем запрос
	rows, err := st.database.Query("SELECT * FROM Orders")
	if err != nil {
		return nil, fmt.Errorf("ошибка при составлении строк: %v", err)
	}
	defer rows.Close()

	// Перебираем каждую строку
	for rows.Next() {
		var order models.Order
		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк: %v", err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (st *storageOrder) FindAllDelivery() (deliverys []models.Delivery, err error) {
	// Составляем запрос
	rows, err := st.database.Query("SELECT * FROM Delivery")
	if err != nil {
		return nil, fmt.Errorf("ошибка при составлении строк: %v", err)
	}
	defer rows.Close()

	// Перебираем каждую строку
	for rows.Next() {
		var delivery models.Delivery
		err = rows.Scan(
			&delivery.OrderUID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк: %v", err)
		}

		deliverys = append(deliverys, delivery)
	}

	return deliverys, nil
}

func (st *storageOrder) FindAllPayment() (payments []models.Payment, err error) {
	// Составляем запрос
	rows, err := st.database.Query("SELECT * FROM Payment")
	if err != nil {
		return nil, fmt.Errorf("ошибка при составлении строк: %v", err)
	}
	defer rows.Close()

	// Перебираем каждую строку
	for rows.Next() {
		var payment models.Payment
		err = rows.Scan(
			&payment.OrderUID,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк: %v", err)
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (st *storageOrder) FindAllItem() (items []models.Item, err error) {
	// Составляем запрос
	rows, err := st.database.Query("SELECT * FROM Item")
	if err != nil {
		return nil, fmt.Errorf("ошибка при сканировании строк: %v", err)
	}
	defer rows.Close()

	// Перебираем каждую строку
	for rows.Next() {
		var item models.Item
		err = rows.Scan(
			&item.OrderUID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк: %v", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func NewStorageOrder(database *sql.DB) storage.StorageOrder {
	return &storageOrder{
		database: database,
	}
}
