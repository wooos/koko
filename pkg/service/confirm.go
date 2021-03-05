package service

import (
	"github.com/jumpserver/koko/pkg/logger"
	"github.com/jumpserver/koko/pkg/model"
)

func checkIfNeedAssetLoginConfirm(userID, assetID, systemUserID,
	sysUsername string) (res checkAssetConfirmResponse, err error) {
	params := map[string]string{
		"user_id":         userID,
		"asset_id":        assetID,
		"system_user_id":  systemUserID,
		"system_username": sysUsername,
	}

	_, err = authClient.Get(AssetLoginConfirmURL, &res, params)
	return
}

func checkIfNeedAppConnectionConfirm(userID, assetID, systemUserID string) (bool, error) {

	return false, nil
}

func GetTicketStatus(ticketID string) (res model.Ticket, err error) {
	params := map[string]string{
		"ticket_id": ticketID,
	}
	_, err = authClient.Get(TicketStatusURL, &res, params)
	if err != nil {
		logger.Errorf("Get Ticket err: %s", err.Error())
	}
	return
}

func closeTicketByUser(userID, ticketID string) {
	params := map[string]string{
		"user_id":   userID,
		"ticket_id": ticketID,
	}
	_, err := authClient.Delete(AssetLoginConfirmURL, nil, params)
	if err != nil {
		logger.Errorf("Close Ticket err: %s", err.Error())
	}
}
