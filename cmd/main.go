package main

import (
	context2 "context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"whisper/init"
	"whisper/internal/controller"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.Trace())
	router.Delims("{[{", "}]}")
	router.LoadHTMLGlob("web/template/*")

	page := router.Group("/")
	{
		page.Static("css", "web/static/css")
		page.Static("js", "web/static/js")
		page.Static("font", "web/static/font")
		page.StaticFile("favicon.ico", "web/static/favicon.ico")
		page.GET("/", context.Handle(controller.SearchBox))
	}

	lol := router.Group("/lol")
	{
		lol.POST("/equipment", context.Handle(controller.Equipment))
		lol.POST("/heroes", context.Handle(controller.Heroes))
		lol.POST("/heroes/attr", context.Handle(controller.HeroesAttribute))
		lol.POST("/rune", context.Handle(controller.Rune))
		lol.POST("/skill", context.Handle(controller.Skill))
	}

	es := router.Group("/")
	{
		es.POST("/query", context.Handle(controller.Query))
		es.POST("/index/build", context.Handle(controller.Build))
	}

	db := router.Group("/")
	{
		db.POST("/alias/heroes", context.Handle(controller.AliasHeroes))
		db.POST("/alias/equip", context.Handle(controller.AliasEquip))
	}
	run.Init()

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", config.GlobalConfig.App.IP, config.GlobalConfig.App.Port),
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Logger.Warnln("Shutdown Server ...")
	ctx, cancel := context2.WithTimeout(context2.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Panic(err)
	}
	log.Logger.Infoln("Server exiting")
}
