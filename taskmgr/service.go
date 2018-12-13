package taskmgr

import (
	"fmt"
	util "github.com/Ankr-network/dccn-hub/util"
	pb "github.com/Ankr-network/dccn-rpc/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Ankr-network/dccn-hub/util"
	pb "github.com/Ankr-network/dccn-rpc/protocol"
	server_rpc "github.com/Ankr-network/dccn-rpc/server_rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct {
	mu        sync.Mutex                      // protects data
	dcstreams map[int]pb.Dccncli_K8TaskServer //datacenterid => stream
}

func (s *server) TaskDetail(ctx context.Context, in *pb.TaskDetailRequest) (*pb.TaskDetailResponse, error) {
	fmt.Printf("received task detail request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		fmt.Printf("add new task fail for user token error \n")

		return &pb.TaskDetailResponse{Body: ""}, nil
	}
	task := util.GetTask(int(in.Taskid))
	return &pb.TaskDetailResponse{Body: task.URL}, nil

}

func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	fmt.Printf("received add task request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		fmt.Printf("add new task fail for user token error \n")
		return &pb.AddTaskResponse{Status: "Failure", Taskid: -1}, nil
	} else {

		// check datacenter name valid
		var stream pb.Dccncli_K8TaskServer = nil
		if len(in.Datacenter) == 0 {
			stream = SelectFreeDatacenter(s)
		} else {
			dc := util.GetDatacenter(in.Datacenter)

			if dc.ID == 0 {
				fmt.Printf("add new task fail for datacenter name does not exist \n")
				return &pb.AddTaskResponse{Status: "Failure", Taskid: -1}, nil

			}
			stream = slelecDatacenterByID(s, int(dc.ID))
		}
		//end check datacenter name

		task := util.Task{Name: in.Name, Datacenter: in.Datacenter, Userid: user.ID, Type: in.Type}
		id := util.AddTask(task)
		task.ID = id

		tastName := util.GetTaskNameAsTaskIDForK8s(task)
		util.UpdateTaskUnqueName(int(id), tastName)

		if in.Replica <= 0 || in.Replica > 100 {
			fmt.Printf("replica is eror %d use default 1  \n", in.Replica)
			in.Replica = 1
		}

		util.UpdateTaskReplica(int(id), int(in.Replica))

		if stream != nil {
			fmt.Printf("GetTaskNameAsTaskIDForK8s  id  %d name %s \n", task.ID, task.Name)
			var message = pb.Task{Type: "NewTask", Taskid: id, Name: tastName, TaskType: task.Type, Image: task.Name, Extra: "nothing"}
			//fmt.Printf("new messsage for add task %s \n", message.Name)
			if err := stream.Send(&message); err != nil {
				fmt.Printf(">>>send add task message %s to data center failed \n", message.Name)
			} else {
				fmt.Printf(">>>send add task message %s to data center success \n", message.Name)
			}
		} else {
			fmt.Printf("no DataCenter available now\n")
		}

		return &pb.AddTaskResponse{Status: "Success", Taskid: id}, nil
	}

}

func slelecDatacenterByName(s *server, dcName string) pb.Dccncli_K8TaskServer {
	dc := util.GetDatacenter(dcName)

	if dc.ID == 0 {
		return nil
	}

	for key, stream := range s.dcstreams {
		if key == int(dc.ID) {
			return stream
		}
	}

	return nil

}

func slelecDatacenterByID(s *server, dcID int) pb.Dccncli_K8TaskServer {

	for key, stream := range s.dcstreams {
		if key == dcID {
			return stream
		}
	}

	return nil

}

func SelectFreeDatacenter(s *server) pb.Dccncli_K8TaskServer {
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

func sendMessageToK8(stream pb.Dccncli_K8TaskServer, taskType string, taskid int64, name string, image string, replica int, extra string) bool {
	fmt.Printf("send sendMessageToK8 id %d name %s image %s replica %d   \n", int(taskid), name, image, replica)
	if stream != nil {
		var message = pb.Task{Type: taskType, Taskid: taskid, Name: name, Image: image, Replica: int64(replica), Extra: extra}
		if err := stream.Send(&message); err != nil {
			fmt.Printf("send message to data center failed \n")
			return false
		} else {
			fmt.Printf("send message to data center successfully \n")
			return true
		}
	}

	return false
}

func (s *server) TaskList(ctx context.Context, in *pb.TaskListRequest) (*pb.TaskListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)
	fmt.Printf("task list reqeust \n")

	if user.ID == 0 {
		fmt.Printf("task list reqeust fail for user token error\n")
		return &pb.TaskListResponse{}, nil
	} else {
		tasks := util.TaskList(int(user.ID))

		dcs := util.GetDatacentersMap()

		var taskList []*pb.TaskInfo
		for i := range tasks {
			task := tasks[i]
			taskInfo := &pb.TaskInfo{}
			taskInfo.Taskid = task.ID
			taskInfo.Taskname = task.Name
			taskInfo.Status = task.Status
			taskInfo.Replica = int64(task.Replica)
			taskInfo.Datacenter = dcs[task.Datacenterid]
			if len(taskInfo.Datacenter) == 0 {
				taskInfo.Datacenter = task.Datacenter // for user assign datacenter name but not startsuccess
			}
			taskList = append(taskList, taskInfo)
			//fmt.Printf("task id : %d %s status %s \n", task.ID,task.Name, task.Status)
		}

		return &pb.TaskListResponse{Tasksinfo: taskList}, nil
	}

}

func (s *server) DataCenterList(ctx context.Context, in *pb.DataCenterListRequest) (*pb.DataCenterListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)
	fmt.Printf("task list reqeust \n")

	if user.ID == 0 {
		fmt.Printf("task list reqeust fail for user token error\n")
		return &pb.DataCenterListResponse{}, nil
	} else {
		dataCenters := util.DataCeterList()

		var dcList []*pb.DataCenterInfo
		for i := range dataCenters {
			dataCenter := dataCenters[i]
			dcInfo := &pb.DataCenterInfo{}
			dcInfo.Id = dataCenter.ID
			dcInfo.Name = dataCenter.Name
			dcList = append(dcList, dcInfo)
			//fmt.Printf("task id : %d %s status %s \n", task.ID,task.Name, task.Status)
		}

		return &pb.DataCenterListResponse{DcList: dcList}, nil
	}

}

func (s *server) CancelTask(ctx context.Context, in *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
	fmt.Printf("received cancel task request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		fmt.Printf("can not find task\n")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	if user.ID == 0 {
		fmt.Printf("cancel task fail for user token error\n")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	if task.Userid != user.ID {
		fmt.Printf("task uid != user id \n")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	fmt.Printf("task %d in DataCenter %d \n", task.ID, int(task.Datacenterid))

	//sendMessageToK8(taskType string, taskid int64, name string, extra string)
	// todo test this function
	datacenter := s.dcstreams[int(task.Datacenterid)]
	if datacenter == nil {
		fmt.Printf("can not find datacenter \n")
		util.UpdateTask(int(in.Taskid), "cancelfailed", 0)
		return &pb.CancelTaskResponse{Status: "Failure"}, nil

	} else {
		fmt.Printf("send cancel message to datacenter id  %d \n", int(task.Datacenterid))
		util.UpdateTask(int(in.Taskid), "cancelling", 0)
		if sendMessageToK8(datacenter, "CancelTask", in.Taskid, task.Uniquename, task.Name, task.Replica, "") == false {
			delete(s.dcstreams, int(task.Datacenterid))
		}
		return &pb.CancelTaskResponse{Status: "Success"}, nil
	}

}

func (s *server) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	fmt.Printf("received update task request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		fmt.Printf("can not find task\n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if user.ID == 0 {
		fmt.Printf("cancel task fail for user token error\n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if task.Userid != user.ID {
		fmt.Printf("task uid != user id \n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if len(task.Uniquename) == 0 {
		fmt.Printf("task does not have Uniquename in mongodb \n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	fmt.Printf("task %d in DataCenter %d \n", task.ID, int(task.Datacenterid))

	datacenter := s.dcstreams[int(task.Datacenterid)]
	if datacenter == nil {
		fmt.Printf("can not find datacenter \n")
		util.UpdateTask(int(in.Taskid), "updateFailed", 0)
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	} else {

		//check replica is valid
		if task.Replica == 0 { // support previous
			task.Replica = 1
		}

		if in.Replica <= 0 || in.Replica > 100 {
			fmt.Printf("replica is eror %d \n", in.Replica)
			in.Replica = int64(task.Replica)
		}

		if len(in.Name) == 0 {
			in.Name = task.Name
		}

		fmt.Printf("send update message to datacenter id  %d  replica  %d  image : %s\n", int(task.Datacenterid), int(in.Replica), task.Name)

		util.UpdateTaskReplica(int(in.Taskid), int(in.Replica))
		util.UpdateTaskImage(int(in.Taskid), in.Name)
		util.UpdateTask(int(in.Taskid), "updating", 0)

		// if they are same use 0 as default value
		if int(in.Replica) == task.Replica {
			in.Replica = 0
		}

		if in.Name == task.Name {
			in.Name = ""
		}

		if sendMessageToK8(datacenter, "UpdateTask", in.Taskid, task.Uniquename, in.Name, int(in.Replica), "") == false {
			delete(s.dcstreams, int(task.Datacenterid))
		}
		return &pb.UpdateTaskResponse{Status: "Success"}, nil
	}

}

func (s *server) K8ReportStatus(ctx context.Context, in *pb.ReportRequest) (*pb.ReportResponse, error) {
	fmt.Printf("received K8ReportStatus request %s\n", in.Report)
	datacenter := util.GetDataCenter(in.Name)
	if datacenter.ID == 0 {
		datacenter := util.DataCenter{Name: in.Name, Report: in.Report}
		id := util.AddDataCenter(datacenter)
		fmt.Printf("insert new  DataCenter id %d \n", id)
	} else {
		datacenter2 := util.DataCenter{Name: in.Name, Report: in.Report}
		util.UpdateDataCenter(datacenter2, int(datacenter.ID))
		fmt.Printf("update  DataCenter id %d \n", datacenter.ID)

	}

	return &pb.ReportResponse{Status: "Success"}, nil

}

//receive message from DataCenter by stream, two type of messages: HeartBeat Task
func (s *server) K8Task(stream pb.Dccncli_K8TaskServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("<<<received  k8s  task : id %d  name:  %s  status: %s \n", in.Taskid, in.Taskname, in.Status)
		s.mu.Lock()
		if in.Type == "HeartBeat" {
			updateDataCenter(s, in, stream)
			fmt.Printf("received  HeartBeat  : datacenter name:  %s report :  %s \n", in.Datacenter, in.Report)
		} else {
			taskId := util.GetTaskIDFromTaskNameForK8s(in.Taskname)
			fmt.Printf("<<<received task  id : %d  status: %s  datacenter : %s \n", taskId, in.Status, in.Datacenter)

			processTaskStatus(taskId, in.Status, in.Datacenter, in.Url)
		}

		s.mu.Unlock()

	}
}

//deal with HeartBeat message from DataCenter
func updateDataCenter(s *server, in *pb.K8SMessage, stream pb.Dccncli_K8TaskServer) {
	datacenter := util.GetDataCenter(in.Datacenter)
	if datacenter.ID == 0 {
		datacenter := util.DataCenter{Name: in.Datacenter, Report: in.Report}
		id := util.AddDataCenter(datacenter)
		fmt.Printf("insert new  DataCenter id %d \n", id)
	} else {
		datacenter2 := util.DataCenter{Name: in.Datacenter, Report: in.Report}
		util.UpdateDataCenter(datacenter2, int(datacenter.ID))
		fmt.Printf("update  DataCenter id %d \n", datacenter.ID)

	}

	datacenter = util.GetDataCenter(in.Datacenter)
	s.dcstreams[int(datacenter.ID)] = stream

}

func processTaskStatus(taskid int64, status string, dcName string, url string) {
	datacenter := util.GetDataCenter(dcName)
	// fmt.Printf the error message for now. They are supposed to be throw as an exception in the future.
	if datacenter.ID == 0 {
		fmt.Printf("datacenter not found\n")
	} else if taskid <= 0 {
		fmt.Printf("Task id is supposed to be larger than zero\n")
	} else {

		fmt.Printf("processTaskStatus %d %s\n", taskid, status)
		if status == "StartSuccess" {
			if len(url) > 0 {
				util.UpdateTaskURL(int(taskid), url)
			}
			util.UpdateTask(int(taskid), "running", int(datacenter.ID))
		}

		if status == "StartFailure" {
			util.UpdateTask(int(taskid), "startFailed", int(datacenter.ID))
		}

		if status == "UpdateSuccess" {
			util.UpdateTask(int(taskid), "running", int(datacenter.ID))
		}

		if status == "UpdateFailure" {
			util.UpdateTask(int(taskid), "updateFailed", int(datacenter.ID))
		}

		if status == "UpdateFailure" {
			util.UpdateTask(int(taskid), "updateFailed", int(datacenter.ID))
		}

		if status == "Cancelled" {
			util.UpdateTask(int(taskid), "cancelled", int(datacenter.ID))
		}

		if status == "Done" {
			util.UpdateTask(int(taskid), "done", int(datacenter.ID))
		}

	}
}

func heartbeat(s *server) {
	for {
		fmt.Printf("send HeartBeat to %d DataCenters \n", len(s.dcstreams))
		for key, stream := range s.dcstreams {
			if sendMessageToK8(stream, "HeartBeat", -1, "", "", 0, "") == false {
				delete(s.dcstreams, key)
			}
		}

		time.Sleep(time.Second * 30)
	}
}

func Serve() {
	if len(os.Args) == 2 {
		util.MongoDBHost = os.Args[1]

	}

	util.WriteLog("this is a test")

	lis, s := server_rpc.Connect(port)
	ss := server{}
	ss.dcstreams = map[int]pb.Dccncli_K8TaskServer{}

	go heartbeat(&ss)

	pb.RegisterDccncliServer(s, &ss)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
