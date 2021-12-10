module github.com/crawlab-team/crawlab-plugin

go 1.16

replace github.com/crawlab-team/crawlab-core => ../crawlab-core

require (
	github.com/apex/log v1.9.0
	github.com/aws/aws-sdk-go v1.34.28 // indirect
	github.com/crawlab-team/crawlab-core v0.6.0-beta.20211009.1458
	github.com/crawlab-team/crawlab-grpc v0.6.0-beta.20211009.1455
	github.com/crawlab-team/go-trace v0.1.0
	github.com/gin-gonic/gin v1.7.1
	github.com/gobuffalo/genny v0.1.1 // indirect
	github.com/gobuffalo/gogen v0.1.1 // indirect
	github.com/karrick/godirwalk v1.10.3 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/viper v1.7.1
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v0.0.0-20180714160509-73f8eece6fdc // indirect
	go.mongodb.org/mongo-driver v1.8.0
)
