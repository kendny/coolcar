package id

// AccountID defines account id object
// Identifier Type设计模式， account_id 类型具有唯一确定性
type AccountID string

func (a AccountID) String() string {
	return string(a)
}

// TripID defines trip id object
type TripID string

func (t TripID) String() string {
	return string(t)
}

// IdentityID  defines identity id object
type IdentityID string

func (i IdentityID) String() string {
	return string(i)
}

// CardID  defines car object
type CardID string

func (c CardID) String() string {
	return string(c)
}

// BlobID  defines car object
type BlobID string

func (c BlobID) String() string {
	return string(c)
}
