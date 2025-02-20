package reader

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/kit101/drone-ext-envs/pkg"
	"github.com/kit101/drone-ext-envs/pkg/loggor"
	"os"
	"path/filepath"
	time2 "time"
)

const interval = 5 * time2.Second

type fileReader struct {
	Filepath string
	sha      string
	mtime    int64
	envs     *pkg.Envs
	raw      []byte
	err      error
}

func NewFileReader(fp string) *fileReader {
	abs, _ := filepath.Abs(fp)
	r := &fileReader{
		Filepath: abs,
	}
	r.watch()
	return r
}

func (r *fileReader) Read() (*pkg.Envs, []byte, error) {
	return r.envs, r.raw, r.err
}

func (r *fileReader) watch() {
	go func() {
		for {
			r.doRead()
			time2.Sleep(interval)
		}
	}()
}

func (r *fileReader) doRead() {
	time, err := fileutil.MTime(r.Filepath)
	if err != nil {
		loggor.Default.Warnf("无法获取文件修改时间: %v", err)
	}
	sha, err := fileutil.Sha(r.Filepath, 256)
	if sha != r.sha {
		// 记录变更后的sha值
		r.mtime = time
		loggor.Default.Infof("file [%s] changed at '%s', sha: '%s' -> '%s'",
			r.Filepath, time2.Unix(r.mtime, 0), r.sha, sha)
		r.sha = sha

		// 读取并保存数据
		raw, err := os.ReadFile(r.Filepath)
		if err != nil {
			r.envs = nil
			r.raw = raw
			r.err = err
		} else {
			r.envs, r.raw, r.err = read(raw)
		}
	} else {
		loggor.Default.Debugf("file [%s] not changed", r.Filepath)
	}
}
