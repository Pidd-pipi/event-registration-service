package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gbevent/internal/constants"
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

// SaveUploadedImage 保存上传图片到 uploadDir，返回可访问的相对路径。
// 不允许类型与超限均以 AppError 形式返回，便于 handler 区分状态码；
// 落盘文件名由服务端生成，避免直接使用用户提供的原始文件名。
func SaveUploadedImage(uploadDir string, maxMB int64, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		return "", NewAppError(constants.CodeUnsupportedType, constants.MsgUnsupportedFileType)
	}
	if file.Size > maxMB*1024*1024 {
		return "", NewAppError(constants.CodeUploadTooLarge, constants.MsgUploadTooLarge)
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}
	name := generateUploadName(ext)
	dst := filepath.Join(uploadDir, name)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open upload file: %w", err)
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create destination file: %w", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		return "", fmt.Errorf("write upload file: %w", err)
	}
	return "/uploads/" + name, nil
}

// generateUploadName 生成服务端文件名：时间戳 + 6 位随机数 + 扩展名。
func generateUploadName(ext string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1_000_000))
	return fmt.Sprintf("%s%06d%s", time.Now().Format("20060102150405"), n.Int64(), ext)
}
