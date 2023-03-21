

# gitout

####  介绍

指定两个版本的 `commit-id` 导出变更文件，同时打增量zip包，方便发版。



#### 使用示例

```
gitout -dir=D:/www/XXX -version=XXX_V1.0.5 -currentId=6f202419 -lastId=64b51dad -outDir=./
```



#### 参数说明:

|   参数    |                   说明                    |          示例值          |
| :-------: | :---------------------------------------: | :----------------------: |
|   debug   |             是否开启调试模式              |  bool 值，默认 `false`   |
|    dir    |            项目所在目录全路径             |        D:/www/XXX        |
|  outDir   |              输出文件夹路径               | 默认使用当前程序所在目录 |
|  version  | 项目版本号，用于导出目录命名和zip打包命名 |        XXX_V1.0.5        |
| currentId |          当前版本的 `commit-id`           |     默认使用 `HEAD`      |
|  lastId   |          上个版本的 `commit-id`           |    默认使用 `HEAD~1`     |



#### 注意

由于 Git 仅仅跟踪文件的变动，不跟踪目录。要导出空目录的话需在空目录下创建 `.gitkeep` 文件，然后在项目的 `.gitignore` 中设置不忽略 `.gitkeep` 。



#### 下载

你可以下载源码并自行构建

```
go build -ldflags="-w -s" -o gitout.exe main.go
```

也可以从以下地址下载已构建好的版本：

- [gitout Download Page (GitHub release)](https://github.com/YuuYong/gitout/releases)


####  参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request