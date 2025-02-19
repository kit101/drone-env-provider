package reader

import (
	"context"
	"fmt"
	"github.com/kit101/drone-ext-envs/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strings"
)

type K8sCMReader struct {
	Configmap string
}

// 从ConfigMap获取数据
func (r *K8sCMReader) Read() (*pkg.Envs, []byte, error) {
	// 构建kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		// 如果在集群外运行，尝试使用本地 kubeconfig
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("无法构建 Kubernetes 配置: %v", err)
		}
	}

	// 创建客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("无法创建Kubernetes客户端: %w", err)
	}

	// 获取ConfigMap
	cmSplits := strings.Split(r.Configmap, "/")
	configMapNS := cmSplits[0]
	configMapName := cmSplits[1]
	cm, err := clientset.CoreV1().ConfigMaps(configMapNS).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("无法获取ConfigMap: %w", err)
	}

	// 假设ConfigMap只有一个键值对，且值为YAML或JSON格式
	for _, data := range cm.Data {
		raw := []byte(data)
		return read(raw)
	}

	return nil, nil, fmt.Errorf("ConfigMap中没有可用数据")
}
