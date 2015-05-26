package helper

import (
	"time"

	"github.com/skynology/cloud-types"
	"github.com/skynology/go-sdk"
)

// 用于云代码中调用其他函数
type FuncHandler interface {
	Call(h *Helper, req types.CloudRequest, app *skynology.App, name string, data map[string]interface{}) (result map[string]interface{}, err error)
}

// 创建新Helper
func NewHelper(handler FuncHandler) *Helper {
	return &Helper{
		HideFields:    []string{},
		ProtectFields: []string{},
		Logs:          []types.CloudLog{},
		Data:          make(map[string]interface{}),
		funcHandler:   handler,
	}
}

// Helper
type Helper struct {
	Type          string
	IsFunction    bool
	HideFields    []string
	ProtectFields []string
	Result        interface{}
	Data          map[string]interface{}
	Logs          []types.CloudLog
	Errors        types.CloudError
	funcHandler   FuncHandler
}

// 云代码中调用函数
func (h *Helper) Call(req types.CloudRequest, app *skynology.App, name string, data map[string]interface{}) (result map[string]interface{}, err error) {
	result, err = h.funcHandler.Call(h, req, app, name, data)
	return
}

// 完成并退出
// result 参数必须为可转为Json的对象
func (h *Helper) Render(result interface{}) {
	// 不允许在表hook中使用
	if h.Type == "hook" {
		return
	}
	h.Result = result
}

// 更新字段
func (h *Helper) Set(fieldName string, value interface{}) {
	h.Data[fieldName] = value
}

// 取消并返回
// 返回 RESTFul Error格式,
// Code: -1, Error: message
func (h *Helper) Cancel(message string) {
	h.Errors = types.CloudError{Code: -1, Message: message}
}

// 取消并返回
// 返回 RESTFul Error格式,
// Code: code, Error: message
func (h *Helper) CancelWithCode(message string, code int) {
	h.Errors = types.CloudError{Code: code, Message: message}
}

// 隐藏字段。不返回给API调用者
func (h *Helper) Hide(fieldName string) {
	h.HideFields = append(h.HideFields, fieldName)
}

// 保护字段， 设置过的字段将无法更新
func (h *Helper) Protect(fieldName string) {
	h.ProtectFields = append(h.ProtectFields, fieldName)
}

// 打印log, 方便调试
func (h *Helper) Log(message string) {
	h.LogWithFlag(message, "info")
}

func (h *Helper) LogWithFlag(message string, flag string) {
	l := types.CloudLog{}
	l.Content = message
	l.CreatedAt = time.Now().Format(time.RFC3339Nano)
	l.Flag = flag
	h.Logs = append(h.Logs, l)
}
