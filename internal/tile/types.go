package tile

type SideType int

const (
	Field SideType = iota
	Road
	City
	River
)

type Direction int

const (
	Top Direction = iota
	Right
	Bottom
	Left
)
