package main

import (
	"whisper/init"
	"whisper/internal/controller"
	"whisper/pkg/context"
	"whisper/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.Trace())

	lol := router.Group("/lol")
	{
		lol.POST("/equipment", context.Handle(controller.Equipment))
		lol.POST("/heroes", context.Handle(controller.Heroes))
		lol.POST("/rune", context.Handle(controller.Rune))
		lol.POST("/skill", context.Handle(controller.Skill))
	}

	run.Init()
	router.Run("0.0.0.0:8123") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
