# BUG_REPRO: 上传错误码丢失且文件名未服务端生成

## Bug 是什么
上传校验错误被降级成通用 400（不支持类型应 415、超限应 413），handler 无法区分；保存文件名沿用客户端文件名。

## 如何触发
1. 上传 .txt 文件 → 400 “Upload image failed: unsupported file type: .txt”（应为 415）。
2. 上传超过 1MB 的文件 → 400（应为 413）。
3. 上传 photo.png → 返回 url 为 /uploads/photo.png（应为服务端时间戳文件名）。

## 真实错误信息
```
--- FAIL: TestUploadUnsupportedTypeReturns415
    expected 415 for unsupported type, got 400 ...
--- FAIL: TestUploadFilenameSanitized
    upload url should use server-generated name, got: {"code":0,"data":{"url":"/uploads/photo.png"},"message":"ok"}
```
