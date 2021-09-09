module github.com/crawlab-team/crawlab-plugin

go 1.15

replace (
	github.com/crawlab-team/crawlab-core => /Users/marvzhang/projects/crawlab-team/crawlab-core
	github.com/crawlab-team/crawlab-grpc => /Users/marvzhang/projects/crawlab-team/crawlab-grpc
)

require (
	github.com/apex/log v1.9.0
	github.com/crawlab-team/crawlab-core v0.6.0-beta.20210802.1344
	github.com/gin-gonic/gin v1.6.3
	go.mongodb.org/mongo-driver v1.4.5
)
