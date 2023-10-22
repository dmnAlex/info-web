package tools

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func Add(a int, b int) int {
	return a + b
}

// returns the value of each field of struct in a slice
func GetFieldValues(s interface{}) []interface{} {
	if s == nil {
		return nil
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out = append(out, v.Field(i).Interface())
	}

	return out
}

// returns the name of each field of struct in a slice
func GetFieldNames(s interface{}) []string {
	if s == nil {
		return nil
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]string, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out = append(out, v.Type().Field(i).Name)
	}

	return out
}

func ToLowerCase(str string) string {
	return strings.ToLower(str)
}

// returns the value of struct's primary key
func GetPrimaryKeyValue(s interface{}) interface{} {
	v := reflect.ValueOf(s)
	var out interface{}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag
		value := tag.Get("gorm")
		if strings.Contains(value, "primaryKey") {
			out = v.Field(i).Interface()
			break
		}
	}

	return out
}

func GetPrimaryKeyName(s interface{}) string {
	v := reflect.ValueOf(s)
	var out string
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag
		value := tag.Get("gorm")
		if strings.Contains(value, "primaryKey") {
			out = v.Type().Field(i).Name
			break
		}
	}

	return out
}

func GetAllTemplates() []string {
	templates := make([]string, 0)
	err := filepath.Walk("internal/template/", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			templates = append(templates, path)
		}
		return err
	})

	if err != nil {
		log.Fatalf("%v", err)
	}

	return templates
}
