package entities

type StaticEntity struct {
	*BaseEntity
	BlocksMovement bool
}

func (s *StaticEntity) Type() string { return "object" }
func (s *StaticEntity) Team() int    { return 0 }
