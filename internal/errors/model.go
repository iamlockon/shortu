package errors

type Error struct {
	Code int
	Msg  string
}

const (
	// Config error
	InvalidConfigError = iota + 1000

	// Database error
	FailedToGetDBConnError = iota + 10000
)
