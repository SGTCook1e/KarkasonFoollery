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

func (c Coord) CoordByDirection(dir Direction) Coord {
	var r Coord
	switch dir {
	case Top:
		r = c.Top()
	case Right:
		r = c.Right()
	case Bottom:
		r = c.Bottom()
	case Left:
		r = c.Left()
	}
	return r
}

func (c Coord) GetCoordsAround() []Coord {
	return []Coord{
		{X: c.X - 1, Y: c.Y - 1},
		c.Top(),
		{X: c.X + 1, Y: c.Y - 1},
		c.Left(),
		c.Right(),
		{X: c.X - 1, Y: c.Y + 1},
		c.Bottom(),
		{X: c.X + 1, Y: c.Y + 1},
	}
}
