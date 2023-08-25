package main

import (
	context2 "context"
	"fmt"
	"github.com/robfig/cron/v3"
	"net/http"
	"os"
	"os/signal"
	"time"
	run "whisper/init"
	"whisper/internal/controller"
	"whisper/internal/logic"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	pprof.Register(router, "dev/pprof")
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.Cors())
	router.Use(middleware.Trace())
	router.LoadHTMLGlob("web/template/whisper/dist/*.html")

	page := router.Group("/")
	{
		page.Static("assets", "web/template/whisper/dist/assets/")
		page.StaticFile("favicon.ico", "web/static/favicon.ico")
		page.GET("/", context.Handle(controller.SearchBox))
		page.GET("/equip", context.Handle(controller.SearchBox))

		page.GET("/version", context.Handle(controller.QueryVersion))
		page.GET("/equip/types", context.Handle(controller.QueryEquipTypes))
	}

	lol := router.Group("/lol")
	{
		lol.POST("/equipment", context.Handle(controller.Equipment))
		lol.POST("/heroes", context.Handle(controller.Heroes))
		lol.POST("/heroes/attr", context.Handle(controller.HeroesAttribute))
		lol.POST("/rune", context.Handle(controller.Rune))
		lol.POST("/rune/type", context.Handle(controller.RuneType))
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

	inner := router.Group("/")
	{
		inner.POST("/cron", context.Handle(controller.Cron))
		inner.POST("/equip/extract", context.Handle(controller.EquipExtract))
		inner.POST("/equip/filter", context.Handle(controller.EquipFilter))
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

	// 启动定时任务
	c := cron.New()
	_, err := c.AddFunc(config.LOLConfig.Cron.Time, func() {
		fmt.Println(time.Now())
		logic.Cron(nil)
	})
	if err != nil {
		panic(err)
	}
	c.Start()

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
