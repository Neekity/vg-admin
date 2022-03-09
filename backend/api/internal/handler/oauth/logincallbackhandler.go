package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"neekity.com/go-admin/api/internal/logic/oauth"
	"neekity.com/go-admin/api/internal/svc"
)

func LoginCallbackHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewLoginCallbackLogic(r.Context(), ctx)
		code := r.URL.Query().Get("code")
		token, err := l.LoginCallback(code)
		if err != nil {
			httpx.Error(w, err)
		} else {
			http.Redirect(w, r, ctx.Config.FrontUrl+token.AccessToken, http.StatusMovedPermanently)
		}
	}
}
