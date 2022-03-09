package logic

import (
	"context"
	"golang.org/x/oauth2"

	"github.com/zeromicro/go-zero/core/logx"
	"neekity.com/go-admin/api/internal/svc"
)

type AttemptLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttemptLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) AttemptLoginLogic {
	return AttemptLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttemptLoginLogic) AttemptLogin() string {
	return l.svcCtx.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}
