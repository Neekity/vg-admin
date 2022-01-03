package middleware

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/tal-tech/go-zero/rest/httpx"
	"neekity.com/go-admin/common/constant"
	"neekity.com/go-admin/common/helper"
	"net/http"
	"strconv"
	"strings"
)

type AuthMiddleware struct {
	AccessSecret string
	enforcer     *casbin.Enforcer
}

type tokenClaim struct {
	Sub uint `json:"sub"`
}

func (t tokenClaim) Valid() error {
	return nil
}

func NewAuthMiddleware(secret string, enforcer *casbin.Enforcer) *AuthMiddleware {
	return &AuthMiddleware{AccessSecret: secret, enforcer: enforcer}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.OkJson(w, map[string]string{
				"code": "2003",
				"msg":  "请求头中auth为空",
			})
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			httpx.OkJson(w, map[string]string{
				"code": "2004",
				"msg":  "请求头中auth格式有误",
			})
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		claims, err := getClaimsByToken(parts[1], m.AccessSecret)
		if err != nil {
			httpx.OkJson(w, map[string]string{
				"code": "2005",
				"msg":  "无效的Token" + err.Error(),
			})
			return
		}

		userId := strconv.Itoa(int(claims.Sub))
		r.Header.Set("UserId", userId)
		if err := m.enforcer.LoadPolicy(); err != nil {
			httpx.OkJson(w, helper.ApiError(err.Error(), nil))
			return
		}
		if isAdmin := m.enforcer.HasGroupingPolicy(userId, constant.AdminRole); isAdmin == false {
			flag, err := m.enforcer.Enforce(userId, r.URL.Path, constant.CasbinPermissionWrite)
			if err != nil {
				httpx.OkJson(w, helper.ApiError(err.Error(), nil))
				return
			}
			if !flag {
				httpx.OkJson(w, helper.ApiError("未拥有该路由的权限", nil))
				return
			}
		}
		next(w, r)
	}
}

func getClaimsByToken(tokenString string, secret string) (*tokenClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*tokenClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
