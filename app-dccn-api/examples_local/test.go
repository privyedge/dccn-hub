package main

import (
	//"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"

	//	"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
	"log"
	//	"time"

	//	"google.golang.org/grpc/metadata"

	//	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)



func main() {

	m1 := make(map[string]int, 0)
	m1["apple"] = 10;

	log.Printf("m1 is %d", m1["apple"])

	// Using a composite literal to initialize a map.
	m2 := map[string]int{}
	m2["apple"] = 20;

	log.Printf("m2 is %d", m2["apple"])


	m3 := new(int)
	(*m3) = 10;

	log.Printf("m3 is %d", (*m3))


	log.Println("END")
}
