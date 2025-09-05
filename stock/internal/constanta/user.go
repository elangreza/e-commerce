package constanta

//go:generate stringer -type=Key
type Key string

const (
	UserIDKey Key = "user_id"
)
