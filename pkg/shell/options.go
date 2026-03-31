package shell

type Config struct {
	Shell     string
	B64Encode bool
	Encrypt   bool
	Stride    int
}

func EncodeSrc(v bool) func(s *ShellObfs) {
	return func(s *ShellObfs) {
		s.Config.B64Encode = v
	}
}

func EncryptSrc(v bool) func(s *ShellObfs) {
	return func(s *ShellObfs) {
		s.Config.Encrypt = v
	}
}

func SetShell(v string) func(s *ShellObfs) {
	return func(s *ShellObfs) {
		s.Config.Shell = v
	}
}

func SetStride(v int) func(s *ShellObfs) {
	return func(s *ShellObfs) {
		s.Config.Stride = v
	}
}
