package id

// AccountID defines account id object.
type AccountID string

func (a AccountID) String() string {
	return string(a)
}
