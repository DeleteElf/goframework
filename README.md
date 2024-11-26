# goframework
- 作者：DeleteElf
- 联系邮箱：48207475@qq.com

go的微服务与低代码配置框架，旨在于提高复用度
支持领域：
1. 电商平台
2. 游戏引擎
3. im系统
4. 生产平台

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


### 版本变更日志
#### 0.1
1.初始化框架功能逻辑