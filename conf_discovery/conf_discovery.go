package conf_discovery

import (
	"github.com/oikomi/FishChatServer2/common/constant"
)

type Heart struct {
}

func (h *Heart) Heart(protocol int) {
	if protocol == constant.ZOOKEEPER {

	}
}
