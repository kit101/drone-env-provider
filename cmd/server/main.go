package main

import (
	"flag"
	"fmt"
	"github.com/kit101/drone-ext-envs/pkg"
	"github.com/kit101/drone-ext-envs/pkg/loggor"
	"github.com/kit101/drone-ext-envs/pkg/reader"
	"net/http"
	"os"
)

// 全局变量
var (
	logLevel string

	listenPort string
	secretKey  string

	from      string
	configMap string
	filePath  string

	version = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("version: " + version)
		return
	}

	// 定义命令行参数
	flag.StringVar(&logLevel, "log-level", os.Getenv("EXT_ENV_LOG_LEVEL"), "日志级别 (debug/info/warn/error)")

	flag.StringVar(&listenPort, "port", os.Getenv("PORT"), "监听端口")
	flag.StringVar(&secretKey, "secret-key", os.Getenv("EXT_ENV_SECRET_KEY"), "访问密钥")

	flag.StringVar(&from, "from", os.Getenv("EXT_ENV_FROM"), "数据源 (file/k8s-cm)")
	flag.StringVar(&configMap, "configmap", os.Getenv("EXT_ENV_CONFIGMAP"), "Kubernetes NS/ConfigMap名称, e.g. your_namespace/your_configmap")
	flag.StringVar(&filePath, "file", os.Getenv("EXT_ENV_FILE"), "本地配置文件路径")

	flag.Parse()

	loggor.Default = loggor.New(logLevel)

	// 检查必要参数
	if from == "" {
		from = "file"
	}
	if secretKey == "" {
		loggor.Default.Warnf("请指定访问密钥 (--secret-key 或 EXT_ENV_SECRET_KEY 环境变量)")
	}
	if listenPort == "" {
		listenPort = "8080"
	}
	if from == "file" && filePath == "" {
		loggor.Default.Errorln("当数据源为文件时，必须指定配置文件路径 (--file 或 EXT_ENV_FILE 环境变量)")
		return
	}
	if from == "k8s-cm" && (configMap == "") {
		loggor.Default.Errorln("当数据源为Kubernetes ConfigMap时，必须指定ConfigMap, e.g. your_ns/your_cm (--configmap 或 EXT_ENV_CONFIGMAP 环境变量)")
		return
	}

	var r pkg.EnvsReader
	if from == "file" {
		r = reader.NewFileReader(filePath)
	} else if from == "k8s-cm" {
		r = &reader.K8sCMReader{Configmap: configMap}
	} else {
		loggor.Default.Errorln("不合法的--from: ", from)
		return
	}
	p := pkg.NewEnvPlugin(r, loggor.Default)
	handler := pkg.Handler(secretKey, p, loggor.Default)
	http.HandleFunc("/envs", handler.ServeHTTP)

	http.HandleFunc("/healthy", healthy)

	// 启动服务器
	loggor.Default.Info("服务器启动，监听端口: ", listenPort)
	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		loggor.Default.Errorln("服务器启动失败: %v", err)
	}
}

func healthy(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
