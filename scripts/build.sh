#!/bin/sh

# 脚本使用说明
# 功能概述：
# 该脚本主要用于编译 Go 项目，支持不同操作系统和架构的编译，并且会根据环境变量或传入参数设置编译版本号。编译后的二进制文件会输出到对应的目录中。
#
# 环境变量：
# - CI_COMMIT_SHA: 用于指定提交的 SHA 值，若未设置，脚本会通过 `git rev-parse HEAD` 命令获取当前 Git 仓库的最新提交 SHA 值。
# - DRONE_TAG: 用于指定版本标签，若设置了该环境变量，编译的版本号会包含此标签信息。
#
# 参数说明：
# - 若传入参数为 "darwin"，脚本会针对 macOS（darwin）系统的 amd64 架构进行编译。
# - 若不传入参数或传入其他值，脚本会默认针对 Linux 系统的 amd64 和 arm64 架构进行编译。
#
# 版本号设置：
# - 若设置了 DRONE_TAG 环境变量或传入了非空的第一个参数作为标签，版本号的格式为 "<标签>-<提交 SHA 值>"。
# - 若未设置 DRONE_TAG 且未传入有效标签，版本号即为提交的 SHA 值。
#
# 使用示例：
# 1. 针对 Linux 系统的 amd64 和 arm64 架构编译，使用默认版本号：
#    ./script.sh
# 2. 针对 macOS（darwin）系统的 amd64 架构编译，使用默认版本号：
#    ./script.sh darwin
# 3. 若设置了 DRONE_TAG 环境变量，例如：
#    DRONE_TAG=v1.0 ./script.sh
#    此时编译的版本号会包含 v1.0 标签信息。
# 4. 若直接第一个参数为TAG，例如：
#    ./script.sh v1.0
#
# 注意事项：
# - 运行此脚本前，请确保 Go 开发环境已正确安装和配置，因为脚本依赖 Go 编译器进行代码编译。
# - 脚本假设项目的入口文件为 cmd/server/main.go，若项目结构不同，请相应修改脚本中的文件路径。
# - 编译后的二进制文件会输出到 dist 目录下对应的操作系统和架构子目录中，确保该目录有足够的权限进行文件写入操作。

sha=${CI_COMMIT_SHA:-${GITHUB_SHA:-$(git rev-parse HEAD 2>/dev/null)}}
tag=${1:-${DRONE_TAG:-$GITHUB_REF_NAME}}
version=$( [ -n "$tag" ] && printf "%s-%s" "$tag" "$sha" || echo $sha )
echo "sha: $sha, tag: $tag, version: $version"

build() {
  local os=$1
  local arch=$2
  local main=$3
  local binary_file=$4

  local binary_dir="dist/${os}/${arch}"
  local binary_path="${binary_dir}/${binary_file}"
  local rm_cmd="rm -f $binary_path"

  local cmd="GOOS=$os GOARCH=$arch go build -a -tags netgo -ldflags=\"-s -w  -extldflags '-static' -X main.version=$version\" -o $binary_path $main"
  echo $rm_cmd
  eval $rm_cmd
  echo $cmd
  eval $cmd
}

build_server() {
  local os=$1
  local arch=$2
  build $os $arch ./cmd/server drone_ext_envs
}

build_client() {
  local os=$1
  local arch=$2
  build $os $arch ./cmd/client drone_ext_envs_client
}

build_client="0"
build_server="0"
if printf "%s\n" "$@" | grep -q -- '--build-client'; then
  build_client="1"
elif printf "%s\n" "$@" | grep -q -- '--build-server'; then
  build_server="1"
else
  build_client="1"
  build_server="1"
fi



if [ "$1" = "darwin" ]; then
  build_server darwin amd64
  exit 0
fi



if [ "$build_server" = "1" ]; then
  build_server linux amd64
  build_server linux arm64
fi

if [ "$build_client" = "1" ]; then
  build_client linux amd64
  build_client linux arm64
fi



