// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func rpcRequest2XML(method string, rpc interface{}) (string, error) {
	buffer := "<methodCall><methodName>"
	buffer += method
	buffer += "</methodName>"
	params, err := rpcParams2XML(rpc)
	buffer += params
	buffer += "</methodCall>"
	return buffer, err
}

func rpcResponse2XML(rpc interface{}) (string, error) {
	buffer := "<methodResponse>"
	params, err := rpcParams2XML(rpc)
	buffer += params
	buffer += "</methodResponse>"
	return buffer, err
}

func rpcParams2XML(rpc interface{}) (string, error) {
	var err error
	buffer := "<params>"
	for i := 0; i < reflect.ValueOf(rpc).Elem().NumField(); i++ {
		var xml string
		buffer += "<param>"
		xml, err = rpc2XML(reflect.ValueOf(rpc).Elem().Field(i).Interface())
		buffer += xml
		buffer += "</param>"
	}
	buffer += "</params>"
	return buffer, err
}

func rpc2XML(value interface{}) (string, error) {
	// if enc, ok := value.(interface {
	// 	MarshalRpcXml(v interface{}) ([]byte, error)
	// }); ok {
	// 	bs, err := enc.MarshalRpcXml(value)
	// 	return string(bs), err
	// }
	if enc, ok := value.(XmlRpcMarshaler); ok {
		res, err := enc.MarshalXmlRpcValue(value)
		return string(res), err
	}
	var buf = &strings.Builder{}
	buf.WriteString("<value>")
	switch curr := value.(type) {
	case int:
		buf.WriteString("<int>")
		buf.WriteString(strconv.Itoa(value.(int)))
		buf.WriteString("</int>")
	case float64:
		s := strconv.FormatFloat(value.(float64), 'f', 6, 64)
		buf.WriteString("<double>")
		buf.WriteString(s)
		buf.WriteString("</double>")
	case string:
		buf.WriteString(string2XML(value.(string)))
	case bool:
		buf.WriteString(bool2XML(value.(bool)))
	case []byte:
		buf.WriteString(base642XML(value.([]byte)))
	case []interface{}, []int, []float64, []string:
		buf.WriteString(array2XML(value))
	case time.Time:
		buf.WriteString(time2XML(curr))
	default:
		if value == nil /* || reflect.ValueOf(value).IsNil() */ {
			buf.WriteString("<nil/>")
		} else {
			rev := reflect.ValueOf(value)
			switch rev.Kind() {
			case reflect.Array, reflect.Slice:
				buf.WriteString(array2XML(value))
			case reflect.Pointer:
				if rev.IsNil() {
					buf.WriteString("<nil/>")
				} else {
					buf.WriteString(struct2XML(value))
				}
			default:
				buf.WriteString(struct2XML(value))
			}
		}
	}
	buf.WriteString("</value>")
	return buf.String(), nil
}

func bool2XML(value bool) string {
	var b string
	if value {
		b = "1"
	} else {
		b = "0"
	}
	return fmt.Sprintf("<boolean>%s</boolean>", b)
}

func string2XML(value string) string {
	value = strings.Replace(value, "&", "&amp;", -1)
	value = strings.Replace(value, "\"", "&quot;", -1)
	value = strings.Replace(value, "<", "&lt;", -1)
	value = strings.Replace(value, ">", "&gt;", -1)
	return fmt.Sprintf("<string>%s</string>", value)
}

func struct2XML(value interface{}) (out string) {
	buf := strings.Builder{}
	buf.WriteString("<struct>")
	for i := 0; i < reflect.TypeOf(value).NumField(); i++ {
		field := reflect.ValueOf(value).Field(i)
		field_type := reflect.TypeOf(value).Field(i)
		var name string
		if field_type.Tag.Get("xml") != "" {
			name = field_type.Tag.Get("xml")
		} else {
			name = field_type.Name
		}
		field_value, _ := rpc2XML(field.Interface())
		field_name := fmt.Sprintf("<name>%s</name>", name)
		buf.WriteString(fmt.Sprintf("<member>%s%s</member>", field_name, field_value))
	}
	buf.WriteString("</struct>")
	out = buf.String()
	return
}

func array2XML(value interface{}) (out string) {
	var buf = strings.Builder{}
	buf.WriteString("<array><data>")
	ref := reflect.ValueOf(value)
	n := ref.Len()
	for i := 0; i < n; i++ {
		item_xml, _ := rpc2XML(ref.Index(i).Interface())
		buf.WriteString(item_xml)
	}
	buf.WriteString("</data></array>")
	out = buf.String()
	return
}

func time2XML(t time.Time) string {
	/*
		// TODO: find out whether we need to deal
		// here with TZ
		var tz string;
		zone, offset := t.Zone()
		if zone == "UTC" {
			tz = "Z"
		} else {
			tz = fmt.Sprintf("%03d00", offset / 3600 )
		}
	*/
	return fmt.Sprintf("<dateTime.iso8601>%04d%02d%02dT%02d:%02d:%02d</dateTime.iso8601>",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func base642XML(data []byte) string {
	str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("<base64>%s</base64>", str)
}
