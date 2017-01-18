package external

const (
	// error
	ErrServerCMD = 90001

	// gateway
	ReqAccessServerCMD = 10001
	// ResSelectAccessServerForClientCMD = 10002

	// acess
	LoginCMD           = 20001
	PingCMD            = 20002
	SendP2PMsgCMD      = 20003
	AcceptP2PMsgAckCMD = 20004
	SendGroupMsgCMD    = 20005
	SyncMsgCMD         = 20006
	NotifyCMD          = 20007
)
