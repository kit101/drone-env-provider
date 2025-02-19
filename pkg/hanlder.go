package pkg

import (
	"bytes"
	"context"
	"github.com/drone/drone-go/plugin/environ"
	"github.com/drone/drone-go/plugin/logger"
	"github.com/kit101/drone-ext-envs/pkg/loggor"
	"strings"
)

type (
	Envs struct {
		Common    []EnvVar            `json:"common" yaml:"common"`
		RepoSlugs map[string][]EnvVar `json:"repoSlugs" json:"repo-slugs" yaml:"repoSlugs" yaml:"repo-slugs"`
	}
	// EnvVar 定义环境变量结构体
	EnvVar struct {
		Name  string                 `json:"name" yaml:"name"`
		Data  string                 `json:"data" yaml:"data"`
		Mask  bool                   `json:"mask" yaml:"mask"`
		Extra map[string]interface{} `json:"extra" yaml:"extra"`
	}

	EnvsReader interface {
		Read() (*Envs, []byte, error)
	}

	EnvsHandler struct {
		reader EnvsReader
		preRaw []byte
		log    logger.Logger
	}
)

func NewEnvHandler(reader EnvsReader, log logger.Logger) *EnvsHandler {
	return &EnvsHandler{
		reader: reader,
		log:    log,
	}
}

func (p *EnvsHandler) List(ctx context.Context, r *environ.Request) ([]*environ.Variable, error) {
	envs, raw, err := p.reader.Read()
	if err != nil {
		loggor.Default.Errorf("raw: \n%s\n", string(raw))
		return nil, err
	}

	if !bytes.Equal(raw, p.preRaw) {
		if p.preRaw == nil {
			p.log.Debugln("init envs: \n", string(raw), "\n")
		} else {
			p.log.Debugln("envs changed: \n", string(raw), "\n")
		}
		p.preRaw = raw
	}

	vars := Convert(envs.Common)

	repoEnvs := envs.RepoSlugs[r.Repo.Slug]
	if repoEnvs != nil {
		vars = append(vars, Convert(repoEnvs)...)
	}

	p.log.Info("return envs, repo:", r.Repo.Slug, ", build:", r.Build.Number, ", vars:", varNamesStr(vars))

	return vars, nil
}

func Convert(src []EnvVar) []*environ.Variable {
	var dest []*environ.Variable
	for _, ev := range src {
		dest = append(dest, &environ.Variable{
			Name: ev.Name,
			Data: ev.Data,
			Mask: ev.Mask,
		})
	}
	return dest
}

func varNames(vars []*environ.Variable) []string {
	var varsNames []string
	for _, v := range vars {
		varsNames = append(varsNames, v.Name)
	}
	return varsNames
}

func varNamesStr(vars []*environ.Variable) string {
	return strings.Join(varNames(vars), ",")
}
