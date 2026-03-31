package shell

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Joey574/shadow/v2/internal/crypto"
)

type ShellObfs struct {
	Config *Config
}

func NewObfuscator(options ...func(s *ShellObfs)) *ShellObfs {
	s := &ShellObfs{
		Config: &Config{
			Shell:     "sh",
			B64Encode: false,
			Encrypt:   false,
			Stride:    5,
		},
	}

	for _, f := range options {
		f(s)
	}

	return s
}

func (s *ShellObfs) Obfuscate(src string) Script {
	code := ParseScript(src)
	code = StripShellScript(code)

	if s.Config.Encrypt {
		code = s.Encrypt(code)
	}

	if s.Config.B64Encode {
		code = s.B64Encode(code)
	}

	return code
}

func ParseScript(src string) Script {
	cut := strings.TrimLeft(src, " \n\t")
	env := ""

	if strings.HasPrefix(cut, "#!") {
		idx := strings.Index(cut, "\n")
		if idx == -1 {
			env = cut
		} else {
			env = cut[0:idx]
		}
	}

	return Script{
		Shebang: env,
		Src:     strings.Replace(src, env, "", 1),
		Key:     nil,
	}
}

// Strip comments and whitespace from shell script
// TODO : Should do proper minimization as well
func StripShellScript(src Script) Script {
	var code strings.Builder
	code.Grow(len(src.Src))

	lines := strings.Split(src.Src, "\n")
	for i := range lines {
		line := lines[i]

		// shell[0] = any code on that line
		split := strings.Split(line, "#")
		shell := strings.TrimSpace(split[0])
		if len(shell) == 0 || shell[0] == '#' {
			continue
		}

		// if last character in line is a '\' we can strip it and the newline
		if shell[len(shell)-1] == '\\' {
			shell = shell[:len(shell)-1]
			shell = strings.TrimSpace(shell)
			code.WriteString(shell)
			code.WriteByte(' ')
			continue
		}

		code.WriteString(shell)
		code.WriteString("\n")
	}

	return Script{
		Shebang: src.Shebang,
		Src:     strings.TrimSpace(code.String()),
		Key:     src.Key,
	}
}

func (s *ShellObfs) B64Encode(src Script) Script {
	shell := s.Config.Shell
	stride := s.Config.Stride

	b64 := base64.StdEncoding.EncodeToString([]byte(src.Src))
	cmd := fmt.Sprintf(`%s -c "$(printf %%s %s | base64 -d)" %s "$@"`, shell, b64, shell)

	size := (len(cmd) + stride - 1) / stride

	keys := make([]string, size)
	values := make([]string, 0, size)

	for i := range keys {
		keys[i] = varName(i)
	}

	// prevent predictable ordering of output variables
	crypto.Shuffle(keys)

	for i := range size {
		sidx := i * stride
		eidx := min(sidx+stride, len(cmd))
		values = append(values, cmd[sidx:eidx])
	}

	// randomize order in which variables are declared
	names := make([]string, len(keys))
	copy(names, keys)
	crypto.ShuffleWith(names, values)

	var compiled strings.Builder
	compiled.Grow(len(cmd) * 2)

	// use names here for randomized declaration order
	for i, n := range names {
		fmt.Fprintf(&compiled, "%s='%s';", n, values[i])
	}
	compiled.WriteString(`eval "`)

	// use keys here for consistent reassembling
	for _, k := range keys {
		fmt.Fprintf(&compiled, "$%s", k)
	}
	compiled.WriteString(`";`)

	return Script{
		Shebang: src.Shebang,
		Src:     compiled.String(),
		Key:     src.Key,
	}
}

func (s *ShellObfs) Encrypt(src Script) Script {
	shell := s.Config.Shell

	b64 := base64.StdEncoding.EncodeToString([]byte(src.Src))
	key := []byte(crypto.B64CHARSET)
	crypto.Shuffle(key)

	ct := crypto.RotCipher([]byte(b64), key, []byte(crypto.B64CHARSET))

	cmd := fmt.Sprintf(`stty -echo;read -p "" KEY;stty echo;%s -c "$(echo "%s" | tr $KEY "%s" | base64 -d)" %s "$@"`, shell, ct, crypto.B64CHARSET, shell)

	return Script{
		Shebang: src.Shebang,
		Src:     cmd,
		Key:     key,
	}
}
