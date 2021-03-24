package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/qibobo/webgo-gin/config"
	"github.com/qibobo/webgo-gin/controller"
	"github.com/qibobo/webgo-gin/db"
	"github.com/qibobo/webgo-gin/db/sqldb"
	_ "github.com/qibobo/webgo-gin/docs"
	"github.com/qibobo/webgo-gin/logging"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {

	var path string
	flag.StringVar(&path, "c", "", "config file")
	flag.Parse()
	if path == "" {
		fmt.Fprintln(os.Stderr, "missing config file")
		os.Exit(1)
	}

	configFile, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to open config file '%s' : %s\n", path, err.Error())
		os.Exit(1)
	}
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to read config file '%s' : %s\n", path, err.Error())
		os.Exit(1)
	}
	var conf *config.Config
	conf, err = config.LoadConfig(configBytes)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to read config: '%s'\n", err.Error())
		os.Exit(1)
	}
	configFile.Close()

	err = conf.Validate()
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to validate configuration : %s\n", err.Error())
		os.Exit(1)
	}

	logger, err := logging.NewLogger("webgo")
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to init logger : %s\n", err.Error())
		os.Exit(1)
	}
	demodb, err := sqldb.NewDemoSQLDB(conf.DB.DemoDB, *logger.Named("demodb"))
	if err != nil {
		logger.Error("failed to connect to demodb", zap.Error(err))
		os.Exit(1)
	}

	server := createServer(logger, demodb)
	server.Run(fmt.Sprintf(":%d", conf.Server.Port))
}

func createServer(logger *zap.Logger, demoDB db.DemoDB) *gin.Engine {
	r := gin.Default()

	c := controller.NewDemoController(logger, demoDB)

	v1 := r.Group("/api/v1")
	{
		demo := v1.Group("/demo")
		{
			demo.GET("/:id", c.GetById)
			demo.POST("", c.CreateDemo)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
