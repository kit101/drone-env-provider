package reader

import (
	"fmt"
	"github.com/kit101/drone-ext-envs/pkg"
	"os"
)

type FileReader struct {
	Filepath string
}

// 从ConfigMap获取数据
func (r *FileReader) Read() (*pkg.Envs, []byte, error) {
	data, err := os.ReadFile(r.Filepath)
	if err != nil {
		return nil, data, fmt.Errorf("无法读取文件: %w", err)
	}
	return read(data)
}
