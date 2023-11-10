package logic

import (
	context2 "context"
	"sync"
	"time"
	"whisper/internal/logic/command"
	"whisper/internal/logic/common"
	"whisper/internal/logic/equipment"
	"whisper/internal/logic/suit"
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
		_, err := equipment.QueryEquipments(ctx, common.PlatformForLOL)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start QueryEquipments PlatformForLOLM...")
		_, err = equipment.QueryEquipments(ctx, common.PlatformForLOLM)
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
	// 装备、英雄 别名
	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
		}()

		log.Logger.Info(ctx, "start AliasHeroes...")
		_, err := AliasHeroes(ctx)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()
	go func() {
		defer func() {
			wg.Done()
		}()

		log.Logger.Info(ctx, "start AliasEquip...")
		_, err := AliasEquip(ctx)
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

		log.Logger.Info(ctx, "start BatchUpdateSuitEquip...")
		suit.NewSuit()(ctx, common.PlatformForLOL).BatchUpdateSuitEquip()
		suit.NewSuit()(ctx, common.PlatformForLOLM).BatchUpdateSuitEquip()

		log.Logger.Info(ctx, "start SuitData2Redis...")
		err := suit.SuitData2Redis(ctx)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		log.Logger.Info(ctx, "start SuitHeroData2Redis...")
		err = SuitHeroData2Redis(ctx)
		if err != nil {
			log.Logger.Error(ctx, err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Logger.Info(ctx, "mongo ExtractKeyWords START...")
		// mongo
		invoker := new(command.Invoker)
		equipForLOLCmd := equipment.NewInnerIns(ctx).WithPlatform(common.PlatformForLOL).NewExtractKeyWordsCmd()
		equipForLOLMCmd := equipment.NewInnerIns(ctx).WithPlatform(common.PlatformForLOLM).NewExtractKeyWordsCmd()
		invoker.AddCommand(equipForLOLCmd, equipForLOLMCmd)
		invoker.NonBlockRun()

		log.Logger.Info(ctx, "mongo ExtractKeyWords END")
	}()
	wg.Wait()

}
