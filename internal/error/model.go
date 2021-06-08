package error

type Error struct {
	Code int
	Msg  string
}

const (
	InvalidConfigError = iota + 1000
)
