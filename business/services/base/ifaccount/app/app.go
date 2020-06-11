package app

import (
	"business/services/base/ifaccount/handler"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
	"net/http"
)

func init() {
	loggers.Info.Printf("IFAccount service init  [\033[0;32;1mOK\033[0m]")
}

func Start(addr string) {
	//Init()
	r := httpserver.NewRouter()
	r.RouteHandleFunc("/login", httpserver.Check(handler.Login)).Methods(http.MethodPost)
	r.RouteHandleFunc("/loginRecord", httpserver.NormalAuth(handler.GetLoginRecord)).Methods(http.MethodGet)
	r.RouteHandleFunc("/users/register", httpserver.Check(handler.Register)).Methods(http.MethodPost)

	r.RouteHandleFunc("/logout", httpserver.NormalAuth(handler.Logout)).Methods(http.MethodGet)
	r.RouteHandleFunc("/users/active", httpserver.NormalAuth(handler.ActiveAccount)).Methods(http.MethodPost)
	r.RouteHandleFunc("/users/resetPassword", handler.ResetPassword).Methods(http.MethodPost)
	r.RouteHandleFunc("/users/me", httpserver.NormalAuth(handler.UserMe)).Methods(http.MethodGet)
	r.RouteHandleFunc("/users/check", handler.CheckAccountExist).Methods(http.MethodGet)
	r.RouteHandleFunc("/settles", httpserver.NormalAuth(handler.Settles)).Methods(http.MethodGet)
	r.RouteHandleFunc("/v2/settles", httpserver.NormalAuth(handler.SettlesV2)).Methods(http.MethodGet)
	r.RouteHandleFunc("/rewards", httpserver.NormalAuth(handler.Rewards)).Methods(http.MethodGet)

	r.RouteHandleFunc("/verifyCode", handler.VerifyCodePost).Methods(http.MethodPost)
	r.RouteHandleFunc("/address", httpserver.NormalAuth(handler.DepositAddress)).Methods(http.MethodGet)
	r.RouteHandleFunc("/rechargeAmount", httpserver.NormalAuth(handler.RechargeAmount)).Methods(http.MethodPost)
	r.RouteHandleFunc("/user/accountName", httpserver.NormalAuth(handler.SetAccountName)).Methods(http.MethodPost)
	r.RouteHandleFunc("/user/avatar", httpserver.NormalAuth(handler.UpdateAccountAvatar)).Methods(http.MethodPost)

	r.RouteHandleFunc("/withdraw", httpserver.AssetAuth(handler.Withdraw)).Methods(http.MethodPost)

	r.RouteHandleFunc("/bindPhone", httpserver.NormalAuth(handler.BindPhone)).Methods(http.MethodPost)
	r.RouteHandleFunc("/bindEmail", httpserver.NormalAuth(handler.BindEmail)).Methods(http.MethodPost)

	// WithdrawAddress
	r.RouteHandleFunc("/withdrawAddress", httpserver.NormalAuth(handler.AddWithdrawAddress)).Queries("action", "add").Methods(http.MethodPost)
	r.RouteHandleFunc("/withdrawAddress", httpserver.NormalAuth(handler.GetWithdrawAddress)).Queries("action", "query").Methods(http.MethodPost)
	r.RouteHandleFunc("/withdrawAddress", httpserver.NormalAuth(handler.UpdateWithdrawAddress)).Queries("action", "update").Methods(http.MethodPost)
	r.RouteHandleFunc("/withdrawAddress", httpserver.NormalAuth(handler.DeleteWithdrawAddress)).Queries("action", "delete").Methods(http.MethodPost)

	// AssetPassword
	r.RouteHandleFunc("/assetPassword", httpserver.NormalAuth(handler.AddAssetPassword)).Queries("action", "add").Methods(http.MethodPost)
	r.RouteHandleFunc("/assetPassword", httpserver.NormalAuth(handler.ResetAssetPassword)).Queries("action", "reset").Methods(http.MethodPost)
	r.RouteHandleFunc("/assetPasswordEffectiveTime", httpserver.NormalAuth(handler.ResetAssetPasswordEffectiveTime)).Queries("action", "reset").Methods(http.MethodPost)

	r.RouteHandleFunc("/antiFishingText", httpserver.NormalAuth(handler.SetAntiFishingText)).Queries("action", "set").Methods(http.MethodPost)

	// Google Authenticator
	r.RouteHandleFunc("/GAKey", httpserver.NormalAuth(handler.GetGAKey)).Queries("action", "query").Methods(http.MethodPost)
	r.RouteHandleFunc("/GAKey", httpserver.NormalAuth(handler.SetGAKey)).Queries("action", "add").Methods(http.MethodPost)
	r.RouteHandleFunc("/GAKey", httpserver.NormalAuth(handler.DeleteGAKey)).Queries("action", "delete").Methods(http.MethodPost)

	r.RouteHandleFunc("/KYCAuth", httpserver.NormalAuth(handler.KYCInfo)).Queries("action", "query").Methods(http.MethodPost)
	r.RouteHandleFunc("/KYCAuth", httpserver.NormalAuth(handler.KYCUpdate)).Queries("action", "update").Methods(http.MethodPost)

	r.RouteHandleFunc("/captchCheck", handler.CaptchCheck)

	go r.ListenAndServeCORS(addr)
	loggers.Info.Printf("IFAccount service start [\033[0;32;1mOK\t%+v\033[0m]", addr)
	select {}
}
