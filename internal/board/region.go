package board

type RegionType string

const (
	RegionCity      RegionType = "city"
	RegionRoad      RegionType = "road"
	RegionMonastery RegionType = "monastery"
)

const NoOwner = -1

type Region struct {
	Type         RegionType
	Districts    []*Feature
	Owner        int
	Complete     bool
	MonasteryCtr int
	// Score    int
}

func MakeRegion(tile *Tile, feature *Feature, owner int, isMonastery bool) *Region {
	if isMonastery {
		return &Region{
			Type:      RegionMonastery,
			Districts: []*Feature{feature},
			Owner:     owner,
		}
	}
	r := Region{
		Type:      RegionType(feature.Type),
		Districts: []*Feature{feature},
		Owner:     owner,
	}
	feature.Region = &r
	return &r
}

func (r *Region) ExpandRegion(newFeature *Feature) {
	r.Districts = append(r.Districts, newFeature)
	newFeature.Region = r
}

func (r *Region) UpdateCompletion() {
	if r.Complete {
		return
	}
	for _, feature := range r.Districts {
		if !feature.Complete {
			return
		}
	}
	r.Complete = true
}
