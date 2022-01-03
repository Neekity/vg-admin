package handler

import (
	"net/http"

	"neekity.com/go-admin/api/internal/logic/permission"
	"neekity.com/go-admin/api/internal/svc"
	"neekity.com/go-admin/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func StorePermissionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StorePermissionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewStorePermissionLogic(r.Context(), ctx)
		resp, err := l.StorePermission(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
