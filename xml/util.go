package xml

import "io"

// type Swriter interface {
// 	WriteString(s string) int
// }
// func floatString(f float)

func _WriteTagString(w io.StringWriter, tagName string, value string) {
	w.WriteString("<")
	w.WriteString(tagName)
	w.WriteString(">")
	w.WriteString(value)
	w.WriteString("</")
	w.WriteString(tagName)
	w.WriteString(">")
}

type XmlRpcParser interface {
	MarshalRpcXml(v interface{}) ([]byte, error)
	UnmarshalRpcXml(raw []byte, v interface{}) error
}
