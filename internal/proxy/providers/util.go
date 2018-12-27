package providers

import (
	"fmt"
	"reflect"
)

var EnablePrintToken bool

func printToken(tok interface{}) {
	if !EnablePrintToken || tok == nil {
		return
	}
	elm := reflect.ValueOf(tok).Elem()
	if elm.Kind() != reflect.Struct {
		return
	}
	for _, name := range []string{"AccessToken", "RefreshToken"} {
		f := elm.FieldByName(name)
		if f.IsValid() {
			fmt.Println(name+":", f.String())
		}
	}
}
