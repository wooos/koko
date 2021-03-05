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

	ticketID string // 获取到的当前工单 id
}

func (c *ConnectionConfirm) WaitLoginConfirm(ctx context.Context) error {
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
			res, err := c.getTicketStatus()
			if err != nil {
				return model.ErrConfirmRequestFailure
			}
			fmt.Println(res)
			switch res.Status {
			case TicketStatusOpen:
				continue
			case TicketStatusClosed:
				switch res.Action {
				case TicketActionApprove:
					return nil
				case TicketActionReject:
					return model.ErrConfirmReject
				case TicketActionClose:
					return model.ErrConfirmReject
				}
				return model.ErrConfirmReject
			default:
				return fmt.Errorf("unkonw status: %s", res.Status)
			}

		}
	}
}

func (c *ConnectionConfirm) CheckIsNeedLoginConfirm() (bool, error) {
	userID := c.option.user.ID
	systemUserID := c.option.systemUser.ID
	systemUsername := c.option.systemUser.Username
	targetID := c.option.targetID
	switch c.option.targetType {
	case model.AppType:
		return checkIfNeedAppConnectionConfirm(userID, targetID, systemUserID)
	default:
		res, err := checkIfNeedAssetLoginConfirm(userID, targetID,
			systemUserID, systemUsername)
		if err != nil {
			return false, err
		}
		if !res.Msg {
			c.ticketID = res.TicketID
		}
		return !res.Msg, nil
	}
}

func (c *ConnectionConfirm) getTicketStatus() (model.Ticket, error) {
	return GetTicketStatus(c.ticketID)
}

func (c *ConnectionConfirm) cancelConfirm() {
	userID := c.option.user.ID
	switch c.option.targetType {
	case model.AppType:
		closeTicketByUser(userID, c.ticketID)
	default:
		closeTicketByUser(userID, c.ticketID)
	}
}

type checkAssetConfirmResponse struct {
	Msg      bool   `json:"msg"`
	Err      string `json:"error,omitempty"`
	TicketID string `json:"ticket_id"`
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
