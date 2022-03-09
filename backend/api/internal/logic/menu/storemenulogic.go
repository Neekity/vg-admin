package logic

import (
	"context"
	"neekity.com/go-admin/common/helper"
	"neekity.com/go-admin/model"

	"neekity.com/go-admin/api/internal/svc"
	"neekity.com/go-admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StoreMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStoreMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) StoreMenuLogic {
	return StoreMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StoreMenuLogic) StoreMenu(req types.StoreMenuRequest) (*types.ApiResponse, error) {
	menuInfo := model.Menu{
		Name:     req.Name,
		Path:     req.Path,
		Icon:     req.Icon,
		ParentId: req.ParentId,
	}
	menu, err := l.svcCtx.MenuModel.UpdateOrCreate(req.Id, menuInfo)
	if err != nil {
		return (*types.ApiResponse)(helper.ApiError(err.Error(), nil)), nil
	}

	return (*types.ApiResponse)(helper.ApiSuccess(menu)), nil
}
