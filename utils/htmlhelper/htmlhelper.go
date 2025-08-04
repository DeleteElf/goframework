package htmlhelper

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"io"
	"os"
)

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}

// MarkdownToHTML 将Markdown格式的内容转化成html格式的内容
func MarkdownToHTML(md string) string {
	bytes := blackfriday.Run([]byte(md))
	theHTML := string(bytes)
	return bluemonday.UGCPolicy().Sanitize(theHTML)
}
