# BUG_REPRO: 报名状态机转换校验缺失

## Bug 是什么
报名状态转换的多处守卫被移除：已签到/已取消仍可取消、已审核仍可再审、已取消可被签到、取消后无法重新报名、活动结束后仍可审核。

## 如何触发
1. 对 `checked_in` 状态的报名执行 Cancel → 应 409，实际成功。
2. 对 `approved` 状态的报名再次 Review → 应 409，实际成功。
3. 对 `cancelled` 状态的报名签到 → 应 409，实际成功。
4. 取消报名后再 Create 同活动报名 → 应成功，实际报“您已报名过该活动”。
5. 对已结束活动的报名 Review → 应 409，实际成功。

## 真实错误信息
```
--- FAIL: TestCancelCheckedInRegistrationFails
    expected AppError code 40906, got: <nil>
--- FAIL: TestResignupAfterCancelSucceeds
    re-signup after cancel should succeed, got: code=40903 message=您已报名过该活动
```
