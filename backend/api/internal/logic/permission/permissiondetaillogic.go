package logic

import (
	"context"
	"neekity.com/go-admin/common/helper"

	"neekity.com/go-admin/api/internal/svc"
	"neekity.com/go-admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) PermissionDetailLogic {
	return PermissionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionDetailLogic) PermissionDetail(req types.PermissionDetailRequest) (*types.ApiResponse, error) {
	data, err := l.svcCtx.PermissionModel.FindOne(req.Id)
	if err != nil {
		return (*types.ApiResponse)(helper.ApiError(err.Error(), nil)), nil
	}
	return (*types.ApiResponse)(helper.ApiSuccess(data)), nil
}
