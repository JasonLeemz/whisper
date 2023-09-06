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

		log.Logger.Info(ctx, "start QueryEquipments PlatformForLOL...")
		_, err := QueryEquipments(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start QueryEquipments PlatformForLOLM...")
		_, err = QueryEquipments(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
		}()

		log.Logger.Info(ctx, "start QueryHeroes PlatformForLOL...")
		_, err := QueryHeroes(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start QueryHeroes PlatformForLOLM...")
		_, err = QueryHeroes(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start HeroAttribute PlatformForLOL...")
		_, err = HeroAttribute(ctx, "0", common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start HeroAttribute PlatformForLOLM...")
		_, err = HeroAttribute(ctx, "0", common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
		}()

		log.Logger.Info(ctx, "start QueryRune PlatformForLOL...")
		_, err := QueryRune(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start QueryRune PlatformForLOLM...")
		_, err = QueryRune(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
		}()

		log.Logger.Info(ctx, "start QuerySkill PlatformForLOL...")
		_, err := QuerySkill(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start QuerySkill PlatformForLOLM...")
		_, err = QuerySkill(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	wg.Wait()

	// ------------------------

	wg.Add(3)

	go func() {
		defer wg.Done()

		log.Logger.Info(ctx, "start building index...")
		err := BuildIndex(ctx, "", config.LOLConfig.Cron.ReBuild)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer wg.Done()

		log.Logger.Info(ctx, "start HeroesPosition PlatformForLOLM...")
		_, err := HeroesPosition(ctx, common.PlatformForLOLM)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start BatchUpdateSuitEquip...")
		err = BatchUpdateSuitEquip(ctx)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start SuitData2Redis...")
		err = SuitData2Redis(ctx)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer wg.Done()

		// mongo
		log.Logger.Info(ctx, "start mongo ExtractKeyWords PlatformForLOL...")
		ExtractKeyWords(ctx, common.PlatformForLOL)
		log.Logger.Info(ctx, "start mongo ExtractKeyWords PlatformForLOLM...")
		ExtractKeyWords(ctx, common.PlatformForLOLM)
	}()
	wg.Wait()

}
