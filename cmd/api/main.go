package main

import (
	"fmt"
	"github.com/Inspirate789/SOK-golang-test-task/cmd/api/middleware"
	"github.com/Inspirate789/SOK-golang-test-task/cmd/api/server"
	transactionsDelivery "github.com/Inspirate789/SOK-golang-test-task/internal/transactions/delivery"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/producer"
	usersDelivery "github.com/Inspirate789/SOK-golang-test-task/internal/users/delivery"
	"github.com/Inspirate789/SOK-golang-test-task/internal/users/repository"
	usersUseCase "github.com/Inspirate789/SOK-golang-test-task/internal/users/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	// Web server
	listenAddrApi string

	// Kafka
	kafkaBrokerUrls string
	kafkaVerbose    bool
	kafkaClientId   string
	kafkaTopic      string

	// Postgres
	pgHost       string
	pgPort       string
	pgUsername   string
	pgPassword   string
	pgDbName     string
	pgSslMode    string
	pgDriverName string
)

func init() {
	flag.StringVar(&listenAddrApi, "listen-address", "0.0.0.0:9000", "Listen address for api")

	flag.StringVar(&kafkaBrokerUrls, "kafka-brokers",
		"localhost:19092,localhost:29092,localhost:39092",
		"Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaClientId, "kafka-client-id", "my-kafka-client", "Kafka client id to connect")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic to push")

	flag.StringVar(&pgHost, "postgres-host", "localhost", "")
	flag.StringVar(&pgPort, "postgres-port", "5432", "")
	flag.StringVar(&pgUsername, "postgres-username", "root", "")
	flag.StringVar(&pgPassword, "postgres-password", "postgres", "")
	flag.StringVar(&pgDbName, "postgres-dbname", "postgres_db", "")
	flag.StringVar(&pgSslMode, "postgres-ssl-mode", "disable", "")
	flag.StringVar(&pgDriverName, "postgres-driver-name", "postgres", "")

	flag.Parse()
}

func NewGinRouter(db *sqlx.DB, logger *zerolog.Logger) *gin.Engine {
	router := gin.Default()
	router.UseRawPath = true
	router.UnescapePathValues = false
	router.Use(middleware.QueryLogger(logger))

	usersDelivery.NewDelivery(router, usersUseCase.NewUseCase(repository.NewPgRepository(db)))
	transactionsDelivery.NewDelivery(router, producer.NewKafkaTransactionProducer(&producer.KafkaProducerConfig{
		BrokerUrls: strings.Split(kafkaBrokerUrls, ","),
		ClientID:   kafkaClientId,
	}))

	return router
}

func ExitServer(logger *zerolog.Logger, srv *server.Server) {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug().Msg("Shutdown Server ...")

	err := srv.Stop()
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Server Shutdown: %v", err))
	}
	logger.Debug().Msg("Server exited")
}

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost,
		pgPort,
		pgUsername,
		pgPassword,
		pgDbName,
		pgSslMode,
	))
	if err != nil {
		logger.Fatal().Err(err)
	}
	defer db.Close()

	srv := server.NewServer(listenAddrApi, NewGinRouter(db, &logger))
	go func() {
		err = srv.Start()
		if err != nil && err != http.ErrServerClosed {
			logger.Error().Msg(fmt.Sprintf("listen: %s\n", err))
		}
	}()
	logger.Debug().Msg("Server started")

	ExitServer(&logger, srv)
}
