package handler

import (
	"business/services/base/ifaccount/adaptor"
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
	"business/support/model/config"
	"text/template"
)

func isIllegal(s string) bool {
	return (s != template.HTMLEscapeString(s) || s != template.JSEscapeString(s))
}

func KYCInfo(r *httpserver.Request) *httpserver.Response {
	account := r.Session.Account
	kycInfo := account.KYCInfo()
	if account.KYCStatus == domain.KYC_STATUS_APPROVED {
		if len(kycInfo.IDNo) >= 6 {
			kycInfo.IDNo = kycInfo.IDNo[:len(kycInfo.IDNo)-6] + "******"
		}
		kycInfo.IDPhoto1 = ""
		kycInfo.IDPhoto2 = ""
		kycInfo.IDPhoto3 = ""
	}
	if kycInfo.Nationality != "" {
		language := r.Language
		if r.Language != "en" && r.Language != "zh-cn" {
			language = "en"
		}
		if area, err := config.GetArea(kycInfo.Nationality, language); err != nil {
			kycInfo.Nationality = ""
		} else {
			kycInfo.NationalityLongName = area.LongName
		}

	}
	resp := httpserver.NewResponse()
	resp.Data = kycInfo
	return resp
}

func KYCUpdate(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var (
		req domain.Account
		err error
	)

	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("KYCUpdate parse request body error %s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if isIllegal(req.Nationality) ||
		isIllegal(req.FirstName) ||
		isIllegal(req.LastName) ||
		isIllegal(req.IDNo) ||
		isIllegal(req.IDPhoto1) ||
		isIllegal(req.IDPhoto2) ||
		isIllegal(req.IDPhoto3) {
		loggers.Warn.Printf("KYCUpdate illegal character")
		return httpserver.NewResponseWithError(errors.IllegalCharacter)
	}

	account := r.Session.Account

	if account.KYCStatus != domain.KYC_STATUS_REJECTED &&
		account.KYCStatus != domain.KYC_STATUS_NOT_CERTIFIED &&
		account.KYCStatus != domain.KYC_STATUS_EDITING {
		loggers.Warn.Printf("KYCUpdate not allow update current status:%d", account.KYCStatus)
		return httpserver.NewResponseWithError(errors.KYCNotAllowUpdate)
	}
	if account.Phone == "" {
		loggers.Warn.Printf("KYCUpdate not bind phone")
		return httpserver.NewResponseWithError(errors.NoBindPhone)
	}

	var action string
	req.KYCType = domain.KYC_TYPE_PERSONAL
	req.AccountId = r.Uid
	switch req.KYCStatus {
	case domain.KYC_STATUS_EDITING:
		action = domain.MESSAGE_ACTION_KYC_UPDATED
		err = adaptor.KYCUpdate(&req)
	case domain.KYC_STATUS_SUBMITTED:
		action = domain.MESSAGE_ACTION_KYC_SUBMITED
		if req.Nationality == "" {
			return httpserver.NewResponseWithError(errors.KYCNoNationality)
		}
		if req.FirstName == "" {
			return httpserver.NewResponseWithError(errors.KYCNoFirstName)
		}
		if req.LastName == "" {
			return httpserver.NewResponseWithError(errors.KYCNoLastName)
		}
		if req.IDType == 0 {
			return httpserver.NewResponseWithError(errors.KYCNoIDType)
		}
		if req.IDNo == "" {
			return httpserver.NewResponseWithError(errors.KYCNoIDNo)
		}
		if req.IDPhoto1 == "" {
			return httpserver.NewResponseWithError(errors.KYCNoPhoto1)
		}
		if req.IDPhoto2 == "" {
			return httpserver.NewResponseWithError(errors.KYCNoPhoto2)
		}
		if req.IDType == domain.ID_TYPE_ID_CARD || req.IDType == domain.ID_TYPE_DRIVER_LICENSE {
			if req.IDPhoto3 == "" {
				return httpserver.NewResponseWithError(errors.KYCNoPhoto3)
			}
		}
		ok, err := adaptor.CheckKYCIDNo(req.IDNo)
		if err != nil {
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
		if !ok {
			return httpserver.NewResponseWithError(errors.KYCIDNoHadExisted)
		}
		err = adaptor.KYCUpdate(&req)
	default:
		return httpserver.NewResponseWithError(errors.BadRequest)
	}

	if !domain.CheckURL(req.IDPhoto1) || !domain.CheckURL(req.IDPhoto2) || !domain.CheckURL(req.IDPhoto3) {
		return httpserver.NewResponseWithError(errors.BadRequest)
	}

	if err != nil {
		loggers.Error.Printf("KYCUpdate error %s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}

	publishAccountRecord(account, action, r.ClientIP(), "", 0)

	return httpserver.NewResponse()
}
