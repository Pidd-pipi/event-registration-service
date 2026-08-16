# BUG 复现说明（event-registration__002）

## Bug 是什么
报名凭证格式校验与报名状态文案错误。

## 如何触发
```bash
go test ./internal/util -run 'TestValidateVoucherFormatTable|TestFormatters' -count=1
```

## 错误信息
```
--- FAIL: TestValidateVoucherFormatTable
    util_test.go:59: ValidateVoucherFormat("GB123") = true, want false
--- FAIL: TestFormatters
    util_test.go:69: RegistrationStatusText = "已取消"
```
