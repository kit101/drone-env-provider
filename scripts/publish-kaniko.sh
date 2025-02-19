#!/bin/sh

# 以下是本脚本的使用说明：

# 脚本功能概述
# 该脚本主要用于构建和推送 Docker 镜像，支持多架构（`amd64` 和 `arm64`）镜像的构建和合并

# 环境变量
# - `REGISTRY`：Docker 镜像仓库地址，默认为 `docker.io`。
# - `IMAGE`：Docker 镜像名称，默认为 `kit101z/drone-ext-envs`。

# 参数说明
# 脚本需要传入三个位置参数：
# 1. `registry_username`：Docker 镜像仓库的用户名。
# 2. `registry_password`：Docker 镜像仓库的密码。
# 3. `tag`：Docker 镜像的标签。

# 使用步骤
# 1. 设置环境变量（可选）
# 如果需要使用自定义的镜像仓库地址或镜像名称，可以设置相应的环境变量。例如：
# ```sh
# export REGISTRY=myregistry.example.com
# export IMAGE=myimage
# ```
# 2. 运行脚本
# 在终端中运行脚本，并传入所需的参数。以下是不同模式的示例：
# ```sh
# ./script.sh myusername mypassword v1.0.0
# ```
# 该命令会先构建 `amd64` 和 `arm64` 架构的镜像，然后将它们合并为一个多架构镜像。

# 注意事项
# - 脚本依赖 `docker`，并会拉取`kaniko executor` 和 `manifest-tool`镜像，请确保这些工具已经正确安装并配置。
# - 请确保当前目录下存在 `Dockerfile` 和 `scripts/publish-kaniko.sh` 文件。

registry=${REGISTRY:-docker.io}
image=${IMAGE:-kit101z/drone-ext-envs}

registry_username="$1"
registry_password="$2"
tag="$3"

# 定义生成认证信息并写入文件的函数
generate_docker_auth() {
    local registry="$1"
    local username="$2"
    local password="$3"
    local output_file="$4"

    if [ $registry == "docker.io" ]; then
      registry="https://index.docker.io/v1/"
    else
      registry="https://${registry}/"
    fi

    # 将用户名和密码用冒号拼接
    credentials="${username}:${password}"
    # 对拼接后的字符串进行 Base64 编码
    encoded=$(echo -n "$credentials" | base64)

    # 创建目录
    mkdir -p $(dirname $output_file)

    # 写入 JSON 内容到文件
    {
        echo   "{"
        echo   "  \"auths\": {"
        printf '    "%s": {\n' $registry
        printf "      \"auth\": \"%s\"\n" "$encoded"
        echo   "    }"
        echo   "  }"
        echo   "}"
    } > "$output_file"

    # 提示用户操作完成
    echo "认证信息已成功写入 $output_file 文件"
    # cat $output_file
}

exec_executor() {
  docker run --rm -it -v .:/workspace -w /workspace --entrypoint sh gcr.io/kaniko-project/executor:debug $0 $@ --executor
}

exec_manifest() {
  docker run --rm -it -v .:/workspace -w /workspace --entrypoint sh mplatform/manifest-tool:alpine $0 $@ --manifest
}

remove_specific_arg() {
    local new_args=""
    local arg_to_remove="$1"
    shift
    for arg in "$@"; do
        if [ "$arg" != "$arg_to_remove" ]; then
            new_args="$new_args $arg"
        fi
    done
    echo "$new_args"
}

if printf "%s\n" "$@" | grep -q -- '--executor'; then
  generate_docker_auth $registry $registry_username $registry_password /kaniko/.docker/config.json
  arch=amd64
  executor --custom-platform linux/$arch --context . -f Dockerfile -d $registry/$image:$tag-$arch --build-arg TARGETOS=linux --build-arg TARGETARCH=$arch
  arch=arm64
  executor --custom-platform linux/$arch --context . -f Dockerfile -d $registry/$image:$tag-$arch --build-arg TARGETOS=linux --build-arg TARGETARCH=$arch
elif printf "%s\n" "$@" | grep -q -- '--manifest'; then
  generate_docker_auth $registry $registry_username $registry_password /root/.docker/config.json
  manifest-tool --debug push from-args --platforms=linux/amd64,linux/arm64 --template $registry/$image:$tag-ARCH --target $registry/$image:$tag
elif printf "%s\n" "$@" | grep -q -- '--only-executor'; then
  new_args=$(remove_specific_arg "--only-executor" "$@")
  set -- $new_args
  exec_executor $new_args
elif printf "%s\n" "$@" | grep -q -- '--only-manifest'; then
  new_args=$(remove_specific_arg "--only-manifest" "$@")
  exec_manifest $new_args
else
  exec_executor $@
  exec_manifest $@
fi