package board

type SideType string

const (
	Field SideType = "field"
	Road  SideType = "road"
	City  SideType = "city"
	River SideType = "river"
)

type MeepleType string

const (
	NoUnit  MeepleType = ""
	Peasant MeepleType = "peasant"
	Priest  MeepleType = "priest"
)

type RegionID int

const NoRegion = 0

type PlayerID int

const NoOwner = -1
