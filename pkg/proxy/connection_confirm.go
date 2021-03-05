package proxy

import (
	"context"

	"github.com/jumpserver/koko/pkg/i18n"
	"github.com/jumpserver/koko/pkg/logger"
	"github.com/jumpserver/koko/pkg/model"
	"github.com/jumpserver/koko/pkg/service"
	"github.com/jumpserver/koko/pkg/utils"
)

// 校验用户登录资产是否需要复核
func validateConnectionLoginConfirm(srv *service.ConnectionConfirm, userCon UserConnection) bool {
	ok, err := srv.CheckIsNeedLoginConfirm()
	if err != nil {
		msg := i18n.T("validate Login confirm err: Core Api failed")
		utils.IgnoreErrWriteString(userCon, msg)
		return false
	}
	if !ok {
		return true
	}

	ctx, cancelFunc := context.WithCancel(userCon.Context())
	defer userCon.Close()
	go func() {
		defer cancelFunc()
		term := utils.NewTerminal(userCon, ">: ")
		msg := "Wait for your admin to confirm login [Y/n]?\r\n"
		_, _ = term.Write([]byte(msg))
		for {
			line, err := term.ReadLine()
			if err != nil {
				return
			}
			switch line {
			case "exit", "quit", "q", "n":
				return
			}
			_, _ = term.Write([]byte(msg))
		}

	}()
	if err = srv.WaitLoginConfirm(ctx); err != nil {
		logger.Error("Check admin Confirm login session failed: " + err.Error())
		utils.IgnoreErrWriteString(userCon, getErrI18nMsg(err))
		return false
	}
	return true
}

func getErrI18nMsg(err error) string {
	var msg string
	switch err {
	case model.ErrConfirmReject:
		msg = i18n.T("Reject login asset")
	case model.ErrConfirmRequestFailure:
		msg = i18n.T("API Core failed")
	default:
		msg = i18n.T("Unknown err: " + err.Error())
	}
	return msg
}
