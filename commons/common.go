package commons

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// TernaryOperation 三元运算
func TernaryOperation(check bool, trueVal, falseVal interface{}) interface{} {
	if check {
		return trueVal
	}
	return falseVal
}

// IsNil checkVal为nil时返回nilVal，否则返回checkVal.(string)
func IsNil(checkVal interface{}, nilVal string) string {
	return TernaryOperation(checkVal == nil, nilVal, checkVal).(string)
}

// GetArrayIndex 获取指定字符串在数组中的位置
func GetArrayIndex(values []string, value string) (int, error) {
	for index, v := range values {
		if v == value {
			return index, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("has no value %s, in array %v", value, values))
}

// FormatFloat64 格式化float64保留指定位小数
func FormatFloat64(v float64, dec int) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf(fmt.Sprintf("%%.%vf", dec), v), 64)
}
