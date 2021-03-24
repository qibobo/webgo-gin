package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/qibobo/webgo-gin/db"
	"github.com/qibobo/webgo-gin/models"
)

// DemoController example
type DemoController struct {
	logger *zap.Logger
	demoDB db.DemoDB
}

// NewController example
func NewDemoController(logger *zap.Logger, demoDB db.DemoDB) *DemoController {
	return &DemoController{
		logger: logger,
		demoDB: demoDB,
	}
}

// PingExample godoc
// @Summary get demo
// @Description get demo
// @Tags demo
// @Accept json
// @Produce json
// @Success 200 {string} json "pong"
// @Failure 400 {string} json "ok"
// @Failure 404 {string} json "ok"
// @Failure 500 {string} json "ok"
// @Router /demo/{id} [get]
func (c *DemoController) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.logger.Error("failed to get id", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "id is invalid",
		})
		return
	}
	demo, err := c.demoDB.GetDemo(id)
	if err != nil {
		c.logger.Error("failed to get demo by id", zap.Error(err), zap.Int("id", id))
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	if demo == nil {
		c.logger.Info("no result when get demo by id", zap.Int("id", id))
		ctx.JSON(http.StatusOK, models.Response{
			Code:    http.StatusNotFound,
			Message: "no result when get demo by id",
		})
		return
	}
	ctx.JSON(http.StatusOK, demo)
}

// PingExample godoc
// @Summary create demo
// @Description create demo
// @Tags demo
// @Accept json
// @Produce json
// @Success 200 {string} json "pong"
// @Failure 400 {string} json "ok"
// @Failure 404 {string} json "ok"
// @Failure 500 {string} json "ok"
// @Router /demo [post]
func (c *DemoController) CreateDemo(ctx *gin.Context) {
	demo := models.Demo{}
	err := ctx.BindJSON(&demo)
	if err != nil {
		c.logger.Error("failed to parse body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "request body is invalid",
		})
		return
	}
	err = c.demoDB.CreateDemo(&demo)
	if err != nil {
		c.logger.Error("failed to create demo", zap.Error(err), zap.Any("demo", demo))
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, demo)
}
