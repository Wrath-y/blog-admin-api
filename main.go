package main

import (
	"blog-admin-api/pkg/config"
	"blog-admin-api/pkg/db"
	"blog-admin-api/pkg/def"
	"blog-admin-api/pkg/goredis"
	"blog-admin-api/pkg/httplib"
	"blog-admin-api/pkg/logging"
	"blog-admin-api/router"
	"context"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setup() {
	config.Setup()
	// 监听nacos变化，发现变化后会自动同步到本地，同时杀掉当前进程（之后pod拉起）
	config.ListenNacos(logging.New(), httplib.NewClient(httplib.WithTimeout(30*time.Second), httplib.WithTransport(httplib.NewClient(httplib.WithNoLog(true)))))
	for {
		if config.HasInit {
			break
		}
		println("wait for nacos sync")
		time.Sleep(time.Second)
	}
	logging.Setup(viper.GetString("app.log.topic"), logger)
	httplib.Setup(viper.GetString("app.log.topic"), logger)
	db.Setup()
	goredis.Setup()
}

func main() {
	var err error

	setup()

	switch viper.GetString("app.env") {
	case def.EnvDevelopment:
		gin.SetMode(gin.DebugMode)
	case def.EnvTesting:
		gin.SetMode(gin.TestMode)
	case def.EnvProduction:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.Register()

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	log.Println(color.GreenString("项目启动地址 %s", srv.Addr))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		if err == context.DeadlineExceeded {
			log.Println("Server Shutdown: timeout of 3 seconds.")
		} else {
			log.Println("Server Shutdown Error: ", err)
		}
	}
	log.Println("Server exited")
}

func logger(data []byte) {
	switch viper.GetString("app.log.output") {
	case "stdout":
		logging.StdoutLogger(data)
	case "file":
		logging.FileLogger(data)
	default:
		logging.StdoutLogger(data)
	}
}
