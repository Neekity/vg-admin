package logic

import (
	"context"
	"neekity.com/go-admin/common/helper"
	"neekity.com/go-admin/model"

	"neekity.com/go-admin/api/internal/svc"
	"neekity.com/go-admin/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type StorePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) StorePermissionLogic {
	return StorePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StorePermissionLogic) StorePermission(req types.StorePermissionRequest) (*types.ApiResponse, error) {
	permissionInfo := model.Permission{
		CasbinPermission:     req.CasbinPermission,
		Name:                 req.Name,
		Route:                req.Route,
		CasbinPermissionType: req.CasbinPermissionType,
	}
	permission, err := l.svcCtx.PermissionModel.UpdateOrCreate(req.Id, permissionInfo)
	if err != nil {
		return (*types.ApiResponse)(helper.ApiError(err.Error(), nil)), nil
	}

	return (*types.ApiResponse)(helper.ApiSuccess(permission)), nil
}
