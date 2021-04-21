package notice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/li4n0/revsuit/internal/record"
)

var _ Bot = (*Weixin)(nil)

type Weixin struct {
	URL string
}

type weixinMarkdown struct {
	Content string `json:"content"`
}

type weixinPayload struct {
	ToUser   string         `json:"touser"`
	MsgType  string         `json:"msgtype"`
	Markdown weixinMarkdown `json:"markdown"`
}

func (w *Weixin) buildPayload(r record.Record) string {
	payload := weixinPayload{
		ToUser:  "@all",
		MsgType: "markdown",
		Markdown: weixinMarkdown{
			Content: "<font color=\"#e96900\" face=\"Fira Code\" size=3>New Connection</font>\n" +
				formatRecordField(r, `> **<font color="#e96900" face="Fira Code">%s: </font>**<font color="#e96900" face="Fira Code">%v</font>`),
		},
	}
	p, err := json.Marshal(&payload)
	if err != nil {
		return ""
	}
	return string(p)
}

func (w *Weixin) notice(r record.Record) error {
	resp, err := http.DefaultClient.Post(w.URL, "application/json", strings.NewReader(w.buildPayload(r)))
	if err != nil {
		return fmt.Errorf("HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read HTTP response body: %v", err)
		}
		return fmt.Errorf("non-success response status code %d with body: %s", resp.StatusCode, data)
	}
	return nil
}