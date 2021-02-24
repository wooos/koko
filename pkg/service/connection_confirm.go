package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jumpserver/koko/pkg/model"
)

type connectionConfirmOption struct {
	user       *model.User
	systemUser *model.SystemUser

	targetType string
	targetID   string
}

func NewConnectionConfirm(opts ...ConfirmOption) ConnectionConfirm {
	var option connectionConfirmOption
	for _, setter := range opts {
		setter(&option)
	}
	return ConnectionConfirm{option: &option}
}

type ConnectionConfirm struct {
	option *connectionConfirmOption
}

func (c *ConnectionConfirm) WaitLoginConfirm(ctx context.Context) error {
	// TODO: 通过登录复核的工单，检查复核
	return c.waitConfirmFinish(ctx)
}

func (c *ConnectionConfirm) waitConfirmFinish(ctx context.Context) error {
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

func (c *ConnectionConfirm) CheckIsNeedLoginConfirm() bool {
	// todo: 获取登录复核的工单信息
	userID := c.option.user.ID
	systemUserID := c.option.systemUser.ID
	targetID := c.option.targetID
	switch c.option.targetType {
	case model.AppType:
		return checkIfNeedAppConnectionConfirm(userID, targetID, systemUserID)
	default:
		return checkIfNeedAssetConnectionConfirm(userID, targetID, systemUserID)
	}
}

func (c *ConnectionConfirm) checkConfirmFinish() (confirmResponse, error) {
	userID := c.option.user.ID
	systemUserID := c.option.systemUser.ID
	targetID := c.option.targetID
	switch c.option.targetType {
	case model.AppType:
		return checkAPPConnectionConfirmFinish(userID, targetID, systemUserID)
	default:
		return checkLoginAssetConfirmFinish(userID, targetID, systemUserID)
	}
}

func (c *ConnectionConfirm) cancelConfirm() {
	userID := c.option.user.ID
	systemUserID := c.option.systemUser.ID
	targetID := c.option.targetID
	switch c.option.targetType {
	case model.AppType:
		cancelAPPConnectionConfirm(userID, targetID, systemUserID)
	default:
		cancelAssetConnectionConfirm(userID, targetID, systemUserID)
	}
}

type confirmResponse struct {
	Msg string `json:"msg"`
	Err string `json:"error,omitempty"`
}

type ConfirmOption func(*connectionConfirmOption)

func ConfirmWithUser(user *model.User) ConfirmOption {
	return func(option *connectionConfirmOption) {
		option.user = user
	}
}

func ConfirmWithSystemUser(sysUser *model.SystemUser) ConfirmOption {
	return func(option *connectionConfirmOption) {
		option.systemUser = sysUser
	}
}

func ConfirmWithTargetType(targetType string) ConfirmOption {
	return func(option *connectionConfirmOption) {
		option.targetType = targetType
	}
}

func ConfirmWithTargetID(targetID string) ConfirmOption {
	return func(option *connectionConfirmOption) {
		option.targetID = targetID
	}
}
