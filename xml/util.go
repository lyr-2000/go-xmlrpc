package xml

import "io"

func _WriteTagString(w io.StringWriter, tagName string, value string) {
	w.WriteString("<")
	w.WriteString(tagName)
	w.WriteString(">")
	w.WriteString(value)
	w.WriteString("</")
	w.WriteString(tagName)
	w.WriteString(">")
}

type XmlRpcMarshaler interface {
	MarshalXmlRpcValue(v interface{}) ([]byte, error)
}
type XmlRpcUnmarshaler interface {
	UnmarshalXmlRpcValue(raw []byte, v interface{}) error
}

type XmlRpcTypeMarshaler interface {
	MarshalXmlRpcElem(v interface{}) ([]byte, error)
}

type XmlRpcTypeUnMarshaler interface {
	UnmarshalXmlRpcElem(v interface{}) ([]byte, error)
}
