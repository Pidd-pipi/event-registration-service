# BUG 复现说明（event-registration__004）

## Bug 是什么
金额保留、数字转字符串、容量显示三个格式化函数联合出错。

## 如何触发
```bash
go test ./internal/service ./internal/util -run 'TestRound2|TestFmtUint|TestFormatters' -count=1
```

## 错误信息
```
--- FAIL: TestFormatters
    util_test.go:72: FormatCapacity = "10/3"
```
