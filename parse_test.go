package linetmpl

import (
	"strings"
	"testing"

	"github.com/tj/assert"
)

func TestCompile(t *testing.T) {
	tpl, err := Compile("test", `/[.日期]/[.车牌]/[.年]/[.月]/[.日]/[.违法代码]_[.序号].jpg`)
	assert.NoError(t, err)
	assert.NotNil(t, tpl)

	var (
		sb  strings.Builder
		ctx = map[string]interface{}{
			"日期":   "20180201",
			"车牌":   "苏A12345",
			"年":    "2018",
			"月":    "02",
			"日":    "01",
			"违法代码": "01",
			"序号":   "1",
		}
	)

	err = tpl.ExecuteTemplate(&sb, "test", ctx)
	assert.NoError(t, err)
	assert.Equal(t, "/20180201/苏A12345/2018/02/01/01_1.jpg", sb.String())
}
