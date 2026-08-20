# BUG_REPRO: 删除活动未级联清理报名/评论/签到

## Bug 是什么
`ActivityService.Delete` 只删 activities 行，报名/评论/签到记录成为孤儿数据残留。

## 如何触发
1. 给活动创建报名、评论、签到记录。
2. 删除该活动。
3. 关联记录仍可按 activity_id 查到。

## 真实错误信息
```
--- FAIL: TestDeleteActivityRemovesSignups
    expected registrations removed, still 2 left
--- FAIL: TestDeleteActivityRemovesComments
    expected comments removed, still 2 left
--- FAIL: TestDeleteActivityRemovesCheckIns
    expected check-in records removed, still 1 left
```
