package game

import b "KarkasonFoollery/internal/board"

type Player struct {
	Id      b.PlayerID
	Score   int
	Meeples map[b.MeepleType]int
}

func NewPlayer(id b.PlayerID) Player {
	p := Player{
		Id:      id,
		Score:   0,
		Meeples: make(map[b.MeepleType]int),
	}
	p.Meeples[b.Peasant] = 7
	p.Meeples[b.Priest] = 1

	return p
}

func (p *Player) Clone() *Player {
	clone := Player{
		Id:      p.Id,
		Score:   p.Score,
		Meeples: make(map[b.MeepleType]int),
	}
	for k, v := range p.Meeples {
		clone.Meeples[k] = v
	}
	return &clone
}

func (p *Player) canPlaceMeeple(mType b.MeepleType) bool {
	return p.Meeples[mType] > 0
}

func (p *Player) TakeMeeple(mType b.MeepleType) bool {
	if !p.canPlaceMeeple(mType) {
		return false
	}
	p.Meeples[mType]--
	return true
}
