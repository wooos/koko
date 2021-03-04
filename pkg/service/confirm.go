package service

func checkIfNeedAssetLoginConfirm(userID, assetID, systemUserID,
	sysUsername string) (string, bool, error) {
	params := map[string]string{
		"user_id":         userID,
		"asset_id":        assetID,
		"system_user_id":  systemUserID,
		"system_username": sysUsername,
	}
	var res struct {
		Msg      bool   `json:"msg"`
		Err      string `json:"error"`
		TicketId string `json:"ticket_id"`
	}
	if _, err := authClient.Get(AssetLoginConfirmURL, &res, params); err != nil {
		return "", false, err
	}
	return res.TicketId, !res.Msg, nil
}

func checkIfNeedAppConnectionConfirm(userID, assetID, systemUserID string) (bool, error) {

	return false, nil
}

func checkTicketFinish(ticketID string) (confirmResponse, error) {



	return confirmResponse{
		Msg: successMsg,
	}, nil
	//return confirmResponse{}, nil
}

func checkAPPConnectionConfirmFinish(userID, appID, systemUserID string) (confirmResponse, error) {

	return confirmResponse{}, nil
}

func cancelAssetConnectionConfirm(userID, assetID, systemUserID string) {

}

func cancelAPPConnectionConfirm(userID, appID, systemUserID string) {
}
