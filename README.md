# dccn-hub
Hub for Ankr DCCN

Ankr Hub consists of two microservices: k8s_handler and task_manager

## k8s_handler 

receive DC status and update to database
deliver task to DC when 
1.  receive new task notification from task_manager service  
2. interval query task for new task
receive task status from DC and update to database

## task_manager 

add task from cli to database
query task list from database and return to cli
notify k8s_handler when new task create


Database:
use MongoDB single instance

Collections:
task
user
datacenter

Ankr Hub connect with cli/k8s  by gRPC (k8s may use ZeroMQ as messaging)
