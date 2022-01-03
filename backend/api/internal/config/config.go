package config

import (
	"github.com/tal-tech/go-zero/rest"
	"github.com/zeromicro/go-queue/dq"
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
