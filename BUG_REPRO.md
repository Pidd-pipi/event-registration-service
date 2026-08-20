# BUG_REPRO: JWT 校验链与 RBAC 状态错位

## Bug 是什么
ParseToken 吞掉解析错误，过期 token 被放行；角色/签发方白名单缺失；RequireRole 判断反转，普通用户可访问管理员路由。

## 如何触发
1. 携带过期签名 token 访问受保护接口 → 应 401，实际 200。
2. 携带未知角色 token → 应 401，实际放行。
3. 携带错误签发方 token → 应 401，实际放行。
4. 普通用户 token GET /api/v1/users → 应 403，实际 200。

## 真实错误信息
```
--- FAIL: TestExpiredTokenRejected
    expected 401 for expired token, got 200 ...
--- FAIL: TestRbacBlocksNormalUserFromAdminList
    expected 403 for normal user on admin route, got 200 ...
```
