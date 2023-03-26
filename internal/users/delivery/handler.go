package delivery

import (
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
	userUseCase "github.com/Inspirate789/SOK-golang-test-task/internal/users/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type delivery struct {
	userUC userUseCase.UseCase
}

func NewDelivery(router *gin.Engine, uuc userUseCase.UseCase) {
	handler := &delivery{
		userUC: uuc,
	}
	router.POST("/user", handler.createUser)
	router.DELETE("/user", handler.getBalance)
	router.GET("/balance", handler.deleteUser)
}

func (d *delivery) createUser(ctx *gin.Context) {
	name := ctx.Query("name")

	id, err := d.userUC.CreateUser(name)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (d *delivery) deleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, models.ErrUrlParamNotFound("id"))
		return
	}

	err = d.userUC.DeleteUser(int(id))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (d *delivery) getBalance(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, models.ErrUrlParamNotFound("id"))
		return
	}

	balance, err := d.userUC.GetUserBalance(int(id))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}
