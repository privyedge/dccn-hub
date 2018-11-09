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

## Workflow

* Users define desired infrastructure and workloads to run on the remote infrastructure, and how workloads can connect to one another.
  * Desired lifetime of resources is expressed via collateral requirements.
* Orders are generated from the deployment requirement.
* Datacenters bid on an open orderbook.
* The bid with lowest price (more heuristics later) gets matched with order to create a contract.
* Once a contract is reached, workloads and topology are delivered to the datacenter(s).
* Datacenter(s) deploy workloads and allow connectivity as specified by the tenant.
* If a datacenter fails to maintain contract, collateral is transferred to tenant, and a new order is created for the desired resources.
* A tenant can close any active deployment at any time
