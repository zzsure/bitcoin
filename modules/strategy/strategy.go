package strategy

import(
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/bitcoin/models"
	"gitlab.azbit.cn/web/bitcoin/modules/strategy/floating"
)

var logger = logging.MustGetLogger("modules/socket")

func Init() {
    strategys, err := models.GetAllStrategys()
    if err != nil {
        logger.Error("get all strategy error...", err)
    }
    for _, s := range strategys {
        if s.Name == "floating" {
            floating.Init(s)
        }
    }
}
