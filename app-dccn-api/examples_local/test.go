package main

	
import "strings"
import "fmt"


func parseError(s1 string) string{
 index := strings.Index(s1, "detail")
 s2 := s1[index+9:]
 index2 := strings.Index(s2, "\"")
 s3 := s2[:index2]
 return s3
}

func main() {
 s1 := "rpc error: code = Unknown desc = {\"id\":\"\",\"code\":0,\"detail\":\"password does not match\",\"status\":\"\"}"


 fmt.Printf("%s\n", parseError(s1) )

}

