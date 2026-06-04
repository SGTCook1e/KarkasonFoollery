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

type Feature struct {
	Region   *Region
	Type     FeatureType `json:"type"`
	Unit     UnitType
	Sides    []Side `json:"sides"`
	Complete bool
}

func (f *Feature) HasSide(dir Direction) bool {
	for _, side := range f.Sides {
		if side.Direction == dir {
			return true
		}
	}
	return false
}

func (f *Feature) UpdateCompletion() {
	if f.Complete {
		return
	}
	for _, side := range f.Sides {
		if !side.Complete {
			return
		}
	}
	f.Complete = true
}
