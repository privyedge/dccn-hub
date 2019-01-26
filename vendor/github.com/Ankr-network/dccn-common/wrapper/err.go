package dccnwrapper

import (
	"fmt"
	"log"
	"path"
	"runtime"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func PbError(rsp **common_proto.Error, err error) {

	if *rsp == nil {
		*rsp = &common_proto.Error{}
	}
	(*rsp).Status = common_proto.Status_FAILURE
	(*rsp).Details = err.Error()
}

// for test user
func IsSuccess(info string, err *common_proto.Error) bool {
	if err == nil {
		return true
	}

	pathInfo := ""
	if _, file, line, ok := runtime.Caller(1); ok {
		pathInfo = fmt.Sprintf(" [%s:%v]: ", path.Base(file), line)
	}
	if err.Status == common_proto.Status_FAILURE {
		log.Fatal(pathInfo, info+": ", err.Details)
		return false
	}
	return true
}
