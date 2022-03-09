package config

import (
	"github.com/zeromicro/go-queue/dq"
	"github.com/zeromicro/go-zero/rest"
	"golang.org/x/oauth2"
)

type Config struct {
	rest.RestConf
	DataSource    string
	FrontUrl      string
	BeanstalkConf []dq.Beanstalk
	oauth2.Config
	UserInfoUrl string
	JwtAuth     struct {
		AccessSecret string
		AccessExpire int64
	}
}
