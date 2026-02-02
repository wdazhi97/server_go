# buildtool — 公共构建

本目录提供 **公共运行时基础镜像**，供各服务目录的 Earthfile 使用。

- **Earthfile**：`runtime` target，基于 Ubuntu 24.04，安装 ca-certificates，创建 `/root/bin`。
- 各服务目录（gateway/、lobby/、…）的 Earthfile 通过 `FROM ../buildtool+runtime` 依赖此 base。
- 各服务**不拷贝父目录**：编译在根 Earthfile 的 `xxx-build` target 里完成（因 go.mod 在根目录），各目录用 `BUILD ../+xxx-build` 取结果，只做 docker 阶段。

单独构建某服务时，在**仓库根目录**执行：`earthly ./gateway+docker`。
