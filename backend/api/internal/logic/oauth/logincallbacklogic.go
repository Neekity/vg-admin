package logic

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"neekity.com/go-admin/api/internal/types"
	"neekity.com/go-admin/common/helper"
	"neekity.com/go-admin/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"neekity.com/go-admin/api/internal/svc"
)

type LoginCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginCallbackLogic {
	return LoginCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginCallbackLogic) LoginCallback(code string) (*types.JwtTokenResponse, error) {
	token, err := l.svcCtx.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	_, userInfo := helper.HttpGet(l.svcCtx.Config.UserInfoUrl, map[string]string{
		"Content-Type": "application/json", "Authorization": "Bearer " + token.AccessToken})
	var userInfoMap types.BasicUserInfo
	logx.Error(userInfo)
	jsonUnmarshalErr := json.Unmarshal([]byte(userInfo), &userInfoMap)
	if jsonUnmarshalErr != nil {
		logx.Error("格式化用户信息:" + userInfo + "失败" + jsonUnmarshalErr.Error())
		return nil, jsonUnmarshalErr
	}
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	if err != nil {
		return nil, err
	}
	user, modelErr := l.svcCtx.UserModel.UpdateOrCreate(userInfoMap.Email, model.User{
		Avatar:   userInfoMap.Photo,
		Password: string(password),
		Name:     userInfoMap.Name,
		Email:    userInfoMap.Email,
	})
	if modelErr != nil {
		return nil, modelErr
	}

	tokenResp, jwtErr := l.Jwt(user.ID)
	if jwtErr != nil {
		return nil, jwtErr
	}
	return tokenResp, nil
}

func (l *LoginCallbackLogic) Jwt(sub uint) (*types.JwtTokenResponse, error) {
	var accessExpire = l.svcCtx.Config.JwtAuth.AccessExpire

	now := time.Now().Unix()
	accessToken, err := l.GenToken(now, l.svcCtx.Config.JwtAuth.AccessSecret, map[string]interface{}{"sub": sub}, accessExpire)
	if err != nil {
		return nil, err
	}

	return &types.JwtTokenResponse{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *LoginCallbackLogic) GenToken(iat int64, secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
