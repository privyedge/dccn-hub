package taskmgr

import (
	"fmt"
	"github.com/Ankr-network/dccn-hub/util"
	ankr_const "github.com/Ankr-network/dccn-rpc"
	"os"
)

type server struct {
	name string // protects data
}

func (s server) Handle(e util.Event) {
	logStr := fmt.Sprintf("this is handle message %s\n", e.Type)
	util.WriteLog(logStr)
	task := util.GetTask(e.TaskID)

	// add router
	logStr2 := fmt.Sprintf("task %d in DataCenter id %d  datacenter name %s  \n", task.ID, int(task.Datacenterid), task.Datacenter)
	util.WriteLog(logStr2)

	util.Send("datacenter_1", e)
}

func StartService() {
	if len(os.Args) == 3 {
		util.MongoDBHost = os.Args[1]

	}

	if len(os.Args) == 3 {
		util.RabbitMQHost = os.Args[2]
	}
	ss := server{}
	util.Receive(ankr_const.TaskManagerQueueName, ss)
}
