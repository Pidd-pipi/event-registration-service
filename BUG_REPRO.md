# BUG 复现说明（event-registration__005）

## Bug 是什么
活动类型文案与签到码解析错误。

## 如何触发
```bash
go test ./internal/util -run 'TestSscanUint64|TestFormatters'
```

## 错误信息
```
--- FAIL: TestFormatters
    util_test.go:75: ActivityTypeText = "聚会"
--- FAIL: TestSscanUint64
    util_test.go:96: SscanUint64("3") = 0, want 3
```
