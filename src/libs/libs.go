package libs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/iconv.v1"
)

var ErrNotCharset = errors.New("HTML Not Exisit Charset")

// 获取网页内容
func GetHTMLContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	ctx, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	t, err := getCharset(resp.Header)
	if err != nil {
		t = "gbk"
	}

	str, err := ConvertToUtf8(t, string(ctx))
	if err != nil {
		return "", err
	}

	id := strings.Index(str, "{")
	if id == -1 {
		return "", fmt.Errorf("Index '{' failed")
	}

	return FormatData(str[id : len(str)-2]), nil
}

func IsValidTelNumber(s string) bool {

	b, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, s)

	return b

}

// 获取网页编码格式
func getCharset(h http.Header) (string, error) {

	d := h.Get("content-type")

	regular := `(?i)(.*?;)charset=(.*?)[;| ]`
	re := regexp.MustCompile(regular)

	s := re.FindAllStringSubmatch(d+" ", -1)
	if s == nil {
		return "", fmt.Errorf("Get Charset failed")
	}

	return s[0][2], nil

}

// 转换成utf-8格式
func ConvertToUtf8(t string, ctx string) (string, error) {
	str := ""

	switch strings.ToLower(t) {
	case "gbk":
		cd, err := iconv.Open("utf-8", "gbk")
		if err != nil {
			return "", err
		}

		defer cd.Close()

		str = cd.ConvString(ctx)
	case "utf-8":
		str = ctx
	default:
		return "", ErrNotCharset
	}

	return str, nil
}

// 格式转换
func FormatData(ctx string) string {
	str := "{"

	regular := `(\w.*?):'(.*?)'`
	re := regexp.MustCompile(regular)

	s := re.FindAllStringSubmatch(ctx, -1)
	arr := make([]string, len(s))

	for k, v := range s {
		arr[k] = fmt.Sprintf(`"%s":"%s"`, v[1], v[2])
	}

	str += strings.Join(arr, ",") + "}"

	return str
}
