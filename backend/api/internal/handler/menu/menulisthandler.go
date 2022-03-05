package handler

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"neekity.com/go-admin/api/internal/logic/menu"
	"neekity.com/go-admin/api/internal/svc"
)

func MenuListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewMenuListLogic(r.Context(), ctx)
		userId, err := strconv.Atoi(r.Header.Get("UserId"))
		resp, err := l.MenuList(uint(userId))
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
