package id

// AccountID defines account id object
// Identifier Type设计模式， account_id 类型具有唯一确定性
type AccountID string

func (a AccountID) String() string {
	return string(a)
}

type TripID string

func (t TripID) String() string {
	return string(t)
}
