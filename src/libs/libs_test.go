package libs

import (
	"testing"
)

func TestGetHTMLContent(t *testing.T) {
	url := "https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel=15904051152"

	s, err := GetHTMLContent(url)
	if err != nil {
		t.Error(err)
	}

	t.Log(s)
}

func TestFormatData(t *testing.T) {
	str := `    
	mts:'1590405',
    province:'辽宁',
    catName:'中国移动',
    telString:'15904051152',
	areaVid:'30498',
	ispVid:'3236139',
	carrier:'辽宁移动'
`

	FormatData(str)
}
