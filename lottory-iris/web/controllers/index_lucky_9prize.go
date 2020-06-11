package controllers

import (
	"lottery-iris/models"
	"lottery-iris/conf"
	"lottery-iris/services"
)

func (api *LuckyApi) prize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift *models.ObjGiftPrize
	giftList := services.NewGiftService().GetAllUse(true)
	for _, gift := range giftList {
		if gift.PrizeCodeA <= prizeCode &&
			gift.PrizeCodeB >= prizeCode {
			// 中奖编码区间满足条件，说明可以中奖
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall {
				prizeGift = &gift
				break
			}
		}
	}
	return prizeGift
}
