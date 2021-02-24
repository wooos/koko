package service

import "time"

func checkIfNeedLoginAssetConfirm(userID, assetID, systemUserID string) bool {

	return true
}

func checkIfNeedLoginAppConfirm(userID, assetID, systemUserID string) bool {

	return false
}

func checkLoginAssetConfirmFinish(userID, assetID, systemUserID string) (confirmResponse, error) {
	time.Sleep(10 * time.Second)
	return confirmResponse{
		Msg: successMsg,
	}, nil
	//return confirmResponse{}, nil
}

func checkLoginAPPConfirmFinish(userID, appID, systemUserID string) (confirmResponse, error) {

	return confirmResponse{}, nil
}

func cancelAssetConfirmLogin(userID, assetID, systemUserID string) {

}

func cancelAPPConfirmLogin(userID, appID, systemUserID string) {
}
