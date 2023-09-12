package dto

//go:generate dbgen -type Actor

// Actor TBD
type Actor struct {
	ID         int64  `db:"id"`
	RootID     *int64 `db:"root_id"`
	Name       string `db:"name"`
	NormalName string `db:"normal_name"`
}

// ActorList TBD
type ActorList []*Actor

// IDs ids slice
func (a ActorList) IDs() []int64 {
	o := make([]int64, len(a))
	for i := range a {
		o[i] = a[i].ID
	}
	return o
}

// NormalNames normalNames slice
func (a ActorList) NormalNames() []string {
	o := make([]string, len(a))
	for i := range a {
		o[i] = a[i].NormalName
	}
	return o
}

// IDToActor TBD
func (a ActorList) IDToActor() map[int64]*Actor {
	id2actor := make(map[int64]*Actor, 0)
	for _, actor := range a {
		id2actor[actor.ID] = actor
	}
	return id2actor
}
