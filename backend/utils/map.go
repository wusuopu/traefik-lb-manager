package utils

import (
	"reflect"
	"strings"
)

// 像 js lodash 一样，获取 map 的值
func MapDeepGet (data interface{}, path string, defaultValue interface{}) interface{} {
	keys := strings.Split(path, ".")

	// 递归遍历 map
	var current interface{} = data
	for _, key := range keys {
		// 获取当前值的类型
		val := reflect.ValueOf(current)

		// 如果当前值是 nil 或者不是 map 类型，返回错误
		if !val.IsValid() || val.Kind() != reflect.Map {
			return defaultValue
		}

		// 获取 map 中对应键的值
		mapKey := reflect.ValueOf(key)
		current = val.MapIndex(mapKey).Interface()
	}

	return current
}