package ankr_default

// RabbitQ name
const (
	MQDeployTask   = "topic.deploy.task"
	MQFeedbackTask = "topic.feedback.task"
	MQMail         = "topic.mail.handler"
)

// To do: Remove this line when usr mgr is ready
const Secret = "ed1605e17374bde6c68864d072c9f5c9"

const APIPort = 50051 // Default port for all gRPC request traffic

const HeartBeatInterval = 30 // Default interval for heartbeat

// Registry server name
const (
	TaskMgrRegistryServerName = "go.micro.srv.v1.task"
	UserMgrRegistryServerName = "go.micro.srv.v1.user"
	EmailRegistryServerName   = "go.micro.srv.v1.email"
	DcMgrRegistryServerName   = "go.micro.srv.v1.data.center"
	DCApiRegistryServerName   = "go.micro.srv.v1.app-dccn-api"
)
