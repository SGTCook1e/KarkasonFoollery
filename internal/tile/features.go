package tile

type FeatureType int

const (
	FeatureField FeatureType = iota
	FeatureRoad
	FeatureCity
	FeatureRiver
	FeatureMonastery
)

type Feature struct {
	ID   int
	Type FeatureType

	Sides []Direction
}

func (f *Feature) HasSide(dir Direction) bool {
	for _, side := range f.Sides {
		if side == dir {
			return true
		}
	}
	return false
}
