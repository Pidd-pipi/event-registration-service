# BUG 复现说明（event-registration__003）

## Bug 是什么
活动列表分页归一化与 offset 错乱。

## 如何触发
```bash
go test ./internal/dto -run 'TestPageQueryNormalizeDefaults' -count=1
```

## 错误信息
```
--- FAIL: TestPageQueryNormalizeDefaults
    pagination_combination_test.go:9: Page = 0, want 1
```
