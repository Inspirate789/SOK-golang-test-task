package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Inspirate789/SOK-golang-test-task/cmd/worker/processor"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
	transactionsRepo "github.com/Inspirate789/SOK-golang-test-task/internal/transactions/repository"
	transactionsUseCase "github.com/Inspirate789/SOK-golang-test-task/internal/transactions/usecase"
	UsersRepo "github.com/Inspirate789/SOK-golang-test-task/internal/users/repository"
	"github.com/jmoiron/sqlx"
	"github.com/namsral/flag"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	// Kafka
	kafkaBrokerUrls    string
	kafkaVerbose       bool
	kafkaTopic         string
	kafkaConsumerGroup string
	kafkaClientId      string

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
	flag.StringVar(&kafkaBrokerUrls, "kafka-brokers",
		"localhost:19092,localhost:29092,localhost:39092",
		"Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic. Only one topic per worker.")
	flag.StringVar(&kafkaConsumerGroup, "kafka-consumer-group", "consumer-group", "Kafka consumer group")
	flag.StringVar(&kafkaClientId, "kafka-client-id", "my-client-id", "Kafka client id")

	flag.StringVar(&pgHost, "postgres-host", "localhost", "")
	flag.StringVar(&pgPort, "postgres-port", "5432", "")
	flag.StringVar(&pgUsername, "postgres-username", "root", "")
	flag.StringVar(&pgPassword, "postgres-password", "postgres", "")
	flag.StringVar(&pgDbName, "postgres-dbname", "postgres_db", "")
	flag.StringVar(&pgSslMode, "postgres-ssl-mode", "disable", "")
	flag.StringVar(&pgDriverName, "postgres-driver-name", "postgres", "")

	flag.Parse()
}

func NewDB() (*sqlx.DB, error) {
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost,
		pgPort,
		pgUsername,
		pgPassword,
		pgDbName,
		pgSslMode)

	sqlDB, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()

	return sqlx.NewDb(sqlDB, pgDriverName), nil
}

func main() {
	logger := zerolog.Logger{}.Level(zerolog.InfoLevel)
	db, err := NewDB()
	if err != nil {
		logger.Fatal().Err(err)
	}
	defer db.Close()
	proc := processor.NewKafkaProcessor(&processor.KafkaProcessorConfig{
		BrokerUrls:      strings.Split(kafkaBrokerUrls, ","),
		ClientID:        kafkaClientId,
		NewConsumerFunc: consumer.NewKafkaTransactionConsumer,
		DetectionTime:   time.Second,
		UseCase:         transactionsUseCase.NewUseCase(UsersRepo.NewPgRepository(db), transactionsRepo.NewPgRepository(db)),
		Logger:          &logger,
	})

	procCtx := context.Background()
	go func() {
		err = proc.ProcessQueue(procCtx)
	}()
	logger.Debug().Msg("Processor started")

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug().Msg("Shutdown Processor ...")

	procCtx.Done()
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Processor Shutdown: %v", err))
	}
	logger.Debug().Msg("Processor exited")
}
