package ecode

import (
	"strconv"
)

type ecode uint32

func (e ecode) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

func (e ecode) String() string {
	return ecodeMessage[e]
}

func (e ecode) Uint32() uint32 {
	return uint32(e)
}

func To(i uint32) error {
	return ecode(i)
}

func From(e error) ecode {
	i, err := strconv.ParseInt(e.Error(), 10, 64)
	if err != nil {
		return ServerErr
	}
	if _, ok := ecodeMessage[ecode(i)]; !ok {
		return ServerErr
	}
	return ecode(i)
}

var (
	ecodeMessage = map[ecode]string{
		// api
		RequestErr: "request err",

		// common
		OK: "ok",
		// common
		ServerErr: "server error", // 服务器错误

		// gateway
		NoAccessServer: "no accessServer",

		// access
		NoToken:         "no token",
		CalcTokenFailed: "calc token failed",

		// register
		UserIsAlreadyExist: "user is already exist",

		// network
		NoData: "no data found",
	}
)
