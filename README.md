# goframework
- 作者：DeleteElf
- 联系邮箱：48207475@qq.com

访问外网如果有困难，可以使用翻墙梯子 XFLDT(主打比较便宜稳定)，并使用Clash进行访问
[链接1](https://dash.xfltd.app/register?code=KXZSOGgs "点击跳转")
[链接2](https://xfltd.net/#/register?code=KXZSOGgs)

go的微服务与低代码配置框架，旨在于提高复用度
#### 支持领域：
1. 电商平台
2. 游戏引擎
3. im系统
4. 生产平台

#### 相关项目示例
示例暂未对外开放，如有需要，请联系作者
##### 非GRPC项目
[单服务示例](https://github.com/DeleteElf/GoWebSite "点击跳转")

##### GRPC项目
[GRPC服务端示例](https://github.com/DeleteElf/goetcdserver "点击跳转")
[GRPC客户端示例](https://github.com/DeleteElf/goetcdclient "点击跳转")

通过shell命令行加载本模块
```shell
#设置这个账户下的不走代理
go env -w GOPRIVATE=github.com/deleteelf
#拉取模块 
go get github.com/deleteelf/goframework
#清理一下代码区域，防止刚拉的代码，被就代码搞混乱了
go mod tidy
```


```go
loghelper 日志工具类，如果没有初始化级别，默认初始化warn级别日志
stringhelper 字符串帮助类
httphelper http工具栏
ado 基于gorm的orm体系数据库操作类
entities 实体类，一些常用的对象关系基类，可用于扩展属性

```


#### 版本变更日志
##### 0.1 初始版本
1. 初始化框架功能逻辑
2. 增加数据库到json免orm的体系支持，建立增删改查机制，并支持事务
