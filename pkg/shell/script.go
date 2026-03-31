package shell

import "strings"

type Script struct {
	Shebang string
	Src     string
	Key     []byte
}

func (s *Script) Dump() string {
	var o strings.Builder
	o.Grow(len(s.Shebang) + len(s.Src) + len(s.Key) + 10)

	if s.Shebang != "" {
		o.WriteString(s.Shebang)
		o.WriteString("\n")
	}

	if s.Src != "" {
		o.WriteString(s.Src)
	}

	return o.String()
}
