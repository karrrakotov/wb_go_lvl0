package app

import (
	"log"
	"net/http"
	"os"

	"github.com/nats-io/stan.go"

	"karrrakotov/wb_go_lvl0/internal/config"
	"karrrakotov/wb_go_lvl0/internal/server"
	"karrrakotov/wb_go_lvl0/internal/service"
	"karrrakotov/wb_go_lvl0/internal/storage/pdb"
	"karrrakotov/wb_go_lvl0/internal/transport/rest"
	natsstreaming "karrrakotov/wb_go_lvl0/pkg/client/nats-streaming"
	postgresql "karrrakotov/wb_go_lvl0/pkg/client/postgreSQL"
)

func Run() {
	// Init config
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Init Router&Server
	router := http.NewServeMux()
	server := new(server.Server)

	// Connect to PostgreSQL
	postgreClient, err := postgresql.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка при подключении базы данных: %v", err)
		return
	}

	// Connect to Nats-Streaming
	sc, err := natsstreaming.NewNatsConnect(cfg)
	if err != nil {
		log.Fatalf("не удалось подключиться к nats-streaming: %v", err)
		return
	}

	defer sc.Close()

	// Init Storage
	storageOrder := pdb.NewStorageOrder(postgreClient)

	// Init Services
	serviceOrder := service.NewServiceOrder(storageOrder)
	serviceNats := service.NewServiceNatsStreaming(storageOrder)

	// Call Services
	if err := serviceOrder.RecoveryInMemory(); err != nil {
		log.Fatalf("не удалось восстановить данные в кэш: %v", err)
	}
	serviceNats.Subscribe(sc)

	// Init Handlers
	handlerOrder := rest.NewHandlerOrder(serviceOrder)

	// Call Handlers
	handlerOrder.Init(router)

	// Call Script for test NATS-STREAMING
	PublishMessage()

	// Run Server
	if err := server.Run("1234", router); err != nil {
		log.Fatalln("ошибка при запуске сервера")
		return
	}
}

// Проверка работы канала
func PublishMessage() {
	sc, err := stan.Connect("microservice", "microservice_b", stan.NatsURL("nats://nats-streaming:4222"))
	if err != nil {
		log.Fatalf("ошибка подключения nats-streaming: %v", err)
	}
	defer sc.Close()

	// Чтение данных из JSON файла
	filePath := "./data.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	// Публикация сообщения
	err = sc.Publish("test", jsonData)
	if err != nil {
		log.Printf("ошибка при публикации сообщения: %v", err)
	} else {
		log.Println("сообщение успешно опубликовано")
	}
}
