package nextdata

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var re = regexp.MustCompile(`<script id="__NEXT_DATA__"[^>]*>(.*?)</script>`)

func Extract(html []byte) (json.RawMessage, error) {
	m := re.FindSubmatch(html)
	if m == nil {
		return nil, fmt.Errorf("__NEXT_DATA__ not found")
	}
	return json.RawMessage(m[1]), nil
}

func Unmarshal(html []byte, v any) error {
	raw, err := Extract(html)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}
