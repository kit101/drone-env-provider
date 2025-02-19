#!/bin/sh

# 脚本使用说明
# 功能概述：
# 此脚本利用 skopeo 工具，将指定 Docker 镜像仓库中的特定镜像从一个标签复制到其他多个标签。
# skopeo 是用于在不同容器镜像存储库之间复制、同步和检查镜像的工具。
# install doc: https://github.com/containers/skopeo/blob/main/install.md
#
# 环境变量：
# - REGISTRY: Docker 镜像仓库地址，若未设置，默认使用 docker.io。
# - IMAGE: Docker 镜像名称，若未设置，默认使用 kit101z/drone-ext-envs。
#
# 参数说明：
# 脚本运行时需传入至少两个参数：
# - 第一个参数 ($1): 作为源的镜像标签，脚本会把该标签对应的镜像复制到其他指定标签。
# - 后续参数: 目标镜像标签，脚本会将源标签的镜像复制到这些目标标签。
#
# 使用示例：
# 假设你要将镜像 kit101z/drone-ext-envs 的 v1.0 标签复制到 v1.1 和 v1.2 标签，可执行以下命令：
# 若使用默认的镜像仓库和镜像名称
#   ./script.sh v1.0 v1.1 v1.2
# 若要指定自定义的镜像仓库和镜像名称
#   REGISTRY=myregistry.example.com IMAGE=myimage ./script.sh v1.0 v1.1 v1.2
#
# 注意事项：
# - 运行此脚本前，请确保 skopeo 工具已正确安装，因为脚本依赖该工具进行镜像复制操作。
# - 运行前先使用skopeo login registry -u username -p password 登录镜像仓库。
# - 确保你对指定的 Docker 镜像仓库有读写权限，否则镜像复制操作可能会失败。
# - 脚本在执行过程中，会先打印出要执行的 skopeo 命令，然后执行该命令。若执行过程中出现错误，需根据错误信息进行排查，可能是网络问题、权限问题或 skopeo 工具本身的问题。


registry=${REGISTRY:-docker.io}
image=${IMAGE:-kit101z/drone-ext-envs}

src="$1"

for arg in "$@"; do
    if [ "$arg" != "$src" ]; then
      cmd="skopeo copy docker://$registry/$image:$src docker://$registry/$image:$arg --all"
      echo $cmd
      exec $cmd
    fi
done