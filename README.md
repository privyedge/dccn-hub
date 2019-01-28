# Start All Micro Services
    $docker-compose up -d datastore consul rabbit
    $sleep 5s
    $docker-compose up -d taskmgr

$ Run test
    $docker-compose up -d taskmgr-client
    $docker-compose stop taskmgr-client
    $docker-compose rm taskmgr-client
    
