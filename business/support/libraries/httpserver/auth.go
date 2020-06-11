package httpserver

import (
	"business/support/domain"
	"business/support/libraries/auth"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
)

// Check http header
func Check(f HandleFunc) HandleFunc {
	return func(r *Request) *Response {
		// 这些判断已经挪到checkRequestHeader了
		if r.Ver == "" {
			loggers.Debug.Printf("No version input")
			return NewResponseWithError(errors.BadRequest)
		}
		if r.Dev == "" {
			loggers.Debug.Printf("No dev input")
			return NewResponseWithError(errors.BadRequest)
		}
		return f(r)
	}
}

func checkRequestHeaderBase(r *Request) *Response {
	if r.Ts == 0 {
		loggers.Debug.Printf("No ts")
		return NewResponseWithError(errors.BadRequest)
	}
	if r.Dev == "" {
		loggers.Debug.Printf("No dev")
		return NewResponseWithError(errors.BadRequest)
	}
	if _, ok := domain.DevMap[r.Dev]; !ok {
		loggers.Debug.Printf("Invalid dev:%s", r.Dev)
		return NewResponseWithError(errors.BadRequest)
	}
	return nil
}

func checkRequestHeader(r *Request) (*domain.Session, *domain.SessionConf, *Response) {
	if r.Ts == 0 {
		loggers.Debug.Printf("No ts")
		return nil, nil, NewResponseWithError(errors.BadRequest)
	}
	if r.Sign == "" {
		loggers.Debug.Printf("No sign")
		return nil, nil, NewResponseWithError(errors.BadRequest)
	}
	if r.Uid == 0 {
		loggers.Debug.Printf("No uid")
		return nil, nil, NewResponseWithError(errors.BadRequest)
	}
	if r.Dev == "" {
		loggers.Debug.Printf("No dev")
		return nil, nil, NewResponseWithError(errors.BadRequest)
	}
	if _, ok := domain.DevMap[r.Dev]; !ok {
		loggers.Debug.Printf("Invalid dev:%s", r.Dev)
		return nil, nil, NewResponseWithError(errors.BadRequest)
	}
	sessionConf := auth.SessionConfAdaptor(r.Ver)
	if nil == sessionConf {
		loggers.Debug.Printf("Invalid version")
		return nil, nil, NewResponseWithError(errors.BadRequest)
	}
	session, err := auth.SAuth(r.Uid, r.Sign, r.Ver, r.Dev, r.Ts)
	if err != nil {
		loggers.Debug.Printf("Session auth error:%s", err.Error())
		return nil, nil, NewResponseWithError(err.(*errors.Error))
	}
	if session.Account.AccountStatus == domain.ACCOUNT_STATUS_PASSWORD_SAFE_WARN {
		loggers.Warn.Printf("Auth uid:%d safe warn", session.Account.AccountId)
		return nil, nil, NewResponseWithError(errors.AccountSafeWarn)
	}
	r.Session = session
	return session, sessionConf, nil
}

// Auth 自适应适配三种授权方式,普通用户授权认证,acceesKey授权认证,cloud授权认证
func Auth(f HandleFunc) HandleFunc {
	return innerAuth(f, int(AUTH_TYPE_NORMAL|AUTH_TYPE_ACCESS|AUTH_TYPE_CLOUD))
}

// NormalAuth 普通用户授权认证
func NormalAuth(f HandleFunc) HandleFunc {
	return func(r *Request) *Response {
		if r.AuthType() == AUTH_TYPE_NORMAL {
			return normalAuthVerify(r, f)
		} else {
			return NewResponseWithError(errors.BadRequest)
		}
	}
}

func innerAuth(f HandleFunc, authType int) HandleFunc {
	return func(r *Request) *Response {
		if r.IsSupportAuth(authType) {
			if r.AuthType() == AUTH_TYPE_NORMAL {
				return normalAuthVerify(r, f)
			}
		}
		return NewResponseWithError(errors.BadRequest)
	}
}

// AssetAuth 带交易密码的验证
func AssetAuth(f HandleFunc) HandleFunc {
	return func(r *Request) *Response {
		if r.AuthType() == AUTH_TYPE_NORMAL {
			return assetAuthVerify(r, f)
		} else {
			return NewResponseWithError(errors.BadRequest)
		}
	}
}

// Keystore private key auth
func KAuth(f HandleFunc) HandleFunc {
	return Auth(func(r *Request) *Response {
		if r.ActionSign == "" {
			loggers.Debug.Printf("No action sign")
			return NewResponseWithError(errors.Forbidden)
		}
		if err := auth.KAuth(string(r.Address), r.ActionSign, r.BodyBuff.Bytes()); err != nil {
			loggers.Debug.Printf("KAuth error")
			return NewResponseWithError(errors.Forbidden)
		}
		return f(r)
	})
}

func CleanSession(uid int64) {
	for dev, _ := range domain.DevMap {
		if err := auth.DeleteSession(uid, dev); err != nil {
			loggers.Error.Printf("CleanSession %d %s error", uid, dev, err.Error())
		}
	}
}
