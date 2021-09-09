package plugin

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Internal struct {
	c        interfaces.GrpcClient
	modelSvc interfaces.GrpcClientModelService
	api      *gin.Engine
	apiSvr   *http.Server
	p        *models.Plugin
}

func (internal *Internal) GetGrpcClient() interfaces.GrpcClient {
	return internal.c
}

func (internal *Internal) GetModelService() interfaces.GrpcClientModelService {
	return internal.modelSvc
}

func (internal *Internal) GetApi() *gin.Engine {
	return internal.api
}

func (internal *Internal) GetApiServer() *http.Server {
	return internal.apiSvr
}

func (internal *Internal) StartApi() {
	if err := internal.apiSvr.ListenAndServe(); err != nil {
		log.Info("plugin stopped")
	}
}

func (internal *Internal) StopApi() {
	_ = internal.apiSvr.Shutdown(context.Background())
}

func (internal *Internal) loadConfig() {
	// load plugin.json
	var p models.Plugin
	data, err := ioutil.ReadFile("./plugin.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &p); err != nil {
		panic(err)
	}

	// set plugin
	internal.p = &p

	// set environment
	_ = os.Setenv("CRAWLAB_PLUGIN_NAME", p.Name)
}

func (internal *Internal) register() {
	// plugin model service
	pluginSvc, err := internal.modelSvc.NewBaseServiceDelegate(interfaces.ModelIdPlugin)
	if err != nil {
		panic(err)
	}

	// attempt to get from db
	doc, err := pluginSvc.Get(bson.M{"name": internal.p.Name}, nil)
	if err != nil {
		if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
			// not exists, add to db
			if err := client.NewModelDelegate(internal.p).Add(); err != nil {
				panic(err)
			}

			return
		}
		panic(err)
	}

	// exists, update
	internal.p.SetId(doc.GetId())
	if err := client.NewModelDelegate(internal.p).Save(); err != nil {
		panic(err)
	}
}

func (internal *Internal) _apiGetFile(c *gin.Context) {
	relPath := strings.Replace(c.Request.URL.Path, "/", "", 1)
	filePath, err := filepath.Abs(relPath)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}
	c.File(filePath)
	c.AbortWithStatus(http.StatusOK)
}

func NewInternal() *Internal {
	var err error

	// service
	internal := &Internal{}

	// load config
	internal.loadConfig()

	// grpc client
	internal.c, err = NewGrpcClient()
	if err != nil {
		panic(err)
	}

	// start grpc client
	if err := internal.c.Start(); err != nil {
		panic(err)
	}

	// model service
	internal.modelSvc, err = client.NewServiceDelegate()
	if err != nil {
		panic(err)
	}

	// api
	internal.api = gin.New()
	internal.api.GET("/*path", internal._apiGetFile)
	_ = middlewares.InitMiddlewares(internal.api)

	// register
	internal.register()

	// api server
	internal.apiSvr = &http.Server{
		Addr:    internal.p.Endpoint,
		Handler: internal.api,
	}

	return internal
}
