package global

import (
	"github.com/oikomi/FishChatServer2/libnet"
)

type SessionMap map[int64]*libnet.Session

var GSessions SessionMap

func init() {
	GSessions = make(SessionMap)
}
