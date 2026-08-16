# BUG 复现说明（event-registration__001）

## Bug 是什么
活动状态/类型枚举校验失效，管理员权限判断错误。

## 如何触发
```bash
go test ./internal/service -run 'TestIsOrganizer|TestStatusValidators' -count=1
```

## 错误信息
```
--- FAIL: TestIsOrganizer
    service_test.go:52: admin can manage any: IsOrganizer = false, want true
--- FAIL: TestStatusValidators
    service_test.go:62: unknown status should be invalid
```
