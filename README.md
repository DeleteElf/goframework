# goframework go的微服务与低代码配置框架，旨在于提高复用度
- 作者：DeleteElf
- 联系邮箱：48207475@qq.com

访问外网如果有困难，可以使用翻墙梯子 XFLDT(主打比较便宜稳定)，并使用Clash进行访问
[链接1](https://dash.xfltd.app/register?code=KXZSOGgs "点击跳转")
[链接2](https://xfltd.net/#/register?code=KXZSOGgs)
[国内可访问免代理](https://xfltd.top/#/register?code=KXZSOGgs)
#### 支持领域：
1. 电商平台
2. 游戏引擎
3. im系统
4. 生产平台
5. 区块链系统
#### 相关项目示例
示例暂未对外开放，如有需要，请联系作者
##### 非GRPC项目
[单服务示例](https://github.com/DeleteElf/GoWebSite "点击跳转")

##### GRPC项目
[GRPC服务端示例](https://github.com/DeleteElf/goetcdserver "点击跳转")
[GRPC客户端示例](https://github.com/DeleteElf/goetcdclient "点击跳转")

###### 方法一、通过shell命令行加载本模块
```shell
# github上public的项目，无需设置这个，
# 如果要访问私有项目，则可以如下方式设置这个账户下的不走代理，不走代理需要使用梯子
go env -w GOPRIVATE=github.com/deleteelf
#拉取模块 
go get github.com/deleteelf/goframework
#清理一下代码区域，防止刚拉的代码，被就代码搞混乱了
go mod tidy
```
###### 方法二、通过多模块配置加载此模块
多模块配置无需使用go get github.com/deleteelf/goframework来获取最新代码
```shell
#初始化多模块工作区，模块是子目录，goland打开时使用父目录打开
go work init ./goframework
#增加第2个工作区，这边只是示例
go work use ./example
```


1. loghelper 日志工具类，如果没有初始化级别，默认初始化warn级别日志
2. stringhelper 字符串帮助类
3. httphelper http工具栏
4. ado 基于gorm的orm体系数据库操作类，同时实现非orm体系支持，以减少model生成需要，减少代码量，最终达到低代码的需求
5. entities 实体类，一些常用的对象关系基类，可用于扩展属性
