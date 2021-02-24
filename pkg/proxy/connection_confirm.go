package proxy

import (
	"io"
	"strings"

	"github.com/jumpserver/koko/pkg/i18n"
	"github.com/jumpserver/koko/pkg/logger"
	"github.com/jumpserver/koko/pkg/model"
	"github.com/jumpserver/koko/pkg/service"
	"github.com/jumpserver/koko/pkg/utils"
)

func checkAdminConfirmConnection(srv *service.ConnectionConfirm, userCon UserConnection) bool {
	if !srv.CheckIsNeedLoginConfirm() {
		return true
	}
	if !waitUserConfirm(userCon) {
		return false
	}
	if err := srv.WaitLoginConfirm(userCon.Context()); err != nil {
		logger.Error("Check admin Confirm login session failed: " + err.Error())
		utils.IgnoreErrWriteString(userCon, getErrI18nMsg(err))
		return false
	}
	return true
}

func waitUserConfirm(rw io.ReadWriteCloser) bool {
	opt := i18n.T("Do you wait for your admin to confirm login [Y/n]? :")
	term := utils.NewTerminal(rw, opt)
	line, err := term.ReadLine()
	if err != nil {
		return false
	}
	switch strings.ToLower(line) {
	case "yes", "y":
		return true
	default:
		return false
	}
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
