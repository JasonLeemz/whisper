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
		page.POST("/query", context.Handle(controller.Query))
		page.POST("/equip/filter", context.Handle(controller.EquipFilter))

		page.GET("/version", context.Handle(controller.QueryVersion))
		page.GET("/equip/types", context.Handle(controller.QueryEquipTypes))
		page.GET("/hotkey", context.Handle(controller.GetHotKey))

		page.POST("/equip/roadmap", context.Handle(controller.GetRoadmap))
		page.POST("/equip/suit", context.Handle(controller.SuitEquip))
		page.POST("/hero/suit", context.Handle(controller.GetHeroSuit))

	}

	inner := router.Group("/")
	{
		inner.POST("/cron", context.Handle(controller.Cron))
		inner.POST("/equip/extract", context.Handle(controller.EquipExtract))
		inner.POST("/equipment", context.Handle(controller.Equipment))
		inner.POST("/heroes", context.Handle(controller.Heroes))
		inner.POST("/heroes/attr", context.Handle(controller.HeroesAttribute))
		inner.POST("/heroes/position", context.Handle(controller.HeroesPosition))
		inner.POST("/rune", context.Handle(controller.Rune))
		inner.POST("/rune/type", context.Handle(controller.RuneType))
		inner.POST("/skill", context.Handle(controller.Skill))
		inner.POST("/index/build", context.Handle(controller.Build))
		inner.POST("/alias/heroes", context.Handle(controller.AliasHeroes))
		inner.POST("/alias/equip", context.Handle(controller.AliasEquip))

		inner.POST("/equip/suit/batch", context.Handle(controller.BatchUpdateSuitEquip))
		inner.POST("/equip/suit/cache", context.Handle(controller.SuitData2Redis))
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
