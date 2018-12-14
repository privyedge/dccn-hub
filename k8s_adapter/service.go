package k8s_adapter

import (
	"fmt"
	"github.com/Ankr-network/dccn-hub/util"
	pb "github.com/Ankr-network/dccn-rpc/protocol_new/k8s"
	"github.com/Ankr-network/dccn-rpc/server_rpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"math/rand"
	//"net"
	"os"
	"sync"
	"time"
)

const (
	port = ":50052"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	mu        sync.Mutex                      // protects data
	dcstreams map[int]pb.Dccnk8S_K8TaskServer //datacenterid => stream
}

func slelecDatacenterByName(s *server, dcName string) pb.Dccnk8S_K8TaskServer {
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

// RouteChat receives a stream of message/location pairs, and responds with a stream of all
// previous messages at each of those locations.
func (s *server) K8Task(stream pb.Dccnk8S_K8TaskServer) error {

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

		// var message = pb.Task{Type:"Task", Taskid: 13, Name:"docker_image", Extra:"extraxxx"}
		// if err := stream.Send(&message); err != nil {
		// 	return err
		// }

	}
}

func updateDataCenter(s *server, in *pb.K8SMessage, stream pb.Dccnk8S_K8TaskServer) {
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
	if datacenter.ID == 0 {
		fmt.Printf("datacenter not found\n")
	} else {

		fmt.Printf("processTaskStatus %d %s\n", taskid, status)
		//util.UpdateDataCenter(datacenter, int(datacenter.ID)
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

		time.Sleep(time.Second * 300)
	}
}

func (s server) Handle(e util.Event) {
	fmt.Printf("this is handle event type %s taskid %d \n", e.Type, e.TaskID)
	task := util.GetTask(e.TaskID)
	if e.Type == util.NewTaskEvent {
		var stream pb.Dccnk8S_K8TaskServer = nil
		if len(task.Datacenter) == 0 {
			fmt.Printf("find new free datacenter %d  \n", len(s.dcstreams))
			stream = SelectFreeDatacenter(&s)
		} else {
			fmt.Printf("find   datacenter  %s \n", task.Datacenter)
			dc := util.GetDatacenter(task.Datacenter)

			if dc.ID == 0 {
				fmt.Printf("add new task fail for datacenter name does not exist \n")

			} else {
				stream = slelecDatacenterByID(&s, int(dc.ID))
			}

		}

		if stream != nil {
			fmt.Printf("GetTaskNameAsTaskIDForK8s  id  %d name %s \n", task.ID, task.Name)
			var message = pb.Task{Type: "NewTask", Taskid: task.ID, Name: task.Uniquename, TaskType: task.Type, Image: task.Name, Extra: "nothing"}
			//fmt.Printf("new messsage for add task %s \n", message.Name)
			if err := stream.Send(&message); err != nil {
				fmt.Printf(">>>send add task message %s to data center failed \n", message.Name)
			} else {
				fmt.Printf(">>>send add task message %s to data center success \n", message.Name)
			}
		} else {
			fmt.Printf("no DataCenter available now\n")
		}
	}

	if e.Type == util.CancelTaskEvent {
		datacenter := s.dcstreams[int(task.Datacenterid)]
		if datacenter == nil {
			fmt.Printf("can not find datacenter \n")
			util.UpdateTask(int(task.ID), "cancelfailed", 0)

		} else {
			fmt.Printf("send cancel message to datacenter id  %d \n", int(task.Datacenterid))
			util.UpdateTask(int(task.ID), "cancelling", 0)

			if sendMessageToK8(datacenter, "CancelTask", task.ID, task.Uniquename, task.Name, task.Replica, "") == false {
				delete(s.dcstreams, int(task.Datacenterid))
			}

		}
	}

	if e.Type == util.UpdateTaskEvent {
		datacenter := s.dcstreams[int(task.Datacenterid)]
		if sendMessageToK8(datacenter, "UpdateTask", task.ID, task.Uniquename, e.Name, e.Replica, "") == false {
			delete(s.dcstreams, int(task.Datacenterid))
		} else {
			util.UpdateTaskReplica(e.TaskID, e.Replica)
			util.UpdateTaskImage(e.TaskID, e.Name)
		}
	}

}

func StartService() {
	if len(os.Args) == 2 {
		util.MongoDBHost = os.Args[1]

	}

	lis, s := server_rpc.Connect(port)
	ss := server{}
	go util.Receive(util.DataCenterName, &ss)
	ss.dcstreams = map[int]pb.Dccnk8S_K8TaskServer{}

	go heartbeat(&ss)

	pb.RegisterDccnk8SServer(s, &ss)
	// Register reflection service on gRPC server.

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
