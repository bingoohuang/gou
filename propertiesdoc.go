// Package properties is used to read or write or modify the properties document.
package gou

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type elementProperties struct {
	//  #   注释行
	//  !   注释行
	//  ' ' 空白行或者空行
	//  =   等号分隔的属性行
	//  :   冒号分隔的属性行
	typo  byte   //  行类型
	value string //  行的内容,如果是注释注释引导符也包含在内
	key   string //  如果是属性行这里表示属性的key
}

// PropertiesDoc The properties document in memory.
type PropertiesDoc struct {
	elems *list.List
	props map[string]*list.Element
}

// NewProperties is used to create a new and empty properties document.
//
// It's used to generate a new document.
func NewProperties() *PropertiesDoc {
	doc := new(PropertiesDoc)
	doc.elems = list.New()
	doc.props = make(map[string]*list.Element)
	return doc
}

func PrettyProperties(src string) (string, error) {
	doc, _ := LoadProperties(bytes.NewBufferString(src))
	buf := bytes.NewBufferString("")
	err := SavePrettyProperties(doc, buf)
	return buf.String(), err
}

// SavePrettyProperties is used to save the doc to file or stream.
func SavePrettyProperties(doc *PropertiesDoc, writer io.Writer) error {
	var err error

	doc.Accept(func(typo byte, value string, key string) bool {
		switch typo {
		case '#', '!', ' ':
			_, err = fmt.Fprintln(writer, value)
		case '=', ':':
			_, err = fmt.Fprintf(writer, "%s%c%s\n", key, typo, value)
		}

		return nil == err
	})

	return err
}

// LoadProperties is used to create the properties document from a file or a stream.
func LoadProperties(reader io.Reader) (doc *PropertiesDoc, err error) {

	//  创建一个Properties对象
	doc = NewProperties()

	//  创建一个扫描器
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		//  逐行读取
		line := scanner.Bytes()

		//  遇到空行
		if 0 == len(line) {
			doc.elems.PushBack(&elementProperties{typo: ' ', value: string("")})
			continue
		}

		//  找到第一个非空白字符
		pos := bytes.IndexFunc(line, func(r rune) bool {
			return !unicode.IsSpace(r)
		})

		//  遇到空白行
		if -1 == pos {
			doc.elems.PushBack(&elementProperties{typo: ' ', value: string("")})
			continue
		}

		//  遇到注释行
		if '#' == line[pos] {
			doc.elems.PushBack(&elementProperties{typo: '#', value: string(line)})
			continue
		}

		if '!' == line[pos] {
			doc.elems.PushBack(&elementProperties{typo: '!', value: string(line)})
			continue
		}

		//  找到第一个等号的位置
		end := bytes.IndexFunc(line[pos+1:], func(r rune) bool {
			return ('=' == r) || (':' == r)
		})

		//  没有=，说明该配置项只有key
		key := ""
		value := ""
		if -1 == end {
			key = string(bytes.TrimRightFunc(line[pos:], func(r rune) bool {
				return unicode.IsSpace(r)
			}))
		} else {
			key = string(bytes.TrimRightFunc(line[pos:pos+1+end], func(r rune) bool {
				return unicode.IsSpace(r)
			}))

			value = string(bytes.TrimSpace(line[pos+1+end+1:]))
		}

		var typo byte = '='
		if end > 0 {
			typo = line[pos+1+end]
		}
		elem := &elementProperties{typo: typo, key: key, value: value}
		listelem := doc.elems.PushBack(elem)
		doc.props[key] = listelem
	}

	if err = scanner.Err(); nil != err {
		return nil, err
	}

	return doc, nil
}

// Get Retrive the value from PropertiesDoc.
//
// If the item is not exist, the exist is false.
func (p PropertiesDoc) Get(key string) (value string, exist bool) {
	e, ok := p.props[key]
	if !ok {
		return "", ok
	}

	return e.Value.(*elementProperties).value, ok
}

// Set Update the value of the item of the key.
//
// Create a new item if the item of the key is not exist.
func (p *PropertiesDoc) Set(key string, value string) {
	e, ok := p.props[key]
	if !ok {
		p.props[key] = p.elems.PushBack(&elementProperties{typo: '=', key: key, value: value})
		return
	}

	e.Value.(*elementProperties).value = value
	return
}

// Del Delete the exist item.
//
// If the item is not exist, return false.
func (p *PropertiesDoc) Del(key string) bool {
	e, ok := p.props[key]
	if !ok {
		return false
	}

	p.Uncomment(key)
	p.elems.Remove(e)
	delete(p.props, key)
	return true
}

// Comment Append comments for the special item.
//
// Return false if the special item is not exist.
func (p *PropertiesDoc) Comment(key string, comments string) bool {
	e, ok := p.props[key]
	if !ok {
		return false
	}

	//  如果所有注释为空
	if len(comments) <= 0 {
		p.elems.InsertBefore(&elementProperties{typo: '#', value: "#"}, e)
		return true
	}

	//  创建一个新的Scanner
	scanner := bufio.NewScanner(strings.NewReader(comments))
	for scanner.Scan() {
		p.elems.InsertBefore(&elementProperties{typo: '#', value: "#" + scanner.Text()}, e)
	}

	return true
}

// Uncomment Remove all of the comments for the special item.
//
// Return false if the special item is not exist.
func (p *PropertiesDoc) Uncomment(key string) bool {
	e, ok := p.props[key]
	if !ok {
		return false
	}

	for item := e.Prev(); nil != item; {
		del := item
		item = item.Prev()

		if ('=' == del.Value.(*elementProperties).typo) ||
			(':' == del.Value.(*elementProperties).typo) ||
			(' ' == del.Value.(*elementProperties).typo) {
			break
		}

		p.elems.Remove(del)
	}

	return true
}

// Accept Traverse every elementProperties of the document, include comment.
//
// The typo parameter special the elementProperties type.
// If typo is '#' or '!' means current elementProperties is a comment.
// If typo is ' ' means current elementProperties is a empty or a space line.
// If typo is '=' or ':' means current elementProperties is a key-value pair.
// The traverse will be terminated if f return false.
func (p PropertiesDoc) Accept(f func(typo byte, value string, key string) bool) {
	for e := p.elems.Front(); e != nil; e = e.Next() {
		elem := e.Value.(*elementProperties)
		continues := f(elem.typo, elem.value, elem.key)
		if !continues {
			return
		}
	}
}

// Foreach Traverse all of the key-value pairs in the document.
// The traverse will be terminated if f return false.
func (p PropertiesDoc) Foreach(f func(value string, key string) bool) {
	for e := p.elems.Front(); e != nil; e = e.Next() {
		elem := e.Value.(*elementProperties)
		if ('=' == elem.typo) ||
			(':' == elem.typo) {
			continues := f(elem.value, elem.key)
			if !continues {
				return
			}
		}
	}
}

// StringDefault   Retrive the string value by key.
// If the elementProperties is not exist, the def will be returned.
func (p PropertiesDoc) StringDefault(key string, def string) string {
	e, ok := p.props[key]
	if ok {
		return e.Value.(*elementProperties).value
	}

	return def
}

// IntDefault   Retrive the int64 value by key.
// If the elementProperties is not exist, the def will be returned.
func (p PropertiesDoc) IntDefault(key string, def int64) int64 {
	e, ok := p.props[key]
	if ok {
		v, err := strconv.ParseInt(e.Value.(*elementProperties).value, 10, 64)
		if nil != err {
			return def
		}

		return v
	}

	return def
}

// UintDefault Same as IntDefault, but the return type is uint64.
func (p PropertiesDoc) UintDefault(key string, def uint64) uint64 {
	e, ok := p.props[key]
	if ok {
		v, err := strconv.ParseUint(e.Value.(*elementProperties).value, 10, 64)
		if nil != err {
			return def
		}

		return v
	}

	return def
}

// FloatDefault   Retrive the float64 value by key.
// If the elementProperties is not exist, the def will be returned.
func (p PropertiesDoc) FloatDefault(key string, def float64) float64 {
	e, ok := p.props[key]
	if ok {
		v, err := strconv.ParseFloat(e.Value.(*elementProperties).value, 64)
		if nil != err {
			return def
		}

		return v
	}

	return def
}

// BoolDefault   Retrive the bool value by key.
// If the elementProperties is not exist, the def will be returned.
// This function mapping "1", "t", "T", "true", "TRUE", "True" as true.
// This function mapping "0", "f", "F", "false", "FALSE", "False" as false.
// If the elementProperties is not exist of can not map to value of bool,the def will be returned.
func (p PropertiesDoc) BoolDefault(key string, def bool) bool {
	e, ok := p.props[key]
	if ok {
		v, err := strconv.ParseBool(e.Value.(*elementProperties).value)
		if nil != err {
			return def
		}

		return v
	}

	return def
}

// ObjectDefault Map the value of the key to any object.
// The f is the customized mapping function.
// Return def if the elementProperties is not exist of f have a error returned.
func (p PropertiesDoc) ObjectDefault(key string, def interface{}, f func(k string, v string) (interface{}, error)) interface{} {
	e, ok := p.props[key]
	if ok {
		v, err := f(key, e.Value.(*elementProperties).value)
		if nil != err {
			return def
		}

		return v
	}

	return def
}

// String Same as StringDefault but the def is "".
func (p PropertiesDoc) String(key string) string {
	return p.StringDefault(key, "")
}

// Int is ame as IntDefault but the def is 0 .
func (p PropertiesDoc) Int(key string) int64 {
	return p.IntDefault(key, 0)
}

// Uint Same as UintDefault but the def is 0 .
func (p PropertiesDoc) Uint(key string) uint64 {
	return p.UintDefault(key, 0)
}

// Float is same as FloatDefault but the def is 0.0 .
func (p PropertiesDoc) Float(key string) float64 {
	return p.FloatDefault(key, 0.0)
}

// Bool is same as BoolDefault but the def is false.
func (p PropertiesDoc) Bool(key string) bool {
	return p.BoolDefault(key, false)
}

// Object is same as ObjectDefault but the def is nil.
//
// Notice: If the return value can not be assign to nil, this function will panic/
func (p PropertiesDoc) Object(key string, f func(k string, v string) (interface{}, error)) interface{} {
	return p.ObjectDefault(key, interface{}(nil), f)
}
