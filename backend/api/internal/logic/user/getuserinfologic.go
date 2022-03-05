package logic

import (
	"context"
	"neekity.com/go-admin/common/helper"

	"neekity.com/go-admin/api/internal/svc"
	"neekity.com/go-admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserInfoLogic {
	return GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(userId uint) (*types.ApiResponse, error) {
	data, err := l.svcCtx.UserModel.GetUserInfo(userId, l.svcCtx.CasbinRuleEnforce)
	if err != nil {
		return (*types.ApiResponse)(helper.ApiError(err.Error(), nil)), nil
	}
	return (*types.ApiResponse)(helper.ApiSuccess(data)), nil
}
