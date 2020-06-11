package handler

import (
	"business/services/base/ifaccount/adaptor"
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
	"strconv"
)

type RewardsRes struct {
	Rewards []*domain.UserAssetReward `json:"rewards"`
}

func Rewards(r *httpserver.Request) *httpserver.Response {
	coin := r.QueryParams.Get("coin")
	rewardType, err := strconv.Atoi(r.QueryParams.Get("type"))
	if err != nil {
		rewardType = 0
	}
	limit, err := strconv.Atoi(r.QueryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(r.QueryParams.Get("offset"))
	if err != nil || offset <= 0 {
		offset = 0
	}
	rewards, err := adaptor.GetUserAssetRewards(domain.USER_REWARD_TYPE(rewardType), coin, r.Uid, limit, offset)
	if err != nil {
		loggers.Error.Printf("Get Rewards error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	resp := httpserver.NewResponse()
	resp.Data = &RewardsRes{Rewards: rewards}
	return resp
}
