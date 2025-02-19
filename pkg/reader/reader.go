package reader

import (
	"encoding/json"
	"fmt"
	"github.com/kit101/drone-ext-envs/pkg"
	"sigs.k8s.io/yaml"
)

func read(raw []byte) (*pkg.Envs, []byte, error) {
	var envs pkg.Envs
	// 尝试解析为JSON
	err := json.Unmarshal(raw, &envs)
	if err == nil {
		return &envs, raw, nil
	}
	// 尝试解析为YAML
	err = yaml.Unmarshal(raw, &envs)
	if err == nil {
		return &envs, raw, nil
	}
	return &envs, raw, fmt.Errorf("无法解析ConfigMap数据为JSON或YAML: %w", err)
}
