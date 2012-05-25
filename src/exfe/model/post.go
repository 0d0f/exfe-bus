package exfe_model

type Post struct {
	Id uint64
	By_identity Identity
	Content string
	Postable_id uint64
	Postable_type string
	Via string
	Relative map[string]string
}
