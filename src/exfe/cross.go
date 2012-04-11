package exfe

type Cross struct {
	Id uint64
	Type string
	Title string
	Description string
	Time CrossTime
	Place Place
	Attribute map[string]string
	Exfee []Invitation
	Widget []interface{}
	Relative []struct {
		Id uint64
		Relation string
	}
	By_identity Identity
}
