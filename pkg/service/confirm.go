package service

import "time"

func checkIfNeedAssetConnectionConfirm(userID, assetID, systemUserID string) bool {

	return true
}

func checkIfNeedAppConnectionConfirm(userID, assetID, systemUserID string) bool {

	return false
}

func checkLoginAssetConfirmFinish(userID, assetID, systemUserID string) (confirmResponse, error) {
	time.Sleep(10 * time.Second)
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
