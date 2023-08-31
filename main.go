package main

import (
	"blog-admin-api/pkg/config"
	"blog-admin-api/pkg/db"
	"blog-admin-api/pkg/goredis"
	"blog-admin-api/pkg/logging"
	"blog-admin-api/router"
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// setupConfigYaml 就绪配置文件
// 环境变量配置 NACOS_SKIP="Y", 可跳过下载配置
// 环境变量:
// NACOS_USE=false
// NACOS_NAMESPACE=""
// NACOS_SERVER=""
// NACOS_USERNAME=""
// NACOS_PASSWORD=""
func setupConfigYaml() {
	viper.AutomaticEnv()
	if envUse := viper.GetBool("NACOS_USE"); !envUse {
		log.Println("跳过从nacos下载配置文件")
		return
	}

	config.SetupNacosClient()
	config.DownloadNacosConfig()

	// 未使用k8s部署时监听nacos（已经被使用的变量不会体现出变化）
	//config.ListenNacos()

	// 使用k8s部署时监听nacos（已经被使用的变量不会体现出变化）
	config.ListenNacos(func(cnf string) {
		log.Println("当前进程将被停止")
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	})
}

func setup() {
	config.Setup()
	logging.Setup()
	goredis.Setup()
	db.Setup()
}

func main() {
	// rand.Seed(time.Now().UnixNano()) // version < 1.20
	flag.Parse()
	setupConfigYaml()

	setup()
	gin.SetMode(gin.ReleaseMode)
	r := router.Register()
	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	logging.New().Info("项目启动成功", srv.Addr, viper.GetString("app.env"))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.New().Info("Shutdown Server ...", srv.Addr, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		if err == context.DeadlineExceeded {
			logging.New().Info("Server Shutdown: timeout of 3 seconds.", srv.Addr, nil)
		} else {
			logging.New().ErrorL("Server Shutdown Error:", srv.Addr, err)
		}
	}
	logging.New().Info("Server exited", srv.Addr, nil)
}
