package proxy

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	defer cancelFunc()
	go func() {
		defer cancelFunc()
		term := utils.NewTerminal(userCon, "")
		defer userCon.Write([]byte("\r\n"))
		for {
			line, err := term.ReadLine()
			if err != nil {
				logger.Errorf("Wait confirm err: %s", err.Error())
				return
			}
			switch line {
			case "quit", "q":
				return
			}
		}
	}()
	go func() {
		msg := i18n.T("Waiting for your admin to confirm, enter q to exit. ")
		delay := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				delayS := fmt.Sprintf("%ds", delay)
				data := strings.Repeat("\x08", len(delayS)+len(msg)) + msg + delayS
				utils.IgnoreErrWriteString(userCon, data)
				time.Sleep(time.Second)
				delay += 1
			}
		}
	}()

	if err = srv.WaitLoginConfirm(ctx); err != nil {
		logger.Error("Check confirm login session failed: " + err.Error())
		utils.IgnoreErrWriteString(userCon, getErrI18nMsg(err)+"\r\n")
		return false
	}
	return true
}

func getErrI18nMsg(err error) string {
	var msg string
	switch err {
	case model.ErrConfirmCancel:
		return i18n.T("Cancel login confirm")
	case model.ErrConfirmReject:
		msg = i18n.T("Reject login asset")
	case model.ErrConfirmRequestFailure:
		msg = i18n.T("API Core failed")
	default:
		msg = i18n.T("Unknown err: " + err.Error())
	}
	return msg
}
