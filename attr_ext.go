package htmlparser

import (
	"strconv"
)

func (t *HtmlTag) HasAttribute(name string) bool {
	_, exists := t.Attributes[name]
	return exists
}

func (t *HtmlTag) GetAttribute(name string) *HtmlAttribute {
	if attr, exists := t.Attributes[name]; exists {
		return &attr
	}
	return nil
}

func (t *HtmlTag) SetAttribute(name, value string) {
	if t.Attributes == nil {
		t.Attributes = make(map[string]HtmlAttribute)
	}
	attr := HtmlAttribute{
		Name:         name,
		Value:        value,
		IsValueExist: true,
	}
	t.Attributes[name] = attr
}

func (t *HtmlTag) RemoveAttribute(name string) {
	if t.Attributes == nil {
		return
	}
	delete(t.Attributes, name)
}

func (a HtmlAttribute) AsString() string {
	return a.Value
}

func (a HtmlAttribute) AsBool() (bool, error) {
	return strconv.ParseBool(a.Value)
}

func (a HtmlAttribute) AsInt() (int, error) {
	v, err := strconv.ParseInt(a.Value, 10, 0)
	return int(v), err
}

func (a HtmlAttribute) AsInt8() (int8, error) {
	v, err := strconv.ParseInt(a.Value, 10, 8)
	return int8(v), err
}

func (a HtmlAttribute) AsInt16() (int16, error) {
	v, err := strconv.ParseInt(a.Value, 10, 16)
	return int16(v), err
}

func (a HtmlAttribute) AsInt32() (int32, error) {
	v, err := strconv.ParseInt(a.Value, 10, 32)
	return int32(v), err
}

func (a HtmlAttribute) AsInt64() (int64, error) {
	return strconv.ParseInt(a.Value, 10, 64)
}

func (a HtmlAttribute) AsUint() (uint, error) {
	v, err := strconv.ParseUint(a.Value, 10, 0)
	return uint(v), err
}

func (a HtmlAttribute) AsUint8() (uint8, error) {
	v, err := strconv.ParseUint(a.Value, 10, 8)
	return uint8(v), err
}

func (a HtmlAttribute) AsUint16() (uint16, error) {
	v, err := strconv.ParseUint(a.Value, 10, 16)
	return uint16(v), err
}

func (a HtmlAttribute) AsUint32() (uint32, error) {
	v, err := strconv.ParseUint(a.Value, 10, 32)
	return uint32(v), err
}

func (a HtmlAttribute) AsUint64() (uint64, error) {
	return strconv.ParseUint(a.Value, 10, 64)
}

func (a HtmlAttribute) AsFloat32() (float32, error) {
	v, err := strconv.ParseFloat(a.Value, 32)
	return float32(v), err
}

func (a HtmlAttribute) AsFloat64() (float64, error) {
	return strconv.ParseFloat(a.Value, 64)
}
