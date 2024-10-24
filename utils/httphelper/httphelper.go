package httphelper

import (
	"bytes"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

/*
required    //这将验证该值不是数据类型的默认零值。数字不为０，字符串不为 " ", slices, maps, pointers, interfaces, channels and functions 不为 nil
isdefault    //这验证了该值是默认值，几乎与所需值相反。
len=10    //对于数字，长度将确保该值等于给定的参数。对于字符串，它会检查字符串长度是否与字符数完全相同。对于切片，数组和map，验证元素个数。
max=10    //对于数字，max将确保该值小于或等于给定的参数。对于字符串，它会检查字符串长度是否最多为该字符数。对于切片，数组和map，验证元素个数。
min=10
eq=10    //对于字符串和数字，eq将确保该值等于给定的参数。对于切片，数组和map，验证元素个数。
ne=10    //和eq相反
oneof=red green (oneof=5 7 9)    //对于字符串，整数和uint，oneof将确保该值是参数中的值之一。参数应该是由空格分隔的值列表。值可以是字符串或数字。
gt=10    //对于数字，这将确保该值大于给定的参数。对于字符串，它会检查字符串长度是否大于该字符数。对于切片，数组和map，它会验证元素个数。
gt    //对于time.Time确保时间值大于time.Now.UTC（）
gte=10    //大于等于
gte    //对于time.Time确保时间值大于或等于time.Now.UTC（）
lt=10    //小于
lt    //对于time.Time确保时间值小于time.Now.UTC（）
lte=10    //小于等于
lte    //对于time.Time确保时间值小于等于time.Now.UTC（）
－－－－

unique    //对于数组和切片，unique将确保没有重复项。对于map，unique将确保没有重复值。
alpha    //这将验证字符串值是否仅包含ASCII字母字符
alphanum    //这将验证字符串值是否仅包含ASCII字母数字字符
alphaunicode    //这将验证字符串值是否仅包含unicode字符
alphanumunicode    //这将验证字符串值是否仅包含unicode字母数字字符
numeric    //这将验证字符串值是否包含基本数值。基本排除指数等...对于整数或浮点数，它返回true。
hexadecimal    //这将验证字符串值是否包含有效的十六进制
hexcolor    //这验证字符串值包含有效的十六进制颜色，包括＃标签（＃）
rgb    //这将验证字符串值是否包含有效的rgb颜色
rgba    //这将验证字符串值是否包含有效的rgba颜色
hsl    //这将验证字符串值是否包含有效的hsl颜色
hsla    //这将验证字符串值是否包含有效的hsla颜色
email    //这验证字符串值包含有效的电子邮件这可能不符合任何rfc标准的所有可能性，但任何电子邮件提供商都不接受所有可能性
file    //这将验证字符串值是否包含有效的文件路径，并且该文件存在于计算机上。这是使用os.Stat完成的，它是一个独立于平台的函数。
url    //这会验证字符串值是否包含有效的url这将接受golang请求uri接受的任何url，但必须包含一个模式，例如http：//或rtmp：//
uri    //这验证了字符串值包含有效的uri。这将接受uri接受的golang请求的任何uri
base64    //这将验证字符串值是否包含有效的base64值。虽然空字符串是有效的base64，但这会将空字符串报告为错误，如果您希望接受空字符串作为有效字符，则可以将此字符串与omitempty标记一起使用。
base64url    //这会根据RFC4648规范验证字符串值是否包含有效的base64 URL安全值。尽管空字符串是有效的base64 URL安全值，但这会将空字符串报告为错误，如果您希望接受空字符串作为有效字符，则可以将此字符串与omitempty标记一起使用。
btc_addr    //这将验证字符串值是否包含有效的比特币地址。检查字符串的格式以确保它匹配P2PKH，P2SH三种格式之一并执行校验和验证
btc_addr_bech32    //这验证了字符串值包含bip-0173定义的有效比特币Bech32地址（https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki）特别感谢Pieter Wuille提供的参考实现。
eth_addr    //这将验证字符串值是否包含有效的以太坊地址。检查字符串的格式以确保它符合标准的以太坊地址格式完全验证被https://github.com/golang/crypto/pull/28阻止
contains=@    //这将验证字符串值是否包含子字符串值
containsany=!@#?    //这将验证字符串值是否包含子字符串值中的任何Unicode code points。
containsrune=@    //这将验证字符串值是否包含提供的符文值。
excludes=@    //这验证字符串值不包含子字符串值。
excludesall=!@#?    //这将验证字符串值在子字符串值中是否包含任何Unicode code points。
*/
// 验证示例模板
type validateTemplate struct {
	Id         string `json:"id" validate:“required”`
	Age        string `json:"age" validate:“gte=1,lte=30”`
	Email      string `json:"email" validate:“required,email”`
	Password   string `json:"password" validate:“required”`
	RePassword string `json:"re_password" validate:“required,eqfield=Password”`
	Date       string `json:"date" validate:“required,datetime=2000-01-01,checkDate”`
	Tag        any    `json:"tag,omitempty"`
	//表示字段值符合日期类型，如果datetime后边不接=?，那么默认为Y-m-d H:i:s，否则验证器会按照指定格式判断，比如 datetime=Y-m、datetime=Y/m/d H:i:s等，可以是Y m d H i s 的随意拼接
}

// 支持验证的http request请求id
type RequestId struct {
	Id string `form:"id" validate:“required”`
}

// 支持验证的http json request请求id
type RequestJsonId struct {
	Id string `json:"id" validate:“required”`
}

type ResponseJsonBody struct {
	Code    int         `json:"code"`             //错误码  没有错误返回0，其他根据业务要求定义
	Message string      `json:"message"`          //设置应答的格式
	Result  interface{} `json:"result,omitempty"` //interface{} 也可以用any代替
}

// 应答，应答结果为json结构
func Response(w http.ResponseWriter, code int, message string, data interface{}) {
	httpx.OkJson(w, ResponseJsonBody{
		Code:    code,
		Message: message,
		Result:  data,
	})
}

// 应答，应答结果为字符串
func ResponseString(w http.ResponseWriter, result string) {
	w.Write(bytes.NewBufferString(result).Bytes())
}

// 成功的请求
func Success(w http.ResponseWriter, data interface{}) {
	Response(w, 0, "请求成功！！", data)
}

// 处理http失败信息，返回自定义消息
func FailMessage(w http.ResponseWriter, message string) {
	Response(w, 1, message, nil)
}

// 处理http失败信息
func Fail(w http.ResponseWriter, data interface{}) {
	Response(w, 1, "请求失败！！", nil)
}

// 根据结果进行分发处理
func HandleResult(w http.ResponseWriter, result interface{}, err interface{}) {
	if err != nil {
		Fail(w, err)
	} else {
		Success(w, result)
	}
}

// 根据规范进行请求信息验证
func ValidateRequest(ctx context.Context, req http.Request) (err error) {
	return validator.New().StructCtx(ctx, req)
}
