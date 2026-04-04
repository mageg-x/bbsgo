package utils

import (
	"bbsgo/errors"
	"encoding/json"
	"log"
	"net/http"
)

// Response 统一响应结构（已废弃，请使用 errors.Response）
// 本文件保留仅用于兼容旧代码，新代码请直接使用 errors 包
type Response struct {
	Code    int         `json:"code"`           // 状态码：0表示成功，其他表示错误
	Message string      `json:"message"`        // 响应消息（已废弃，统一使用错误码）
	Data    interface{} `json:"data,omitempty"` // 响应数据，可选
}

// Success 返回成功响应
// w: HTTP 响应写入器
// data: 要返回的数据对象
func Success(w http.ResponseWriter, data interface{}) {
	log.Printf("Success: [DEBUG] start, data type=%T, data=%+v", data, data)
	resp := Response{
		Code:    0,
		Message: "",
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Success: [DEBUG] before Encode, resp=%+v", resp)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode response failed: %v", err)
	}
	log.Printf("Success: [DEBUG] after Encode")
}

// Error 返回错误响应
// w: HTTP 响应写入器
// code: 错误码（已废弃，请使用 errors.CodeXXX）
// message: 错误消息（已废弃，固定传空字符串，由前端根据code查找翻译）
func Error(w http.ResponseWriter, code int, message string) {
	errors.Error(w, code, message)
}

// ErrorWithStatus 返回带 HTTP 状态码的错误响应
// w: HTTP 响应写入器
// httpStatus: HTTP 状态码（如 401, 403, 500 等）
// code: 业务错误码（已废弃，请使用 errors.CodeXXX）
// message: 错误消息（已废弃，固定传空字符串）
func ErrorWithStatus(w http.ResponseWriter, httpStatus int, code int, message string) {
	errors.ErrorWithStatus(w, httpStatus, code, message)
}
