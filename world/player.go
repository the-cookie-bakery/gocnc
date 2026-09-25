package world

type Team int

const (
	TeamNeutral Team = 0
	TeamGDI     Team = 1
	TeamNOD     Team = 2
)

type Player struct {
	ID       int
	Team     Team
	Credits  int
	PowerIn  int
	PowerOut int
}

func NewPlayer(id int, team Team) *Player {
	return &Player{
		ID:      id,
		Team:    team,
		Credits: 5000,
	}
}

func (p *Player) HasPower() bool {
	return p.PowerIn >= p.PowerOut || p.PowerOut == 0
}

func (p *Player) CanAfford(cost int) bool {
	return p.Credits >= cost
}

func (p *Player) Spend(amount int) {
	if p.Credits >= amount {
		p.Credits -= amount
	}
}

func (p *Player) AddCredits(amount int) {
	p.Credits += amount
}
