# cy-306 活动报名通 执行记录

- 项目编号/名称：cy-306 活动报名通（gbevent）
- 日期：2026-08-16
- 短名：gbevent
- 端口：前端 18506 / 后端 19506 / MySQL 57506
- 技术栈：Vue 3 + TypeScript + Element Plus + Vite + Pinia；Go 1.22 + Gin + GORM；MySQL 8.0；JWT + RBAC

## Docker Compose 结果

| 容器 | 状态 | 端口 |
| --- | --- | --- |
| gbevent-db | Up (healthy) | 57506:3306 |
| gbevent-backend | Up (healthy) | 19506:8080 |
| gbevent-frontend | Up | 18506:80 |

`docker compose config --quiet` 通过；`docker compose up -d --build` 一键启动成功。

## 关键 API 冒烟结果（42/42 通过）

| 接口 | 方法 | 状态码 | 结果摘要 |
| --- | --- | --- | --- |
| /healthz | GET | 200 | ok |
| /api/healthz（Nginx 反代） | GET | 200 | ok |
| /api/v1/healthz（Nginx 反代） | GET | 200 | ok |
| /api/v1/auth/register | POST | 200 | 注册成功 |
| /api/v1/auth/login | POST | 200 | 返回 JWT + 用户 |
| /api/v1/users/me | GET | 200 | 当前用户 |
| /api/v1/users/me（未授权） | GET | 401 | 未登录 |
| /api/v1/activities/mine（普通用户） | GET | 403 | 越权拦截 |
| /api/v1/activities?page_size=5 | GET | 200 | 活动列表 |
| /api/v1/activities/calendar?month=2026-08 | GET | 200 | 日历活动 |
| /api/v1/activities/1 | GET | 200 | 活动详情 |
| /api/v1/activities/1/comments | GET | 200 | 评论+平均分 |
| /api/v1/activities/1/comments | POST | 200 | 发表评论 |
| /api/v1/registrations | POST | 200 | 在线报名成功 |
| /api/v1/registrations（重复） | POST | 409 | 已报名过该活动 |
| /api/v1/registrations/mine | GET | 200 | 我的报名 |
| /api/v1/registrations/:id/cancel | POST | 200 | 取消成功 |
| /api/v1/registrations/:id/cancel（重复） | POST | 409 | 状态冲突 |
| /api/v1/activities/:id/favorite | POST | 200 | 收藏成功 |
| /api/v1/activities/:id/favorite（重复） | POST | 409 | 已收藏 |
| /api/v1/favorites/mine | GET | 200 | 收藏列表 |
| /api/v1/notifications/mine | GET | 200 | 通知列表 |
| /api/v1/notifications/:id/read | POST | 200 | 已读 |
| /api/v1/notifications/read-all | POST | 200 | 全部已读 |
| /api/v1/activities | POST | 200 | 组织者创建活动 |
| /api/v1/activities/:id/publish | POST | 200 | 发布成功 |
| /api/v1/activities/:id/publish（重复） | POST | 409 | 状态冲突 |
| /api/v1/activities/:id/stats | GET | 200 | 报名/签到统计 |
| /api/v1/activities/mine | GET | 200 | 我的活动 |
| /api/v1/registrations/:id/review | POST | 200 | 审核通过 |
| /api/v1/registrations/:id/review（重复） | POST | 409 | 审核状态冲突 |
| /api/v1/check-ins?activity_id=3（凭证） | POST | 200 | 签到成功 |
| /api/v1/check-ins（重复签到） | POST | 409 | 已签到 |
| /api/v1/check-ins（扫码） | POST | 409 | 已签到拦截 |
| /api/v1/check-ins（无效凭证） | POST | 404 | 凭证号无效 |
| /api/v1/check-ins?activity_id=3 | GET | 200 | 签到记录 |
| /api/v1/registrations/export?activity_id=3 | GET | 200 | CSV 导出 |
| /api/v1/auth/login（admin） | POST | 200 | 管理员登录 |
| /api/v1/users?page_size=5（admin） | GET | 200 | 用户列表 |
| /api/v1/users/me | PUT | 200 | 更新资料 |

异常路径覆盖：401 未授权、403 越权、404 无效凭证、409 重复报名/取消/审核/发布/签到/收藏。

## 浏览器验证结论（内置 Playwright，无外部 Chrome）

- 首页 /calendar 打开正常，月历展示「2026-08-23 1 场」等真实活动计数，数据来自后端 /api/v1/activities/calendar。
- /activities 列表渲染 5 张真实活动卡片（Go 语言企业级开发实战讲座、新员工安全培训、秋季团队趣味运动会、黑客松编程竞赛（草稿）、上季度读书分享会（已结束）），含类型标签与「名额 x/y」。
- 登录流程可用：/login 输入 organizer/User@123 点击登录后进入日历，头部显示「活动组织者」并出现发布管理/报名与签到菜单。
- /organizer/activities 渲染真实活动表格与状态（已发布/草稿/已结束）。
- /organizer/registrations 页面正常（报名名单/签到面板 Tab、表格表头、导出/补录按钮）。
- /profile 渲染 5 个 Tab（基本资料/我的报名/我的收藏/我的评论/消息通知）与用户资料。
- 截图：output/gbevent_calendar.png、output/gbevent_activities.png、output/gbevent_organizer_activities.png、output/gbevent_calendar2.png。

## README 检查项

- Docker Compose 一键启动命令在最前；本地开发命令；技术栈表格（后端 Go 1.22 + Gin + GORM）；目录结构；环境变量；部署说明；License。
- 枚举出现位置清单：ActivityStatus、RegistrationStatus、ActivityType 前后端出现位置已列出。

## 其他质量项

- 后端 `go build ./...` 通过；`go vet` 通过。
- 单元测试：internal/service 与 internal/util 表驱动测试通过（go test ./... ok）。
- 前端 `npm run build` 零错误（vue-tsc + vite build 通过）。
- database/init.sql 含建表与种子数据（3 用户 + 5 活动 + 3 报名 + 签到/评论/收藏/通知），容器首次启动自动执行。

- 提交记录：init commit 81ae99d；本文件独立提交 docs commit（见 git log）。
