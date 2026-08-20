# BUG_REPRO: 审计范围/大小写与分页 total 错位

## Bug 是什么
审计日志把 GET 与失败请求也写入、Action 变为小写、EntityType 用了完整路径；分页响应 total 恒为 0。

## 如何触发
1. GET /api/v1/activities → 审计表出现 GET 记录（不应有）。
2. POST /api/v1/auth/login（失败）→ 审计表出现 POST 记录（不应有）。
3. 任意写操作 → 审计 Action 为小写 "post"（应为 "POST"）。
4. GET /api/v1/activities（有 2 条）→ total=0（应为 2）。

## 真实错误信息
```
--- FAIL: TestPaginationTotalCorrect
    expected total=2, got body={"data":{"list":[...],"total":0,...}}
--- FAIL: TestAuditActionUppercase
    audit action should be uppercase, found lowercase rows=1
```
