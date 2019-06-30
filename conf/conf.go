package conf

import (
	"strings"
	"time"

	"github.com/koding/multiconfig"
	"gitlab.azbit.cn/web/bitcoin/library/util"
)

const (
	// default to be disabled, please DON'T enable it unless it's officially announced.
	ENABLE_PRIVATE_SIGNATURE bool = false

	// generated the key by: openssl ecparam -name prime256v1 -genkey -noout -out privatekey.pem
	// only required when Private Signature is enabled
	// todo: replace with your own PrivateKey from privatekey.pem
	PRIVATE_KEY_PRIME_256 string = `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
	ACCOUNT_ID                   = "8806108"
)

// API请求地址, 不要带最后的/
const (
	//todo: replace with real URLs and HostName
	MARKET_URL string = "https://api.huobi.pro"
	TRADE_URL  string = "https://api.huobi.pro"
	HOST_NAME  string = "api.huobi.pro"
)

type ConfigTOML struct {
	Server struct {
		Listen             string         `required:"true" flagUsage:"服务监听地址"`
		Env                string         `default:"Pro" flagUsage:"服务运行时环境"`
		MaxHttpRequestBody int64          `default:"4" flagUsage:"最大允许的http请求body，单位M"`
		TimeLocation       *time.Location `flagUsage:"用于time.ParseInLocation"`
	}

	Auth struct {
		Secret  string            `flagUsage:"跳过鉴权的hack"`
		Account map[string]string `flagUsage:"复杂验证，apiKey=>apiSecret"`
	}

	Database struct {
		HostPort     string `required:"true" flagUsage:"数据库连接，eg：tcp(127.0.0.1:3306)"`
		UserPassword string `required:"true" flagUsage:"数据库账号密码"`
		DB           string `required:"true" flagUsage:"数据库"`
		Conn         struct {
			MaxLifeTime int `default:"600" flagUsage:"连接最长存活时间，单位s"`
			MaxIdle     int `default:"10" flagUsage:"最多空闲连接数"`
			MaxOpen     int `default:"80" flagUsage:"最多打开连接数"`
		}
	}

	Huobi struct {
		BuyRates  float64 `required:"true" default:0.0003 flagUsage:"买入手续费"`
		SaleRates float64 `required:"true" default:0.0002 flagUsage:"卖出手续费"`
	}

	Strategy struct {
		Type     string `required:"true" default:"floating" flagUsage:"floating浮动买入"`
		Floating struct {
			TotalAmount float64 `required:"true" default:0.0 flagUsage:"总金额仓位"`
			FloatRate   float64 `required:"true" default:0.01 flagUsage:"上下浮动的比例"`
			Depth       int     `required:"true" default:8 flagUsage:"最多下降和上升多少次"`
			Interval    int64   `required:"true" default:600 flagUsage:"间隔多少s再次启用策略"`
		}
	}

	KLineData struct {
		Symbol   string `required:"true" default:"market.btcusdt.kline.1min" default:"K线类型"`
		From     int64  `required:"true" default:1501174800 flagUsage:"socket获取数据开始时间戳"`
		To       int64  `required:"true" default:1560355200 flagUsage:"socket获取数据结束时间戳"`
		Duration int64  `required:"true" default:5 flagUsage:"每多少秒获取行情数据"`
	}

	Log struct {
		Type  string `default:"json" flagUsage:"日志格式，json|raw"`
		Level int    `default:"5" flagUsage:"日志级别：0 CRITICAL, 1 ERROR, 2 WARNING, 3 NOTICE, 4 INFO, 5 DEBUG"`
	} `flagUsage:"服务日志配置"`
}

func (c *ConfigTOML) IsProduction() bool {
	return strings.ToLower(c.Server.Env) == "pro"
}

var Config *ConfigTOML

func Init(tomlPath, args string) {
	var err error
	var loaders = []multiconfig.Loader{
		&multiconfig.TagLoader{},
		&multiconfig.TOMLLoader{Path: tomlPath},
		&multiconfig.EnvironmentLoader{},
		//&multiconfig.FlagLoader{RawArgs: args},
	}
	m := multiconfig.DefaultLoader{
		Loader:    multiconfig.MultiLoader(loaders...),
		Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
	}
	Config = new(ConfigTOML)
	m.MustLoad(&Config)

	Config.Server.TimeLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	util.PrettyPrint(Config)
}
