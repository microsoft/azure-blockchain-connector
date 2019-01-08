package util

import (
	"golang.org/x/oauth2"
	"net/url"
	"reflect"
)

func FieldsToMap(v interface{}, tag string) map[string]string {
	m := make(map[string]string)
	elm := reflect.ValueOf(v).Elem()
	typ := elm.Type()

	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		if f.Type() != reflect.TypeOf("") {
			continue
		}
		k, v := typ.Field(i).Tag.Get(tag), elm.Field(i).String()
		if k == "" || v == "" {
			continue
		}
		m[k] = v
	}
	return m
}

func FieldsToOAuthParams(v interface{}, tag string) []oauth2.AuthCodeOption {
	m := FieldsToMap(v, tag)
	params := make([]oauth2.AuthCodeOption, len(m))
	cnt := 0
	for k, v := range m {
		params[cnt] = oauth2.SetAuthURLParam(k, v)
		cnt++
	}
	return params
}

func FieldsToStringKVPs(v interface{}, tag string) []StringKVP {
	m := FieldsToMap(v, tag)
	params := make([]StringKVP, len(m))
	cnt := 0
	for k, v := range m {
		params[cnt] = StringKVP{k, v}
		cnt++
	}
	return params
}

func AssignURLValues(a, b url.Values) {
	for k, vals := range b {
		for _, val := range vals {
			a.Set(k, val)
		}
	}
}
