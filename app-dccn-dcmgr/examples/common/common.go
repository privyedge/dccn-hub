package testCommon

import (
	"reflect"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func IsEqual(origin, dst *common_proto.Task) bool {
	ok := origin.Id == dst.Id &&
		origin.UserId == dst.UserId &&
		origin.Type == dst.Type &&
		origin.Name == dst.Name &&
		origin.Image == dst.Image &&
		origin.Replica == dst.Replica &&
		origin.DataCenter == dst.DataCenter &&
		origin.DataCenterId == dst.DataCenterId &&
		origin.UniqueName == dst.UniqueName &&
		origin.Url == dst.Url &&
		origin.Hidden == dst.Hidden &&
		origin.Uptime == dst.Uptime &&
		origin.CreationDate == dst.CreationDate
	if origin.Extra != nil && dst.Extra != nil {
		ok = ok && reflect.DeepEqual(origin.Extra, dst.Extra)
	}
	return ok
}

func MockTasks() []common_proto.Task {
	// id generated by uuid
	return []common_proto.Task{
		{
			UserId:       "1",
			Name:         "task01",
			Type:         "web",
			Image:        "nginx",
			Replica:      2,
			DataCenter:   "dc01",
			DataCenterId: "1",
		},
		{
			UserId:       "1",
			Name:         "task02",
			Type:         "web",
			Image:        "nginx",
			Replica:      2,
			DataCenter:   "dc02",
			DataCenterId: "1",
		},
		{
			UserId:       "2",
			Name:         "task01",
			Type:         "web",
			Image:        "nginx",
			Replica:      2,
			DataCenter:   "dc01",
			DataCenterId: "1",
		},
	}
}
