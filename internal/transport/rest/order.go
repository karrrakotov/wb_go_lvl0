package rest

import (
	"net/http"

	"karrrakotov/wb_go_lvl0/internal/service"
	"karrrakotov/wb_go_lvl0/internal/transport"
)

type handlerOrder struct {
	serviceOrder service.ServiceOrder
}

func (h *handlerOrder) Init(router *http.ServeMux) {
	// TODO - GET
	router.HandleFunc("/findOrder", h.findOrder)
}

// TODO GET - /findOrder
func (h *handlerOrder) findOrder(w http.ResponseWriter, r *http.Request) {
	// Проверка входящего запроса
	if r.Method != http.MethodGet {
		responseError := ResponseError{
			Status:  405,
			IsOk:    false,
			Message: "Метод не разрешен",
		}
		ResponseJson(w, 405, responseError)
		return
	}

	// Получение order_uid с входящего запроса
	order_uid := r.URL.Query().Get("order_uid")

	// Проверка order_uid с входящего запроса
	if order_uid == "" {
		responseError := ResponseError{
			Status:  400,
			IsOk:    false,
			Message: "Неверный order_uid: id заказа не должен быть пустым",
		}
		ResponseJson(w, 400, responseError)
		return
	}

	// Обращение к логике-сервера - поиск данных в БД
	responseDTO, err := h.serviceOrder.FindOne(order_uid)
	if err != nil {
		responseError := ResponseError{
			Status:  404,
			IsOk:    true,
			Message: "Ответ: " + err.Error(),
		}
		ResponseJson(w, 404, responseError)
		return
	}

	// Ответ
	responseOk := ResponseOk{
		Data:    responseDTO,
		Status:  200,
		IsOk:    true,
		Message: "Success!",
	}
	ResponseJson(w, 200, responseOk)
}

func NewHandlerOrder(serviceOrder service.ServiceOrder) transport.HandlerOrder {
	return &handlerOrder{
		serviceOrder: serviceOrder,
	}
}
