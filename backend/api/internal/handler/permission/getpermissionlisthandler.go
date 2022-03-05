package handler

import (
	"net/http"

	"neekity.com/go-admin/api/internal/logic/permission"
	"neekity.com/go-admin/api/internal/svc"
	"neekity.com/go-admin/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPermissionListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchPermissionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetPermissionListLogic(r.Context(), ctx)
		resp, err := l.GetPermissionList(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
