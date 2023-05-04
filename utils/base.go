package utils

import (
	"fmt"
	"reflect"
)

// ToMap 结构体转为Map[string]interface{}
func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			// 创建的时候 以下判断的字段 不能为空和0等空值
			// if submit == "create" {
			// } else {
			// 	out[tagValue] = v.Field(i).Interface()
			// }
			if tagValue == "ports" {
				// fmt.Println("ports   === ", v.Field(i).Interface(), v.Field(i).Len())
				if v.Field(i).Len() > 0 {
					out[tagValue] = v.Field(i).Interface()
				}
			} else if tagValue == "conn_limit" {
				if v.Field(i).Interface() != 0 {
					out[tagValue] = v.Field(i).Interface()
				}
			} else if tagValue == "weight" {
				if v.Field(i).Interface() != 0 {
					out[tagValue] = v.Field(i).Interface()
				}
			} else {
				out[tagValue] = v.Field(i).Interface()
			}

		}
	}
	// delete(out, "description") //
	return out, nil
}
