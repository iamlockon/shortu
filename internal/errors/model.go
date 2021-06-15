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
	QueryRowError
	QueryError
	ExecSQLError
	BeginError
	CommitError
	ZeroAffectedSQLError
	URLNotFoundError
	ScanError
	RowsError
	// Cache error
	CacheSetTextFailedError = iota + 20000
)
