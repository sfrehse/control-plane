package api

import (
	"context"
	"control-plane/controller"
	"control-plane/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type V1Router struct {
	ctrl    *controller.Controller
	limiter *redis_rate.Limiter

	perMinute int
}

func NewRouter(controller *controller.Controller) *V1Router {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	limiter := redis_rate.NewLimiter(client)

	return &V1Router{
		ctrl:      controller,
		limiter:   limiter,
		perMinute: 20,
	}
}

func (router *V1Router) Router(v1 *gin.RouterGroup) {
	v1.POST("/generationTask", router.postGenerationTask())
	v1.GET("/generationTask/:id/status", router.getGenerationTaskStatus())
}

func (router *V1Router) postGenerationTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var generationTask models.GenerationTask
		if err := c.BindJSON(&generationTask); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if !router.isAllowedToSchedule() {
			c.Status(http.StatusTooManyRequests)
			return
		}

		status, err := router.ctrl.CreateNewTask(context.Background(), generationTask)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusCreated, status)
	}
}

func (router *V1Router) isAllowedToSchedule() bool {
	res, _ := router.limiter.Allow(context.Background(), "generation_task_limiter", redis_rate.PerMinute(router.perMinute))
	if res.Allowed == 0 {
		log.Warn("Exceed limit to schedule generation tasks.")
		return false
	}

	return true
}

func (router *V1Router) getGenerationTaskStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if !controller.ValidGenerationTaskId(id) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("not a valid id"))
			return
		}

		status, err := router.ctrl.GetGenerationTaskStatus(context.Background(), id)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, status)
	}
}
