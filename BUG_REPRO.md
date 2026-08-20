# BUG_REPRO: 未命中的取消收藏/标记已读被静默当成成功

## Bug 是什么
`DeleteByUserActivity` 与 `MarkRead` 不检查 RowsAffected：取消未收藏的活动、标记不存在或他人通知已读都返回成功（200），而不是 404。

## 如何触发
1. DELETE /api/v1/activities/99999/favorite（未收藏）→ 200 “已取消收藏”。
2. POST /api/v1/notifications/99999/read（不存在）→ 200。
3. POST /api/v1/notifications/1/read（他人通知）→ 200。

## 真实错误信息
```
--- FAIL: TestUnfavoriteMissingActivityErrors
    expected 404 removing non-favorited activity, got 200 body={"code":0,"data":null,"message":"已取消收藏"}
--- FAIL: TestMarkReadNonExistentErrors
    expected 404 marking missing notification, got 200 body={"code":0,"data":null,"message":"ok"}
```
