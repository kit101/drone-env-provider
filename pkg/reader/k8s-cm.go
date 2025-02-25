package reader

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kit101/drone-ext-envs/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sCMReader struct {
	namespace     string
	configmapName string
	cs            kubernetes.Clientset
}

func K8sCMReader(configmap string) (pkg.EnvsReader, error) {
	cmSplits := strings.Split(configmap, "/")
	if len(cmSplits) != 2 {
		return nil, fmt.Errorf("illegal configmap: %s, should be {namespace}/{configmap_name}", configmap)
	}
	namespace := cmSplits[0]
	configMapName := cmSplits[1]

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
			return nil, fmt.Errorf("not create kubernetes client: %v", err)
		}
	}

	// 创建客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("not create kubernetes client: %v", err)
	}
	r := &k8sCMReader{
		namespace:     namespace,
		configmapName: configMapName,
		cs:            *clientset,
	}
	r.watch()
	return r, nil
}

// 从ConfigMap获取数据
func (r *k8sCMReader) Read() (*pkg.Envs, []byte, error) {
	return r.doRead()
}

func (r *k8sCMReader) watch() {
	// TODO not implement yet.
}

func (r *k8sCMReader) doRead() (*pkg.Envs, []byte, error) {
	cm, err := r.cs.CoreV1().ConfigMaps(r.namespace).Get(context.TODO(), r.configmapName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("无法获取ConfigMap: %w", err)
	}

	// 假设ConfigMap只有一个键值对，且值为YAML或JSON格式
	for _, data := range cm.Data {
		raw := []byte(data)
		return parse(raw)
	}

	return nil, nil, fmt.Errorf("ConfigMap中没有可用数据")
}
