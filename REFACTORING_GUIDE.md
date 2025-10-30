# Huma重构指导文档

## 概述

本文档提供了将剩余服务从fiber迁移到huma的完整指导。所有DTO文件已创建完成，需要重构Handler和Router文件。

## 已完成的服务模板

参考以下已完成的服务作为模板：
- `/workspace/internal/handler/tag.go` - Tag服务Handler示例
- `/workspace/internal/router/tag.go` - Tag服务Router示例
- `/workspace/internal/handler/category.go` - Category服务Handler示例
- `/workspace/internal/router/category.go` - Category服务Router示例

## 快速重构步骤

### 第1步：重构Handler

对于每个服务（以Article为例）：

```bash
# 当前文件：internal/handler/article.go
# 需要修改的内容：

1. 修改import：
   - 移除：github.com/gofiber/fiber/v2
   - 添加：context
   - 添加：github.com/hcd233/aris-blog-api/internal/protocol/dto

2. 修改接口定义：
   将：HandleXxx(c *fiber.Ctx) error
   改为：HandleXxx(ctx context.Context, req *dto.XxxRequest) (*protocol.HumaHTTPResponse[*dto.XxxResponse], error)

3. 修改方法实现：
   - 移除fiber.Ctx相关的Locals调用
   - 改用ctx.Value(constant.CtxKeyUserID)获取userID
   - 改用req.Body、req.Path等访问请求参数
   - 转换service返回的protocol类型为dto类型
   - 使用util.WrapHTTPResponse包装返回值
```

### 第2步：重构Router

```bash
# 当前文件：internal/router/xxx.go
# 需要修改的内容：

1. 修改函数签名：
   将：func initXxxRouter(r fiber.Router)
   改为：func initXxxRouter(xxxGroup *huma.Group)

2. 移除fiber中间件，改用huma：
   将：middleware.JwtMiddleware()
   改为：middleware.JwtMiddlewareForHuma()

3. 使用huma.Register注册路由：
   huma.Register(xxxGroup, huma.Operation{
       OperationID: "operationName",
       Method:      http.MethodGet,
       Path:        "/path",
       Summary:     "Summary",
       Description: "Description",
       Tags:        []string{"tagName"},
       Security: []map[string][]string{
           {"jwtAuth": {}},
       },
   }, handler.HandleXxx)
```

### 第3步：更新router.go

在`/workspace/internal/router/router.go`中：

```go
// 1. 从v1Router中移除旧的initXxxRouter调用
// 2. 在v1Group下添加新的huma路由组：

xxxGroup := huma.NewGroup(v1Group, "/xxx")
initXxxRouter(xxxGroup)
```

## 特殊处理说明

### AI服务的SSE流式响应

AI服务的流式响应需要特殊处理，暂时保持使用fiber的方式，或者需要实现huma的SSE支持。

建议方案：
1. 保留AI服务的fiber路由用于SSE端点
2. 其他非流式AI端点使用huma重构

### Asset服务的文件上传

文件上传需要特殊处理：
1. UploadImage端点可能需要保留fiber实现
2. 或者使用huma的multipart form支持

### 中间件适配

确保所有中间件都有huma版本：
- JWT中间件：`middleware.JwtMiddlewareForHuma()` - 已完成
- 权限检查中间件：需要创建huma版本
- 限流中间件：需要创建huma版本

## 快速重构命令

对于每个服务，执行以下步骤：

```bash
# 1. 备份原文件
cp internal/handler/xxx.go internal/handler/xxx.go.bak
cp internal/router/xxx.go internal/router/xxx.go.bak

# 2. 参考模板文件重构
# 使用tag.go和category.go作为参考

# 3. 更新router.go

# 4. 测试编译
go build ./...
```

## DTO到Protocol类型映射

由于service层仍然使用protocol类型，Handler需要进行类型转换：

### 请求转换（DTO -> Protocol）
```go
svcReq := &protocol.XxxRequest{
    UserID: userID,  // 从ctx获取
    Field1: req.Body.Field1,  // 从DTO获取
    Field2: req.PathParam,    // 从DTO获取
}
```

### 响应转换（Protocol -> DTO）
```go
rsp := &dto.XxxResponse{
    Data: &dto.DataType{
        Field1: svcRsp.Data.Field1,
        Field2: svcRsp.Data.Field2,
    },
}
```

## 注意事项

1. 所有路径参数必须在DTO的struct tag中使用`path:"paramName"`
2. 所有查询参数必须在DTO的struct tag中使用`query:"paramName"`
3. 所有Body字段使用`json:"fieldName"`
4. 分页参数使用指针类型`*int`以支持可选
5. 从context获取userID：`ctx.Value(constant.CtxKeyUserID).(uint)`
6. 所有错误处理使用`util.WrapHTTPResponse`

## 测试清单

重构完成后，确认以下内容：
- [ ] 代码可以编译通过
- [ ] 所有路由都已注册到huma
- [ ] JWT认证正常工作
- [ ] API文档自动生成（访问/docs）
- [ ] 所有接口返回正确的HTTP状态码
- [ ] 错误处理正确

## 参考资料

- Huma文档: https://github.com/danielgtaylor/huma
- 已完成的服务: tag.go, category.go, user.go, token.go, oauth2.go
- DTO文件: /workspace/internal/protocol/dto/
