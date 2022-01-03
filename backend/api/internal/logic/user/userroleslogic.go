package logic

import (
	"context"
	"neekity.com/go-admin/common/helper"

	"neekity.com/go-admin/api/internal/svc"
	"neekity.com/go-admin/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserRolesLogic {
	return UserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRolesLogic) UserRoles(req types.UserRolesRequest) (*types.ApiResponse, error) {
	roles, err := l.svcCtx.UserModel.GetUserRoles(req.UserId, l.svcCtx.CasbinRuleEnforce)
	if err != nil {
		return (*types.ApiResponse)(helper.ApiError(err.Error(), nil)), nil
	}
	return (*types.ApiResponse)(helper.ApiSuccess(roles)), nil
}
