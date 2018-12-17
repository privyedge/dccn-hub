# Micro service taskMgr

add task from cli to database query task list from database and return to cli notify k8s_handler when new task create
Database: use MongoDB single instance
Collections: task user datacenter
Ankr Hub connect with cli/k8s by gRPC (k8s may use ZeroMQ as messaging)
