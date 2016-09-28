package external

const (
	// error
	ErrServer = 90001

	// gateway
	ReqAccessServerCMD                = 10001
	ResSelectAccessServerForClientCMD = 10002

	// msg_server
	ReqLoginCMD      = 20001
	ReqSendP2PMsgCMD = 20002
)
