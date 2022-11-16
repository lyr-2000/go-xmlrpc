// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type SubStructRpc2Xml struct {
	Foo  int
	Bar  string
	Data []int
}

type StructRpc2Xml struct {
	Int    int
	Float  float64
	Str    string
	Bool   bool
	Sub    SubStructRpc2Xml
	Time   time.Time
	Base64 []byte
}

func TestRPC2XML(t *testing.T) {
	req := &StructRpc2Xml{123, 3.145926, "Hello, World!", false, SubStructRpc2Xml{42, "I'm Bar", []int{1, 2, 3}}, time.Date(2012, time.July, 17, 14, 8, 55, 0, time.Local), []byte("you can't read this!")}
	xml, err := rpcRequest2XML("Some.Method", req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodCall><methodName>Some.Method</methodName><params><param><value><int>123</int></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param><param><value><dateTime.iso8601>20120717T14:08:55</dateTime.iso8601></value></param><param><value><base64>eW91IGNhbid0IHJlYWQgdGhpcyE=</base64></value></param></params></methodCall>"
	if xml != expected {

		t.Error("RPC2XML conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)

		os.WriteFile("./expect.xml", []byte(expected), 0666)
		os.WriteFile("./curr.xml", []byte(xml), 0666)
	}
}

type StructSpecialCharsRpc2Xml struct {
	String1 string
}

func TestRPC2XMLSpecialChars(t *testing.T) {
	req := &StructSpecialCharsRpc2Xml{" & \" < > "}
	xml, err := rpcResponse2XML(req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><string> &amp; &quot; &lt; &gt; </string></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

type StructNilRpc2Xml struct {
	Ptr *int
}

func TestRpc2XmlNil(t *testing.T) {
	req := &StructNilRpc2Xml{nil}
	xml, err := rpcResponse2XML(req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><nil/></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

func struct2XML_test1(value interface{}) (out string) {
	out += "<struct>"
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
		out += fmt.Sprintf("<member>%s%s</member>", field_name, field_value)
	}
	out += "</struct>"
	return
}

func struct2XML_test2(value interface{}) (out string) {
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

type User struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Content string `json:"content"`
	Title   string `json:"title"`
	ID      int    `json:"id"`
}

func Benchmark_plusstring(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := User{ID: i, Title: strings.Repeat("aaa", 100)}

		_ = struct2XML_test1(user)
	}
	// _ = s
}
func Benchmark_stringsbuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// s := ""
		user := User{ID: i, Title: strings.Repeat("aaa", 100)}
		_ = struct2XML_test2(user)
	}

}
