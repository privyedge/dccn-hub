package main

import (
	"log"
	"io"
	"net"
  "fmt"
	"sync"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "dccn-hub/protocol"
	"google.golang.org/grpc/reflection"
	util "dccn-hub/util"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{
	 mu  sync.Mutex // protects data
	 dcstream  pb.Dccncli_K8TaskServer  //datacenter stream
}


func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
        fmt.Printf("received add task request\n")
				token := in.Usertoken
				user := util.GetUser(token)

        if user.ID == 0 {
					fmt.Printf("add new task fail for user token error\n")
					return &pb.AddTaskResponse{Status:"Failure", Taskid: -1}, nil
				}else{
					task := util.Task{Name: in.Name, Region: in.Region, Zone: in.Zone, Userid: user.ID}
					id := util.AddTask(task)
					fmt.Printf("add new task ID : %d\n", id)

					if s.dcstream != nil {
						var message = pb.Task{Taskid: 255, Name:"docker_image", Extra:"extraxxx"}
						if err := s.dcstream.Send(&message); err != nil {
							  fmt.Printf("send message to data center failed \n")
						}else{
							  fmt.Printf("send message to data center successfully \n")
						}
					}

           return &pb.AddTaskResponse{Status:"Success", Taskid: id}, nil
				}

}


func (s *server) TaskList(ctx context.Context, in *pb.TaskListRequest) (*pb.TaskListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		fmt.Printf("task list reqeust fail for user token error\n")
		return &pb.TaskListResponse{}, nil
	}else{
		tasks := util.TaskList(int(user.ID))

    var taskList []*pb.TaskInfo
		for i := range tasks {
				task := tasks[i]
				taskInfo := &pb.TaskInfo{}
				taskInfo.Taskid = task.ID
				taskInfo.Taskname = task.Name
				taskInfo.Status = task.Status
				taskList = append(taskList, taskInfo)
				fmt.Printf("task id : %d %s status %s \n", task.ID,task.Name, task.Status)
     }

 		return &pb.TaskListResponse{Tasksinfo: taskList}, nil
	}

}

func (s *server) CancelTask(ctx context.Context, in *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
        fmt.Printf("received add task request\n")
				token := in.Usertoken
				user := util.GetUser(token)

				task := util.GetTask(int(in.Taskid))

         if task.ID == 0 {
					 fmt.Printf("can not find task\n")
					 return &pb.CancelTaskResponse{Status:"Failure"}, nil
				 }

        if user.ID == 0 {
					fmt.Printf("add new task fail for user token error\n")
					return &pb.CancelTaskResponse{Status:"Failure"}, nil
				}

				if task.Userid != user.ID {
					fmt.Printf("task uid != user id \n")
					return &pb.CancelTaskResponse{Status:"Failure"}, nil
				}


				return &pb.CancelTaskResponse{Status:"Success"}, nil


}



func (s *server) K8ReportStatus(ctx context.Context, in *pb.ReportRequest) (*pb.ReportResponse, error) {
         fmt.Printf("received K8ReportStatus request %s\n", in.Report)
				 datacenter := util.GetDataCenter(in.Name)
				 if datacenter.ID == 0{
					 	datacenter := util.DataCenter{Name: in.Name, Report:in.Report, Host: in.Host, Port: in.Port}
						id := util.AddDataCenter(datacenter)
						fmt.Printf("insert new  DataCenter id %d \n", id)
				 }else{
					    datacenter2 := util.DataCenter{Name: in.Name, Report:in.Report, Host: in.Host, Port: in.Port}
              util.UpdateDataCenter(datacenter2, int(datacenter.ID))
							fmt.Printf("update  DataCenter id %d \n", datacenter.ID)

				 }

         return &pb.ReportResponse{Status:"Success"}, nil


}


func (s *server) K8QueryTask(ctx context.Context, in *pb.QueryTaskRequest) (*pb.QueryTaskResponse, error) {
         fmt.Printf("received K8QueryTask request\n")
				 datacenter := util.GetDataCenter(in.Name)
				 if datacenter.ID == 0 {
					 fmt.Printf("datacenter not found\n")
					 return &pb.QueryTaskResponse{}, nil
				 }else{
              //util.UpdateDataCenter(datacenter, int(datacenter.ID)

				 }

				 task := util.GetNewTask()
				 if task.ID == 0 {
					 fmt.Printf("No new task\n")
					 return &pb.QueryTaskResponse{}, nil
				 }else{
					  fmt.Printf("GetNewTask %d\n", task.ID)
					  util.UpdateTask(int(task.ID), "running", int(datacenter.ID))
            return &pb.QueryTaskResponse{Taskid:task.ID, Name:task.Name, Extra:""}, nil
				 }
}


// RouteChat receives a stream of message/location pairs, and responds with a stream of all
// previous messages at each of those locations.
func (s *server) K8Task(stream pb.Dccncli_K8TaskServer) error {
	s.dcstream = stream
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		taskid := in.Taskid

		s.mu.Lock()
		if taskid == -1 {
				fmt.Printf("received error task status  : %d %s \n", in.Taskid, in.Status)
		}else{
        fmt.Printf("received task status  : %d %s \n", in.Taskid, in.Status)
		}

		s.mu.Unlock()

     var message = pb.Task{Taskid: 13, Name:"docker_image", Extra:"extraxxx"}
			if err := stream.Send(&message); err != nil {
				return err
			}

	}
}


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDccncliServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
