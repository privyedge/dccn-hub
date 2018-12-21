package handler

import (
	"fmt"
	"github.com/Ankr-network/refactor/app_dccn_k8s/proto"
	s "github.com/Ankr-network/refactor/gateway/proto/k8s"
	"github.com/Ankr-network/refactor/util"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/random"
	"io"
	"math/rand"
	"sync"
)


type K8sHandler struct{
	mu        sync.Mutex                      	// protects data
	dcstreams map[int]k8s.K8S_K8TaskService		//datacenterid => stream
}

func (p *K8sHandler) K8Task(stream s.K8S_K8TaskServer) error {
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
	return nil
}

func slelecDatacenterByID(s *K8sHandler, dcID int) k8s.K8S_K8TaskService {

	for key, stream := range s.dcstreams {
		if key == dcID {
			return stream
		}
	}

	return nil

}

func SelectFreeDatacenter(s *K8sHandler) k8s.K8S_K8TaskService {
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

func sendMessageToK8(stream k8s.K8S_K8TaskService, taskType string, taskid int64, name string, image string, replica int, extra string) bool {
	logStr := fmt.Sprintf("send sendMessageToK8 id %d name %s image %s replica %d  ", int(taskid), name, image, replica)
	util.WriteLog(logStr)
	if stream != nil {
		var message = k8s.Task{Type: taskType, Taskid: taskid, Name: name, Image: image, Replica: int64(replica), Extra: extra}
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
