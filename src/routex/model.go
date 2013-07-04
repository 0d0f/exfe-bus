package routex

type Location struct {
	Timestamp int64  `json:"timestamp"`
	Lng       string `json:"lng"`
	Lat       string `json:"lat"`
}

type CrossToken struct {
	TokenType  string `json:"token_type"`
	CrossId    uint64 `json:"cross_id"`
	IdentityId uint64 `json:"identity_id"`
	UserId     int64  `json:"user_id"`
	CreatedAt  int64  `json:"created_time"`
	UpdatedAt  int64  `json:"updated_time"`
}

type LocationRepo interface {
	Save(id string, crossId uint64, l Location) error
	Load(id string, crossId uint64) ([]Location, error)
}

type RouteRepo interface {
	Save(crossId uint64, content string) error
	Load(crossId uint64) (string, error)
}
