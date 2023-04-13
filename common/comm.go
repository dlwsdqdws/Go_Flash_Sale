package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// Map data to a struct based on the SQL tag in the struct and convert the types.
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		// Retrieve the value corresponding to the SQL.
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		// Retrieve the name of the corresponding field.
		name := objValue.Type().Field(i).Name
		// Retrieve the type of the corresponding field.
		structFieldType := objValue.Field(i).Type()
		// Retrieve the type of the variable, or it can be directly specified as string type.
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			// Type conversion
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {

			}
		}
		// Set the value of a variable with a specific type.
		objValue.FieldByName(name).Set(val)
	}
}

// Type conversion
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	return reflect.ValueOf(value), errors.New("Unknown：" + ntype)
}
