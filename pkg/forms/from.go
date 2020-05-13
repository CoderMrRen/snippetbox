package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

//EmailRX 邮箱正则表达式
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//Form Form
type Form struct {
	url.Values
	Errors errors
}

//New create Form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

//Required 检查表单中的特定字段
//检测指定字段是否为空，为空添加消息显示表单不能为空。
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

//MaxLength 查表单中的特定字段指允许最大字符长度。如果大于指定长度，则添加相应的消息以解决表单错误。
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field to long(maximum id %d)", d))
	}
}

//MinLength 查表单中的特定字段指允许小于字符长度。如果小于指定长度，则添加相应的消息以解决表单错误。
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field to short(minimum id %d)", d))
	}
}

//PermittedValues 以检查表单中的特定字段,匹配一组特定的允许值之一。如果检查失败将错误消息添加到表单错误中。
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

//MatchesPattern 实现MatchesPattern方法以检查表单中的特定字段匹配正则表达式。如果检查失败，则添加相应的消息以解决表单错误
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

//Valid 判断是否有错，没有错返回true,有错返回False
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
