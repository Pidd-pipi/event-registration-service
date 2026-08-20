// Code generated for scoring tests of record 007.
package z7_test

import (
	"bytes"
	"gbevent/internal/handler"
	"gbevent/internal/testutil"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)




func uploadReq(t *testing.T, filename string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	cfg := testutil.TestConfig(t)
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatal(err)
	}
	w.Close()

	engine := gin.New()
	engine.POST("/api/v1/uploads", handler.NewUploadHandler(cfg, testutil.NewLogger()).UploadImage)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/uploads", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

func TestUploadUnsupportedTypeReturns415(t *testing.T) {
	rec := uploadReq(t, "a.txt", []byte("hello"))
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 for unsupported type, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestUploadTooLargeReturns413(t *testing.T) {
	rec := uploadReq(t, "big.png", bytes.Repeat([]byte("x"), 2*1024*1024))
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413 for oversized file, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestUploadFilenameSanitized(t *testing.T) {
	rec := uploadReq(t, "photo.png", []byte("png-bytes"))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	re := regexp.MustCompile(`"url":"/uploads/[0-9]+\.png"`)
	if !re.MatchString(rec.Body.String()) {
		t.Fatalf("upload url should use server-generated name, got: %s", rec.Body.String())
	}
}
