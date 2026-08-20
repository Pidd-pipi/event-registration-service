# BUG_REPRO: 角色注册与路由鉴权缺失（权限错位）

## Bug 是什么
注册可申请 admin 角色、service 层角色白名单缺失、/users、/activities/mine、/comments/mine、/favorites/mine 路由鉴权缺失，普通用户可越权访问管理接口。

## 如何触发
1. POST /api/v1/auth/register 带 role=admin → 成功创建管理员（应 400）。
2. 普通用户 GET /api/v1/users → 200（应 403）。
3. 普通用户 GET /api/v1/activities/mine → 200（应 403）。
4. 未登录 GET /api/v1/comments/mine、/api/v1/favorites/mine → 200（应 401）。

## 真实错误信息
```
--- FAIL: TestUserCannotListUsers
    expected 403 for normal user listing users, got 200 ...
--- FAIL: TestCommentMineRequiresAuth
    expected 401 without token for comment mine, got 200 ...
```
