package main

import (
	context2 "context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
	run "whisper/init"
	"whisper/internal/controller"
	"whisper/internal/logic"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/ip"
	"whisper/pkg/log"
	"whisper/pkg/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/grafana/pyroscope-go"
	"github.com/robfig/cron/v3"
)

func main() {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: "whisper",

		// replace this with the address of pyroscope server
		ServerAddress: "http://192.168.31.91:4040",

		// you can disable logging by setting this to nil
		Logger: nil,

		// you can provide static tags via a map:
		Tags: map[string]string{
			"ip": ip.GetLocalIP(),
		},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

	router := gin.New()
	pprof.Register(router, "dev/pprof")
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	router.Use(middleware.Trace())
	router.Use(middleware.Proc(), middleware.Params(), middleware.Auth())
	router.LoadHTMLGlob("web/template/whisper/dist/*.html")

	page := router.Group("/")
	{
		page.Static("assets", "web/template/whisper/dist/assets/")
		page.StaticFile("favicon.ico", "web/static/favicon.ico")

		// 自动填充
		page.POST("/auto/complete", context.Handle(controller.AutoComplete))

		page.GET("/", context.Handle(controller.SearchBox))
		page.POST("/hero/skins", context.Handle(controller.GetHeroSkins))

		page.GET("/search", context.Handle(controller.SearchBox))
		page.POST("/query", context.Handle(controller.Query))

		page.GET("/equip", context.Handle(controller.SearchBox))
		page.POST("/equip/filter", context.Handle(controller.EquipFilter))

		page.POST("/version/list", context.Handle(controller.VersionList))
		page.POST("/version/detail", context.Handle(controller.VersionDetail))

		page.GET("/version", context.Handle(controller.QueryVersion))
		page.POST("/version", context.Handle(controller.QueryVersion))
		page.GET("/equip/types", context.Handle(controller.QueryEquipTypes))
		page.GET("/hotkey", context.Handle(controller.GetHotKey))

		// 查询当前装备的合成路线和可合成路线
		page.POST("/equip/roadmap", context.Handle(controller.GetRoadmap))
		// 通过heroID来查询适配的装备，这个接口是将format数据写入db，接口返回的是第三方数据源的数据
		page.POST("/equip/suit", context.Handle(controller.SuitEquip))
		// 页面上通过heroID来查询适配的装备
		page.POST("/hero/suit", context.Handle(controller.GetHeroSuit))

		// 通过装备id查看适配英雄
		page.POST("/equip/hero/suit", context.Handle(controller.GetEquipHeroSuit))
		// 通过符文id查看适配英雄
		page.POST("/rune/hero/suit", context.Handle(controller.GetRuneHeroSuit))
		// 通过召唤师技能id查看适配英雄
		page.POST("/skill/hero/suit", context.Handle(controller.GetSkillHeroSuit))

	}

	inner := router.Group("/")
	{
		inner.POST("/cron", context.Handle(controller.Cron))
		inner.POST("/equip/extract", context.Handle(controller.EquipExtract))
		inner.POST("/equipment", context.Handle(controller.Equipment))
		inner.POST("/heroes", context.Handle(controller.Heroes))
		inner.POST("/heroes/attr", context.Handle(controller.HeroesAttribute))
		inner.POST("/rune", context.Handle(controller.Rune))
		inner.POST("/rune/type", context.Handle(controller.RuneType))
		inner.POST("/skill", context.Handle(controller.Skill))
		inner.POST("/index/build", context.Handle(controller.Build))
		inner.POST("/alias/heroes", context.Handle(controller.AliasHeroes))
		inner.POST("/alias/equip", context.Handle(controller.AliasEquip))

		{
			// 英雄的适配装备

			// 写DB
			// 1. LOLM将英雄适合的位置写入heroes_position（批量执行）
			// 2. LOL将英雄适合的位置写入heroes_position
			// 3. LOL将英雄适合的装备写入heroes_suit
			// 4. LOLM将英雄适合的装备写入heroes_suit
			inner.POST("/equip/suit/batch", context.Handle(controller.BatchUpdateSuitEquip))

			// 写Cache
			// 页面查询英雄合适的装备是从redis中获取的,要提前执行这个才能拿到数据
			// 这个接口依赖/equip/suit/batch，需要先执行/equip/suit/batch
			inner.POST("/suit/hero/cache", context.Handle(controller.SuitData2Redis))
		}

		// 装备、符文、技能适配英雄列表，汇总db然后写入redis
		// 页面查询：1.根据装备id获取适配的英雄
		inner.POST("/suit/high_rate/cache", context.Handle(controller.SuitHeroData2Redis))

		// 缓存heroes的attribute
		//inner.POST("/attr/hero/cache", context.Handle(controller.AttrData2Redis))

		{
			// 抓取外站攻略
			inner.POST("/strategy/grab", context.Handle(controller.GrabStrategy))

			// 查询英雄攻略
			inner.POST("/strategy/hero", context.Handle(controller.StrategyHero))
			// 查询装备攻略
			inner.POST("/strategy/equip", context.Handle(controller.StrategyEquip))
			// 查询符文攻略
			inner.POST("/strategy/rune", context.Handle(controller.StrategyRune))
		}

		{
			inner.POST("/equipment/setbit", context.Handle(controller.SetBit))
			inner.POST("/equipment/contain", context.Handle(controller.Contain))
		}
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
	defer func() {
		log.Logger.Sync()
		log.RpcLogger.Sync()
		log.GLogger.Sync()
		log.ELogger.Sync()
		log.MLogger.Sync()
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Panic(err)
	}
	log.Logger.Infoln("Server exiting")
}
