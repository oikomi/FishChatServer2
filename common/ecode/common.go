package ecode

const (
	//api
	RequestErr ecode = 10001

	// common error code
	OK ecode = 0
	// server
	ServerErr ecode = 90001

	//gateway
	NoAccessServer ecode = 92001

	// access_server
	NoToken         ecode = 93001
	CalcTokenFailed ecode = 93002

	// network
	NoData ecode = 91001
	//
)
