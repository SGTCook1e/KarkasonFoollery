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

func MakeRegion(tile *Tile, feature *Feature, owner int) *Region {
	r := &Region{
		Districts: []*Feature{feature},
		Owner:     owner,
		Complete:  false,
	}
	if tile.Monastery {
		r.Type = RegionMonastery
	} else {
		r.Type = RegionType(feature.Type)
		feature.Region = r
	}
	return r
}

func (r *Region) ExpandRegion(newFeature *Feature) {
	newFeature.Region = r
	r.Districts = append(r.Districts, newFeature)
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
