/*
错误处理包

提供统一的错误处理方式：
- 所有 API 错误必须使用 Error() 或 ErrorWithStatus() 方法
- 第一个参数为错误码，message 传空字符串
- 前端通过错误码查找翻译显示

示例：
    errors.Error(w, errors.CodeUsernameExists, "")
    errors.ErrorWithStatus(w, 401, errors.CodeNoPermission, "")

错误码定义见 codes.go
*/
package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response API统一响应结构
type Response struct {
	Code    int         `json:"code"`           // 错误码，0表示成功
	Message string      `json:"message"`        // 错误消息，前端根据code查找翻译
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// Success 返回成功响应
func Success(w http.ResponseWriter, data interface{}) {
	resp := Response{
		Code:    0,
		Message: "",
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode response failed: %v", err)
	}
}

// Error 返回错误响应
// code: 错误码，见 codes.go
// message: 固定传空字符串，由前端根据code查找翻译
func Error(w http.ResponseWriter, code int, message string) {
	resp := Response{
		Code:    code,
		Message: "", // 前端根据 code 查找翻译
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode error response failed: %v", err)
	}
	// 日志记录错误码和中文描述
	if code != 0 {
		log.Printf("API Error: code=%d, message=%s", code, GetMessage(code))
	}
}

// ErrorWithStatus 返回带HTTP状态码的错误响应
// 用于需要返回特定HTTP状态码的场景（如401未认证、403禁止等）
func ErrorWithStatus(w http.ResponseWriter, httpStatus int, code int, message string) {
	resp := Response{
		Code:    code,
		Message: "",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode error response failed: %v", err)
	}
	if code != 0 {
		log.Printf("API Error: httpStatus=%d, code=%d, message=%s", httpStatus, code, GetMessage(code))
	}
}

// ErrorWithData 返回带数据的错误响应
// 用于某些需要返回额外数据的错误场景
func ErrorWithData(w http.ResponseWriter, code int, data interface{}) {
	resp := Response{
		Code:    code,
		Message: "",
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode error response failed: %v", err)
	}
	log.Printf("API Error: code=%d, message=%s, data=%+v", code, GetMessage(code), data)
}
