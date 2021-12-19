package plugin

import (
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Internal struct {
	c        interfaces.GrpcClient
	modelSvc interfaces.GrpcClientModelService
	api      *gin.Engine
	apiSvr   *http.Server
	p        *models.Plugin
	eventSvc EventServiceInterface
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

func (internal *Internal) GetEventService() EventServiceInterface {
	return internal.eventSvc
}

func (internal *Internal) StartApi() {
	if err := internal.apiSvr.ListenAndServe(); err != nil {
		trace.PrintError(err)
	}
	log.Info("plugin stopped")
}

func (internal *Internal) StopApi() {
	if err := internal.apiSvr.Close(); err != nil {
		trace.PrintError(err)
	}
}

func (internal *Internal) Wait() {
	utils.DefaultWait()
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
	relPath := strings.Replace(c.Request.URL.Path, "/_ui/", "", 1)
	filePath, err := filepath.Abs(relPath)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}
	c.File(filePath)
	c.AbortWithStatus(http.StatusOK)
}

func (internal *Internal) _apiGetLang(c *gin.Context) {
	relPath := internal.p.LangUrl
	dirPath, err := filepath.Abs(relPath)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}
	res := bson.M{}
	for _, f := range utils.ListDir(dirPath) {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		filePath := path.Join(dirPath, f.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			controllers.HandleErrorInternalServerError(c, err)
			return
		}
		var langData bson.M
		if err := json.Unmarshal(data, &langData); err != nil {
			controllers.HandleErrorInternalServerError(c, err)
			return
		}
		lang := strings.Replace(f.Name(), ".json", "", 1)
		res[lang] = langData
	}
	controllers.HandleSuccessWithData(c, res)
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
	internal.api.GET("/_ui/*path", internal._apiGetFile)
	internal.api.GET("/_lang", internal._apiGetLang)
	_ = middlewares.InitMiddlewares(internal.api)

	// register
	if viper.GetBool("plugin.register") {
		internal.register()
	}

	// api server
	internal.apiSvr = &http.Server{
		Addr:    internal.p.Endpoint,
		Handler: internal.api,
	}

	// event service
	internal.eventSvc = NewEventService(internal)
	if err := internal.eventSvc.Subscribe(); err != nil {
		panic(err)
	}

	return internal
}
