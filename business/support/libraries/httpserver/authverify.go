package httpserver

import (
	"business/support/libraries/auth"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"fmt"
	"math"
	"time"
)

// normalAuthVerify
func normalAuthVerify(r *Request, f HandleFunc) *Response {
	session, sessionConf, res := checkRequestHeader(r)
	if nil != res {
		return res
	}
	//if auth.IsInLimit(r.Uid) {
	//	return NewResponseWithError(errors.SysAccess2OftenError)
	//}
	resp := f(r)
	if time.Since(*session.CreatedAt) >= sessionConf.Timeout || resp.IsResetToken {
		session, err := auth.NewSession(session.Account, r.Ver, r.Dev, r.Ts)
		if err != nil {
			return NewResponseWithError(errors.InternalServerError)
		}
		resp.SetSession(session)
	}
	return resp
}

// assetAuthVerify
func assetAuthVerify(r *Request, f HandleFunc) *Response {
	session, sessionConf, res := checkRequestHeader(r)
	if nil != res {
		return res
	}

	//if auth.IsInLimit(r.Uid) {
	//	return NewResponseWithError(errors.SysAccess2OftenError)
	//}

	// 资金密码已设置,并且存在有效时间
	if session.Account.AssetPasswordEffectiveTime > 0 {
		if auth.IsDisableAssetOP(session.Account.AccountId) {
			ts := auth.GetDisableAssetOPTime(session.Account.AccountId)
			t := math.Floor(float64(time.Now().Unix()-ts) / 3600)
			f := fmt.Sprintf("%g", 24-t)
			loggers.Warn.Printf("assetAuthVerify account %d optime %d countdown %s", session.Account.AccountId, ts, f)
			if ts <= 1 {
				return NewResponseWithError(errors.DisableAssetOP)
			}
			return NewResponseWithFormatError(errors.DisableAssetOPCountDown, []interface{}{f})
		}
		if r.Header.Get("bifund-AssetPassword") != "" {
			if r.Header.Get("bifund-AssetPassword") != session.Account.AssetPassword {
				count, err := auth.IncorrectAssetPasswordIncr(session.Account.AccountId)
				if err != nil {
					loggers.Error.Printf("IncorrectAssetPasswordIncr error %s", err.Error())
				}
				if count >= 6 {
					auth.DisableAssetOP(session.Account.AccountId)
				}
				l := 6 - count
				if l <= 0 {
					//账号或密码错误，请在24小时后重试
					return NewResponseWithError(errors.DisableAssetOP)
				}

				//账户或密码输入错误，您还有x次尝试机会
				return NewResponseWithFormatError(errors.IncorretAssetPassword, []interface{}{l})
			}
		} else {
			if !auth.AssetTokenAuth(session.Token, time.Duration(session.Account.AssetPasswordEffectiveTime)*time.Second) {
				return NewResponseWithError(errors.AssetPermissionDenied)
			}
		}
	}

	// 资金密码已设置,单次有效
	if session.Account.AssetPasswordEffectiveTime == 0 {
		if auth.IsDisableAssetOP(session.Account.AccountId) {
			ts := auth.GetDisableAssetOPTime(session.Account.AccountId)
			t := math.Floor(float64(time.Now().Unix()-ts) / 3600)
			f := fmt.Sprintf("%g", 24-t)
			loggers.Warn.Printf("assetAuthVerify account %d optime %d countdown %s", session.Account.AccountId, ts, f)
			if ts <= 1 {
				return NewResponseWithError(errors.DisableAssetOP)
			}
			return NewResponseWithFormatError(errors.DisableAssetOPCountDown, []interface{}{f})
		}
		if r.Header.Get("bifund-AssetPassword") == "" {
			return NewResponseWithError(errors.AssetPermissionDenied)
		}
		if r.Header.Get("bifund-AssetPassword") != session.Account.AssetPassword {
			count, err := auth.IncorrectAssetPasswordIncr(session.Account.AccountId)
			if err != nil {
				loggers.Error.Printf("IncorrectAssetPasswordIncr error %s", err.Error())
			}
			if count >= 6 {
				auth.DisableAssetOP(session.Account.AccountId)
			}

			l := 6 - count
			if l <= 0 {
				//账号或密码错误，请在24小时后重试
				return NewResponseWithError(errors.DisableAssetOP)
			}

			//账户或密码输入错误，您还有x次尝试机会
			return NewResponseWithFormatError(errors.IncorretAssetPassword, []interface{}{l})
		}
	}
	resp := f(r)
	if time.Since(*session.CreatedAt) >= sessionConf.Timeout || resp.IsResetToken {
		var err error
		session, err = auth.NewSession(session.Account, r.Ver, r.Dev, r.Ts)
		if err != nil {
			return NewResponseWithError(errors.InternalServerError)
		}
		resp.SetSession(session)
	}

	// 校验成功,并且有效时间不为0,重置资金密码有效时间
	if session.Account.AssetPasswordEffectiveTime > 0 {
		auth.SetAssetSession(session.Token, time.Duration(session.Account.AssetPasswordEffectiveTime)*time.Second)
	}
	return resp
}
