# BUG_REPRO: 缺失用户/活动/报名被误报 500

## Bug 是什么
repository 层把 `ErrNotFound` 哨兵降级成了普通 `fmt.Errorf`，错误链断裂后 `errors.Is` 恒为 false，导致「查不到的用户/活动/报名」被当成服务器内部错误返回 500，而不是 404/业务错误。

## 如何触发
1. 注册一个新手机号（此前未注册过）→ 报“服务器内部错误”。
2. 用不存在的用户名登录 → 500。
3. GET /api/v1/activities/99999（不存在）→ 500。
4. POST /api/v1/registrations/99999/cancel → 500。

## 真实错误信息
```
go test ./z1 -run '^TestFreshUserSignupSucceeds$' -count=1
--- FAIL: TestFreshUserSignupSucceeds
    error_chain_test.go: register new phone should succeed, got: User register check username failed: username "freshman" not found
GET /api/v1/activities/99999 -> HTTP 500 {"code":50000,"data":null,"message":"服务器内部错误"}
```
