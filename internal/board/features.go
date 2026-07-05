package board

type FeatureType string

const (
	FeatureCity      FeatureType = "city"
	FeatureRoad      FeatureType = "road"
	FeatureMonastery FeatureType = "monastery"
	FeatureRiver     FeatureType = "river"
	FeatureField     FeatureType = "field"
)

type Side struct {
	Direction Direction `json:"dir"`
	Complete  bool
}

type MeepleSlot struct {
	Type  MeepleType
	Owner PlayerID
}

type Feature struct {
	RegionID RegionID
	Type     FeatureType `json:"type"`
	Shield   bool        `json:"shield"`
	Meeple   MeepleSlot
	Sides    []Side `json:"sides"`
}

func (f *Feature) HasSide(dir Direction) bool {
	for _, side := range f.Sides {
		if side.Direction == dir {
			return true
		}
	}
	return false
}

func (f *Feature) IsOtherSidesComplete(dir Direction) bool {
	for _, s := range f.Sides {
		if s.Direction == dir {
			continue
		}
		if !s.Complete {
			return false
		}
	}
	return true
}
