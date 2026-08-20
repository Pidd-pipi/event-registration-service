# BUG_REPRO: nil 路径级联导致 panic 与 null 响应

## Bug 是什么
多处零值/nil 路径未防护：AuthRequired 吞掉解析错误后解引用 nil claims 崩溃；评论列表向 nil map 写入 panic；通知列表返回 nil 切片导致 JSON null；ErrorHandler 对空错误列表解引用。

## 如何触发
1. 携带畸形 token 访问受保护接口 → panic/500。
2. GET /api/v1/activities/:id/comments → panic。
3. GET /api/v1/notifications/mine（空）→ data.list 为 null。
4. 任意无 c.Error 的请求 → ErrorHandler panic。

## 真实错误信息
```
[Recovery] panic recovered:
assignment to entry in nil map
.../internal/handler/comment_handler.go:43
	(*CommentHandler).List: data["list"] = list
--- FAIL: TestErrorHandlerNoErrorsPasses
    unexpected panic in error handler: runtime error: invalid memory address or nil pointer dereference
```
