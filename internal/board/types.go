package board

type SideType string

const (
	Field SideType = "field"
	Road  SideType = "road"
	City  SideType = "city"
	River SideType = "river"
)

type UnitType string

const (
	NoUnit  UnitType = ""
	Peasant UnitType = "peasant"
	Priest  UnitType = "priest"
)

type RegionID int

const NoRegion = 0

type OwnerID int

const NoOwner = -1
