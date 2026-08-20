# BUG_REPRO: 并发下请求 ID 生成与限流器数据竞争

## Bug 是什么
请求 ID 用非原子自增计数（`seq++`），限流器 `buckets` map 在无锁状态下读写，高并发请求触发 data race，可能 panic 且请求 ID 重复。

## 如何触发
并发发起大量请求（含限流路由），例如 8 个 goroutine 同时生成请求 ID / 命中限流。

## 真实错误信息
```
WARNING: DATA RACE
Read at 0x0001048291b8 by goroutine 15:
  gbevent/internal/middleware.newRequestID()
      .../internal/middleware/request_id.go:28 +0xb4
Previous write at 0x0001048291b8 by goroutine 9:
  gbevent/internal/middleware.newRequestID()
      .../internal/middleware/request_id.go:28 +0xcc
```
