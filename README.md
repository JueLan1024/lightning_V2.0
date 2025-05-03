# Lightning Web Application

## 项目简介
Lightning 是一个基于 Go (Golang) 的社区管理 Web 应用程序，支持用户管理、社区管理、帖子发布和投票功能。项目采用CLD架构和gin框架，使用 MySQL、Redis、Kafka、Canal和Docker等组件，支持高并发的社区互动体验。

---

## 功能特性
- **用户管理**：支持用户注册、登录和权限管理。
- **社区管理**：创建、查看和管理社区。
- **帖子功能**：支持帖子发布、编辑和删除。
- **投票功能**：用户可以对帖子进行投票（赞成或反对）。
- **实时消息**：通过 Kafka 实现实时消息处理。

---

## 技术栈
- **后端**：Go (Golang)
- **框架**：gin
- **数据库**：MySQL
- **缓存**：Redis
- **DTS系统**：Canal
- **消息队列**：Kafka
- **前端**：Swagger 文档用于 API 测试
- **容器引擎**：Docker

---

## 详细介绍
- **logger**：采用zap日志库实现快速结构化日志记录和printf风格的日志记录
- **配置获取**：采用Viper获取配置信息
- **参数检验**：采用Validator进行参数检验
- **登录认证**: 采用JWT鉴权以AccessToken和RefreshToken认证的方式进行登录认证
- **用户黑名单**: redis中存储RefreshToken控制用户登录资格
- **数据库**: 采用sqlx执行数据库操作
- **缓存**: 采用redis的String、Hash、Set、ZSet数据格式存储数据
- **缓存与数据库一致性**: 采用数据库binlog->canal->kafka->redis的方式保证一致性
- **避免缓存击穿**: 采用SingleFlight处理同名Key，避免多条请求打到数据库
- **避免缓存穿透**: 采用bloom过滤器，在项目启动时和数据库更新时添加数据ID到过滤器中
- **避免缓存雪崩**: 将社区信息、排名、帖子排名永久存储在redis中,高热度帖子过期时间刷新等方法避免
- **优化查询速度**: 设置Mysql索引，将数据缓存到redis，优先查找缓存
- **消息队列**: 采用Goroutine异步读取发送到Kafka中的消息
- **游标查询**：帖子列表使用游标分页查询
- **顺序查询**：根据帖子热度或发帖时间查询
- **算法评价系统**：实现了随时间权重下降的算法评论系统
- **投票数据持久化**：采用更新redis->发送消息到kafka->读取消息存储到mysql的异步存储方式
- **限流策略**：采用令牌桶进行限流
- **优雅关机**：使用channel接收系统信号延时关闭
- **接口文档**：使用Swagger注释生成接口文档
- **项目发布**：使用docker-compose创建并关联多个组件的容器运行项目

---

## 项目结构
lightning_v2.0/
├── canal/                              # Canal 配置文件目录
│   ├── conf/                           # Canal 的配置文件
│   │   ├── canal.properties            # Canal 主配置文件
│   │   ├── example/                    # Canal 实例配置
│   │   │   ├── instance.properties     # Canal 实例的具体配置
├── mysql/                  # MySQL 相关文件
│   ├── conf/               # MySQL 配置文件
│   │   ├── my.cnf          # MySQL 配置文件
│   ├── init/               # MySQL 初始化脚本
│   │   ├── init.sql        # 数据库初始化 SQL 文件
├── web_app/                            # Web 应用程序代码
│   ├── conf/                           # 配置文件目录
│   |   ├── config.yaml/                # 配置文件
│   ├── controller/                     # 控制器层，提供功能接口
│   │   ├── code.go                     # 定义返回响应代码
│   │   ├── community.go                # 社区管理功能
│   │   ├── doc_response_models.go      # Swagger 返回响应模型
│   │   ├── post.go                     # 帖子管理功能
│   │   ├── request.go                  # 获取*gin.Context信息
│   │   ├── response.go                 # 返回响应方法和模型
│   │   ├── user.go                     # 用户管理功能
│   │   ├── validator.go                # validator自定义
│   │   ├── vote.go                     # 投票功能
│   ├── dao/                            # 数据访问层，封装数据库和缓存操作
│   │   ├── mysql/                      # MySQL 相关操作
|   |   |   ├── community.go            # 社区表管理 
|   |   |   ├── error_code.go           # 错误代码定义
|   |   |   ├── mysql.go                # mysql初始化
|   |   |   ├── post.go                 # 帖子表管理
|   |   |   ├── user.go                 # 用户表管理 
|   |   |   ├── vote.go                 # 投票表管理   
│   │   ├── redis/                      # Redis 相关操作
|   |   |   ├── community.go            # 社区数据管理
|   |   |   ├── error_code.go           # 错误代码定义
|   |   |   ├── keys.go                 # key定义和获取方法
|   |   |   ├── post.go                 # 帖子数据管理
|   |   |   ├── redis.go                # redis初始化
|   |   |   ├── user.go                 # 用户数据管理
|   |   |   ├── vote.go                 # 投票数据管理
│   ├── docs/                           # Swagger 文档目录
│   ├── kafka/                          # Kafka 消息处理逻辑
│   │   ├── community.go                # 社区消息管理
│   │   ├── consumer.go                 # 消费者创建和消息读取方法
│   │   ├── error_code.go               # 错误代码定义
│   │   ├── kafka.go                    # kafka初始化
│   │   ├── post.go                     # 帖子消息管理
│   │   ├── producer.go                 # 生产者创建和消息发送方法
│   │   ├── vote.go                     # 投票消息管理
│   ├── logger/                         # zap日志工具
│   ├── logic/                          # 业务逻辑层
│   │   ├── community.go                # 社区相关逻辑
│   │   ├── cookie.go                   # refreshToken认证逻辑
│   │   ├── error_code.go               # 错误代码定义
│   │   ├── post.go                     # 帖子相关逻辑
│   │   ├── user.go                     # 用户相关逻辑
│   │   ├── vote.go                     # 投票相关逻辑
│   ├── middlewares/                    # 中间件
│   │   ├── auth.go                     # JWT认证中间件
│   │   ├── rateLimit.go                # 限流中间件
│   ├── models/                         # 数据库模型和 SQL 文件
│   │   ├── community.go                # 社区模型
│   │   ├── create_table.sql            # 创建表SQL
│   │   ├── message.go                  # 消息模型
│   │   ├── pagination.go               # 游标分页模型
│   │   ├── params.go                   # request参数模型
│   │   ├── post.go                     # 帖子模型
│   │   ├── user.go                     # 用户模型
│   │   ├── vote.go                     # 投票模型
│   ├── pkg/                            # 公共库
│   │   ├── bloom/                      # 布隆过滤器
│   │   ├── jwt/                        # jwt工具
│   │   ├── snowflake/                  # 雪花ID生成器
│   ├── routes/                         # 路由层，定义 API 路由
│   │   ├── routes.go                   # 路由注册文件
│   ├── settings/                       # 配置初始化
│   ├── tool/                           # 工具类函数
│   ├── main.go                         # 应用程序的入口文件
│   ├── dockerfile                      # Dockerfile 文件
│   ├── go.mod                          # Go 模块依赖管理文件
│   ├── go.sum                          # Go 模块依赖校验文件
│   ├── web_app.log                     # 应用日志文件
├── docker-compose.yaml     # Docker Compose 配置文件
├── README.md               # 项目说明文件

---

## 快速启动
1.在根目录执行 docker-compose up -d
2.在mysql容器中执行./mysql/init/init.sql 中的所有sql语句
3.启动lightning_app容器。如果有报错是因为Kafka的topic和group_id在初始化，重启lightning_app容器即可
4.程序在本机的8081端口运行，访问http://127.0.0.1:8081/swagger/index.html 查看接口文档