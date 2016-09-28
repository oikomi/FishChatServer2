package external

const (
	// error
	ErrServer = 90001

	// gateway
	ReqMsgServerCMD                = 10001
	ResSelectMsgServerForClientCMD = 10002

	// msg_server
	ReqLoginCMD      = 20001
	ReqSendP2PMsgCMD = 20002
)
