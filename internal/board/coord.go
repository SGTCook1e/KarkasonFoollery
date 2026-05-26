package board

type Coord struct {
	X int
	Y int
}

func (c Coord) Top() Coord {
	return Coord{X: c.X, Y: c.Y - 1}
}

func (c Coord) Right() Coord {
	return Coord{X: c.X + 1, Y: c.Y}
}

func (c Coord) Bottom() Coord {
	return Coord{X: c.X, Y: c.Y + 1}
}

func (c Coord) Left() Coord {
	return Coord{X: c.X - 1, Y: c.Y}
}
