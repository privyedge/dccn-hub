.PHONY: proto build test run deploy local clean echo grpc

GOPATH:=$(shell go env GOPATH)


APPS=$(wildcard */v1/*.proto)
RECURSIVE_PROTO = for app in $(APPS); do \
        @protoc -I. \
	    -I$(GOPATH)/src/github.com/Ankr-network/dccn-common/protos \
	    --micro_out=Mtask.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	    --go_out=Mevent.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	    $(subdir); \
        done


proto:
	@protoc -I. \
	--micro_out=Mcommon/event.proto=github.com/Ankr-network/dccn-common/protos/common,Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	--go_out=Mcommon/event.proto=github.com/Ankr-network/dccn-common/protos/common,Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	email/v1/micro/email.proto;

	@protoc -I. \
	--micro_out=Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	--go_out=Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	usermgr/v1/micro/usermgr.proto;

	@protoc -I. \
	-Icommon \
	--micro_out=Mcommon/common.proto=github.com/Ankr-network/dccn-common/protos/common,Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	--go_out=Mcommon/common.proto=github.com/Ankr-network/dccn-common/protos/common,Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	taskmgr/v1/micro/taskmgr.proto

	@protoc -I. \
	-Icommon \
	--micro_out=Mcommon/event.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	--go_out=Mcommon/event.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	dcmgr/v1/micro/dcmgr.proto

grpc:
	@cp email/v1/micro/email.proto email/v1/grpc
	@protoc -I. \
	--go_out=plugins=grpc,Mcommon/event.proto=github.com/Ankr-network/dccn-common/protos/common,Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	email/v1/grpc/email.proto;

	@cp usermgr/v1/micro/usermgr.proto usermgr/v1/grpc
	@protoc -I. \
	--go_out=plugins=grpc,Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	usermgr/v1/grpc/usermgr.proto;

	@cp taskmgr/v1/micro/taskmgr.proto taskmgr/v1/grpc
	@protoc -I. \
	-Icommon \
	--go_out=plugins=grpc,Mcommon/common.proto=github.com/Ankr-network/dccn-common/protos/common,Mcommon/error.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	taskmgr/v1/grpc/taskmgr.proto

	@cp dcmgr/v1/micro/dcmgr.proto dcmgr/v1/grpc/
	@protoc -I. \
	-Icommon \
	--go_out=plugins=grpc,Mcommon/event.proto=github.com/Ankr-network/dccn-common/protos/common:. \
	dcmgr/v1/grpc/dcmgr.proto
