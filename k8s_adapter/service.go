package k8s_adapter

import (
	"fmt"
	ankr_const "github.com/Ankr-network/dccn-common"
	pb "github.com/Ankr-network/dccn-common/protocol/k8s"
	server_rpc "github.com/Ankr-network/dccn-common/server_rpc"
	"github.com/Ankr-network/dccn-hub/util"
	"google.golang.org/grpc/reflection"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	port = ":" + ankr_const.DefaultPort
)

type server struct {
	mu        sync.Mutex                      // protects data
	dcstreams map[int]pb.Dccnk8S_K8TaskServer //datacenterid => stream
}

func slelecDatacenterByID(s *server, dcID int) pb.Dccnk8S_K8TaskServer {

	for key, stream := range s.dcstreams {
		if key == dcID {
			return stream
		}
	}

	return nil

}

func SelectFreeDatacenter(s *server) pb.Dccnk8S_K8TaskServer {
	keys := []int{}
	for key, _ := range s.dcstreams {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil
	}

	index := rand.Intn(len(keys))
	key := keys[index]
	return s.dcstreams[key]

}

func sendMessageToK8(stream pb.Dccnk8S_K8TaskServer, taskType string, taskid int64, name string, image string, replica int, extra string) bool {
	logStr := fmt.Sprintf("send sendMessageToK8 id %d name %s image %s replica %d  ", int(taskid), name, image, replica)
	util.WriteLog(logStr)
	if stream != nil {
		var message = pb.Task{Type: taskType, Taskid: taskid, Name: name, Image: image, Replica: int64(replica), Extra: extra}
		if err := stream.Send(&message); err != nil {
			util.WriteLog("send message to data center failed")
			return false
		} else {
			util.WriteLog("send message to data center successfully")
			return true
		}
	}

	return false
}

func (s *server) K8Task(stream pb.Dccnk8S_K8TaskServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		logStr := fmt.Sprintf("<<<received  k8s  task : id %d  name:  %s  status: %s", in.Taskid, in.Taskname, in.Status)
		util.WriteLog(logStr)

		s.mu.Lock()
		if in.Type == "HeartBeat" {
			updateDataCenter(s, in, stream)
			logStr := fmt.Sprintf("received  HeartBeat  : datacenter name:  %s report :  %s", in.Datacenter, in.Report)
			util.WriteLog(logStr)
		} else {
			taskId := util.GetTaskIDFromTaskNameForK8s(in.Taskname)
			logStr := fmt.Sprintf("<<<received task  id : %d  status: %s  datacenter : %s", taskId, in.Status, in.Datacenter)
			util.WriteLog(logStr)

			processTaskStatus(taskId, in.Status, in.Datacenter, in.Url)
		}

		s.mu.Unlock()

	}
}

func updateDataCenter(s *server, in *pb.K8SMessage, stream pb.Dccnk8S_K8TaskServer) {
	datacenter := util.GetDataCenter(in.Datacenter)
	if datacenter.ID == 0 {
		datacenter := util.DataCenter{Name: in.Datacenter, Report: in.Report}
		id := util.AddDataCenter(datacenter)
		logStr := fmt.Sprintf("insert new  DataCenter id %d", id)
		util.WriteLog(logStr)
	} else {
		datacenter2 := util.DataCenter{Name: in.Datacenter, Report: in.Report}
		util.UpdateDataCenter(datacenter2, int(datacenter.ID))
		logStr := fmt.Sprintf("update  DataCenter id %d", datacenter.ID)
		util.WriteLog(logStr)

	}

	datacenter = util.GetDataCenter(in.Datacenter)
	s.dcstreams[int(datacenter.ID)] = stream

}

func processTaskStatus(taskid int64, status string, dcName string, url string) {
	datacenter := util.GetDataCenter(dcName)
	if datacenter.ID == 0 {
		util.WriteLog("datacenter not found")
	} else {
		logStr := fmt.Sprintf("processTaskStatus %d %s", taskid, status)
		util.WriteLog(logStr)

		if status == ankr_const.DataCenterTaskStartSuccess {
			if len(url) > 0 {
				util.UpdateTaskURL(int(taskid), url)
			}
			util.UpdateTask(int(taskid), ankr_const.TaskStatusRunning, int(datacenter.ID))
		}

		if status == ankr_const.DataCenterTaskStartFailure {
			if len(url) > 0 {
				util.UpdateTaskURL(int(taskid), url)
			}
			util.UpdateTask(int(taskid), ankr_const.TaskStatusStartFailed, int(datacenter.ID))
		}

		if status == ankr_const.DataCenterTaskUpdateFailure {
			util.UpdateTask(int(taskid), ankr_const.TaskStatusStartFailed, int(datacenter.ID))
		}

		if status == ankr_const.DataCenterTaskUpdateSuccess {
			util.UpdateTask(int(taskid), ankr_const.TaskStatusRunning, int(datacenter.ID))
		}

		if status == ankr_const.DataCenterTaskCancelled {
			util.UpdateTask(int(taskid), ankr_const.TaskStatusCancelled, int(datacenter.ID))
		}

		if status == ankr_const.DataCenterTaskDone {
			util.UpdateTask(int(taskid), ankr_const.TaskStatusDone, int(datacenter.ID))
		}

	}
}

func heartbeat(s *server) {
	for {
		logStr := fmt.Sprintf("send HeartBeat to %d DataCenters ", len(s.dcstreams))
		util.WriteLog(logStr)

		successList := []int64{}

		for key, stream := range s.dcstreams {
			if sendMessageToK8(stream, "HeartBeat", -1, "", "", 0, "") == false {
				delete(s.dcstreams, key)
			} else {
				successList = append(successList, int64(key))
			}
		}

		util.UpdataDataCentersStatus(successList)

		time.Sleep(time.Second * 30)
	}
}

func (s server) Handle(e util.Event) {
	logStr := fmt.Sprintf("this is handle event type %s taskid %d ", e.Type, e.TaskID)
	util.WriteLog(logStr)
	task := util.GetTask(e.TaskID)
	if e.Type == ankr_const.NewTaskEvent {
		var stream pb.Dccnk8S_K8TaskServer = nil
		if len(task.Datacenter) == 0 {
			logStr := fmt.Sprintf("find new free datacenter %d ", len(s.dcstreams))
			util.WriteLog(logStr)
			stream = SelectFreeDatacenter(&s)
		} else {
			logStr := fmt.Sprintf("find   datacenter  %s", task.Datacenter)
			util.WriteLog(logStr)
			dc := util.GetDatacenter(task.Datacenter)

			if dc.ID == 0 {
				util.WriteLog("add new task fail for datacenter name does not exist")

			} else {
				stream = slelecDatacenterByID(&s, int(dc.ID))
			}

		}

		if stream != nil {
			logStr := fmt.Sprintf("GetTaskNameAsTaskIDForK8s  id  %d name %s", task.ID, task.Name)
			util.WriteLog(logStr)
			var message = pb.Task{Type: "NewTask", Taskid: task.ID, Name: task.Uniquename, TaskType: task.Type, Replica: int64(task.Replica), Image: task.Name, Extra: "nothing"}
			//util.WriteLog("new messsage for add task %s", message.Name)
			if err := stream.Send(&message); err != nil {
				logStr := fmt.Sprintf(">>>send add task message %s to data center failed", message.Name)
				util.WriteLog(logStr)
			} else {
				logStr := fmt.Sprintf(">>>send add task message %s to data center success", message.Name)
				util.WriteLog(logStr)
			}
		} else {
			util.WriteLog("no DataCenter available now")
		}
	}

	if e.Type == ankr_const.CancelTaskEvent {
		datacenter := s.dcstreams[int(task.Datacenterid)]
		if datacenter == nil {
			util.WriteLog("can not find datacenter")
			util.UpdateTask(int(task.ID), "cancelfailed", 0)

		} else {
			logStr := fmt.Sprintf("send cancel message to datacenter id  %d", int(task.Datacenterid))
			util.WriteLog(logStr)
			util.UpdateTask(int(task.ID), "cancelling", 0)

			if sendMessageToK8(datacenter, "CancelTask", task.ID, task.Uniquename, task.Name, task.Replica, "") == false {
				delete(s.dcstreams, int(task.Datacenterid))
			}

		}
	}

	if e.Type == ankr_const.UpdateTaskEvent {
		datacenter := s.dcstreams[int(task.Datacenterid)]
		if sendMessageToK8(datacenter, "UpdateTask", task.ID, task.Uniquename, e.Name, e.Replica, "") == false {
			delete(s.dcstreams, int(task.Datacenterid))
		} else {
			if e.Replica != 0 {
				util.UpdateTaskReplica(e.TaskID, e.Replica)
			}
			if len(e.Name) > 0 {
				util.UpdateTaskImage(e.TaskID, e.Name)
			}

		}
	}

}

func StartService() {
	if len(os.Args) == 3 {
		util.MongoDBHost = os.Args[1]

	}

	if len(os.Args) == 3 {
		util.RabbitMQHost = os.Args[2]
	}

	lis, s := server_rpc.Connect(port)
	ss := server{}
	go util.Receive(ankr_const.DataCenterName, &ss)
	ss.dcstreams = map[int]pb.Dccnk8S_K8TaskServer{}

	go heartbeat(&ss)

	pb.RegisterDccnk8SServer(s, &ss)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		logStr := fmt.Sprintf("failed to serve: %v", err)
		util.WriteLog(logStr)
	}

}
