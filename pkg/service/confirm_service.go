package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jumpserver/koko/pkg/model"
)

type confirmOption struct {
	user       *model.User
	systemUser *model.SystemUser

	targetType string
	targetID   string
}

func NewConfirmService(opts ...ConfirmOption) ConfirmService {
	var option confirmOption
	for _, setter := range opts {
		setter(&option)
	}
	return ConfirmService{option: &option}
}

type ConfirmService struct {
	option *confirmOption
}

func (c *ConfirmService) WaitLoginConfirm(ctx context.Context) error {
	// TODO: 通过登录复核的工单，检查复核
	return c.waitConfirmFinish(ctx)
}

func (c *ConfirmService) waitConfirmFinish(ctx context.Context) error {
	// 10s 请求一次
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			c.cancelConfirm()
			return model.ErrConfirmCancel
		case <-t.C:
			res, err := c.checkConfirmFinish()
			if err != nil {
				return model.ErrConfirmRequestFailure
			}
			if res.Err != "" {
				switch res.Err {
				case ErrSessionLoginConfirmWait:
					continue
				case ErrSessionLoginConfirmRejected:
					return model.ErrConfirmReject
				}
				return fmt.Errorf("unkonw err: %s", res.Err)
			}
			if res.Msg == successMsg {
				return nil
			}
		}
	}
}

func (c *ConfirmService) CheckIsNeedLoginConfirm() bool {
	// todo: 获取登录复核的工单信息
	userID := c.option.user.ID
	systemUserID := c.option.systemUser.ID
	targetID := c.option.targetID
	switch c.option.targetType {
	case model.AppType:
		return checkIfNeedLoginAppConfirm(userID, targetID, systemUserID)
	default:
		return checkIfNeedLoginAssetConfirm(userID, targetID, systemUserID)
	}
}

func (c *ConfirmService) checkConfirmFinish() (confirmResponse, error) {
	userID := c.option.user.ID
	systemUserID := c.option.systemUser.ID
	targetID := c.option.targetID
	switch c.option.targetType {
	case model.AppType:
		return checkLoginAPPConfirmFinish(userID, targetID, systemUserID)
	default:
		return checkLoginAssetConfirmFinish(userID, targetID, systemUserID)
	}
}

func (c *ConfirmService) cancelConfirm() {
	userID := c.option.user.ID
	systemUserID := c.option.systemUser.ID
	targetID := c.option.targetID
	switch c.option.targetType {
	case model.AppType:
		cancelAPPConfirmLogin(userID, targetID, systemUserID)
	default:
		cancelAssetConfirmLogin(userID, targetID, systemUserID)
	}
}

type confirmResponse struct {
	Msg string `json:"msg"`
	Err string `json:"error,omitempty"`
}

type ConfirmOption func(*confirmOption)

func ConfirmWithUser(user *model.User) ConfirmOption {
	return func(option *confirmOption) {
		option.user = user
	}
}

func ConfirmWithSystemUser(sysUser *model.SystemUser) ConfirmOption {
	return func(option *confirmOption) {
		option.systemUser = sysUser
	}
}

func ConfirmWithTargetType(targetType string) ConfirmOption {
	return func(option *confirmOption) {
		option.targetType = targetType
	}
}

func ConfirmWithTargetID(targetID string) ConfirmOption {
	return func(option *confirmOption) {
		option.targetID = targetID
	}
}
