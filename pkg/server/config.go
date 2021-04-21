package server

import (
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/mysql"
	"github.com/li4n0/revsuit/pkg/rhttp"
)

type noticeConfig struct {
	DingTalk string
	Lark     string
	WeiXin   string
	Slack    string
}

type Config struct {
	Addr     string
	Token    string
	Database string
	LogLevel string
	Notice   noticeConfig
	rhttp.Config
	DNS   dns.Config
	Mysql mysql.Config
}