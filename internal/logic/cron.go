package logic

import (
	context2 "context"
	"sync"
	"time"
	"whisper/internal/logic/common"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

func Cron(ctx *context.Context) {
	if ctx == nil {
		ctx = context.NewContext()
	}
	_, cancelFunc := context2.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	wg := &sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer func() {
			wg.Done()
		}()

		_, err := QueryEquipments(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
		_, err = QueryEquipments(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
		}()
		_, err := QueryHeroes(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
		_, err = QueryHeroes(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		_, err = HeroAttribute(ctx, "0", common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
		_, err = HeroAttribute(ctx, "0", common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
		}()

		_, err := QueryRune(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
		_, err = QueryRune(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
		}()

		_, err := QuerySkill(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
		_, err = QuerySkill(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	wg.Wait()
	log.Logger.Info(ctx, "start building index...")
	err := BuildIndex(ctx, "", config.LOLConfig.Cron.ReBuild)
	if err != nil {
		log.Logger.Error(ctx, err)
	}

	// mongo
	ExtractKeyWords(ctx, common.PlatformForLOL)
	ExtractKeyWords(ctx, common.PlatformForLOLM)
}
