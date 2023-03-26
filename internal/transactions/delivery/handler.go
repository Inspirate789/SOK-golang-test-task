package delivery

import (
	"context"
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/producer"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type delivery struct {
	transactionProducer producer.TransactionProducer
}

func NewDelivery(router *gin.Engine, tp producer.TransactionProducer) {
	handler := &delivery{
		transactionProducer: tp,
	}
	router.POST("/transaction", handler.performTransaction)
}

func (d *delivery) performTransaction(ctx *gin.Context) {
	var transaction consumer.TransactionDTO
	err := ctx.BindJSON(&transaction)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, models.ErrCannotParseJSON("TransactionDTO"))
		return
	}

	err = d.transactionProducer.Produce(context.Background(), strconv.Itoa(transaction.UserId), transaction)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
