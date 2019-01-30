package ankr_default

import "errors"

// List execution errors
var (
	ErrDataCenterNotExist        = errors.New("dataCenter does not exist")
	ErrUserNotExist              = errors.New("token error, can not find user")
	ErrTaskNotExist              = errors.New("task does not exist")
	ErrUserNotOwn                = errors.New("user does not own this task")
	ErrUpdateFailed              = errors.New("task can not be updated")
	ErrUserAlreadyExist          = errors.New("user already existed")
	ErrPasswordError             = errors.New("password does not match")
	ErrNamePasswordEmpty         = errors.New("name or Password is empty")
	ErrStatusNotSupportOperation = errors.New("current status not support operation")
	ErrReplicaTooMany            = errors.New("replica too many")
	ErrUnknown                   = errors.New("unknown operation or code")
	ErrSyncTaskInfo              = errors.New("sync task info error")
	ErrPublish                   = errors.New("mq publish message error")
	ErrConnection                = errors.New("connection error")
	ErrNoAvailableDataCenter     = errors.New("no available data center")
)
