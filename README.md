# ChatBot - 自动聊天机器人

一个用Go语言实现的简单自动聊天机器人项目，支持中文交互。

## 项目结构

```
chatbot/
├── cmd/
│   └── chatbot/          # 主应用程序
│       └── main.go       # 应用入口点
├── internal/
│   ├── bot/             # 机器人核心逻辑
│   │   ├── bot.go       # 机器人实现
│   │   └── bot_test.go  # 测试文件
│   ├── config/          # 配置管理
│   │   └── config.go    # 配置实现
│   └── handler/         # 输入处理
│       └── input_handler.go  # 输入处理器
├── chatbot.exe          # 编译后的可执行文件
├── go.mod               # Go模块文件
└── README.md            # 项目文档
```

## 功能特性

- 基本的聊天功能
- 关键词匹配响应
- 支持中英文输入
- 简单的配置管理
- 命令行交互界面
- 可扩展的架构

## 支持的命令

- `hello`/`hi`/`你好`/`嗨` - 问候
- `how are you`/`你好吗`/`怎么样` - 询问机器人状态
- `name`/`名字` - 询问机器人名字
- `help`/`帮助` - 获取帮助信息
- `time`/`时间` - 询问时间（仅返回预设响应）
- `weather`/`天气` - 询问天气（仅返回预设响应）
- `exit`/`退出` - 退出程序

## 如何运行

1. **构建项目**
   ```bash
   go build -o chatbot.exe ./cmd/chatbot
   ```

2. **运行程序**
   ```bash
   ./chatbot.exe
   ```

3. **与机器人聊天**
   ```
   聊天机器人已启动。输入'exit'或'退出'退出程序。
   ----------------------------------------
   你: 你好
   Bot: 你好！今天我能帮你什么忙？
   你: 你叫什么名字？
   Bot: 我的名字是ChatBot。很高兴认识你！
   你: 退出
   再见！
   ```

## 运行测试

```bash
go test ./internal/bot
```

## Web 页面（聊天对话界面）

项目新增了一个 Web 入口（内置页面资源，无需额外安装前端依赖），用于在浏览器里进行聊天对话。

1. **直接运行（开发调试）**

```bash
go run ./cmd/web
```

2. **构建并运行**

```bash
go build -o chatbot-web.exe ./cmd/web
./chatbot-web.exe
```

3. **打开页面**

在浏览器访问 `http://localhost:8080/`。

> 端口可通过环境变量 `CHATBOT_ADDR` 修改，例如 `CHATBOT_ADDR=:9000`。

## 扩展指南

1. **添加新的响应规则**
   - 在 `internal/bot/bot.go` 的 `generateResponse` 函数中添加新的 case 语句

2. **修改配置**
   - 在 `internal/config/config.go` 中修改默认配置值
   - 可以扩展为从环境变量或配置文件加载配置

3. **添加新功能**
   - 可以添加自然语言处理能力
   - 可以集成外部API获取实时信息
   - 可以添加对话历史记录功能

## 技术栈

- Go 语言
- 标准库

## 许可证

本项目使用 MIT 许可证。
