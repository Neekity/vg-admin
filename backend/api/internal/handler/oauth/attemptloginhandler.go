package handler

import (
	"net/http"

	"neekity.com/go-admin/api/internal/logic/oauth"
	"neekity.com/go-admin/api/internal/svc"
)

func AttemptLoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAttemptLoginLogic(r.Context(), ctx)
		url := l.AttemptLogin()
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
