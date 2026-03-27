package errcode

import (
	"errors"
	"strings"
)

// 错误码定义
const (
	CodeSuccess       = 200
	CodeBadRequest    = 400
	CodeUnauthorized = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeServerError   = 500
)

// 错误消息映射
var errMsgMap = map[string]string{
	// 外键约束错误
	"foreign key constraint fails":                   "操作失败：关联数据不存在",
	"fk_pm_product_owner":                           "请先添加负责人",
	"fk_pm_requirement_product":                     "请先选择有效的产品",
	"fk_pm_requirement_creator":                     "创建人不存在",
	"fk_pm_requirement_assignee":                    "负责人不存在",
	"fk_pm_product_owner_id":                        "负责人不存在，请先创建用户",
	"fk_pm_requirement_product_id":                  "产品不存在，请先创建产品",
	"fk_pm_requirement_creator_id":                 "创建人不存在",
	"fk_pm_requirement_assignee_id":                 "负责人不存在",

	// 唯一约束错误
	"duplicate entry":                               "数据已存在，请勿重复添加",
	"unique constraint":                            "数据已存在",

	// 非空约束错误
	"null":                                         "必填字段不能为空",
	"cannot be null":                               "不能为空",

	// 数据截断错误
	"data too long":                                "输入数据过长，请缩短",
	"truncated":                                    "数据被截断，请检查输入长度",

	// 通用错误
	"invalid":                                       "数据格式不正确",
	"out of range":                                 "数据超出范围",
}

// ParseError 解析错误并返回友好消息
func ParseError(err error) string {
	if err == nil {
		return "操作成功"
	}

	errStr := err.Error()
	errStrLower := strings.ToLower(errStr)

	// 先检查具体的错误类型
	for key, msg := range errMsgMap {
		if strings.Contains(errStrLower, strings.ToLower(key)) {
			return msg
		}
	}

	// 检查 MySQL 错误码
	if strings.Contains(errStr, "Error 1452") || strings.Contains(errStr, "Error 1451") {
		// 外键约束错误
		if strings.Contains(errStr, "pm_product") && strings.Contains(errStr, "owner_id") {
			return "请先添加负责人"
		}
		if strings.Contains(errStr, "pm_requirement") && strings.Contains(errStr, "product_id") {
			return "请先选择有效的产品"
		}
		if strings.Contains(errStr, "pm_requirement") && strings.Contains(errStr, "creator_id") {
			return "创建人不存在"
		}
		if strings.Contains(errStr, "pm_requirement") && strings.Contains(errStr, "assignee_id") {
			return "负责人不存在"
		}
		return "操作失败：关联数据不存在"
	}

	if strings.Contains(errStr, "Error 1062") {
		return "数据已存在，请勿重复添加"
	}

	if strings.Contains(errStr, "Error 1048") || strings.Contains(errStr, "Error 1364") {
		return "必填字段不能为空"
	}

	// 如果是自定义错误，直接返回
	var ErrInvalid *InvalidError
	if errors.As(err, &ErrInvalid) {
		return ErrInvalid.Message
	}

	// 默认返回简短错误
	return "操作失败"
}

// InvalidError 自定义无效数据错误
type InvalidError struct {
	Message string
}

func (e *InvalidError) Error() string {
	return e.Message
}

func NewInvalidError(msg string) *InvalidError {
	return &InvalidError{Message: msg}
}
