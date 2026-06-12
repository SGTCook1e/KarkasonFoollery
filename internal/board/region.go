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
		r := Region{
			Type:      RegionMonastery,
			Districts: []*Feature{feature},
			Owner:     owner,
		}
		tile.MonasteryRegion = &r
		return &r
	}
	r := Region{
		Type:      RegionType(feature.Type),
		Districts: []*Feature{feature},
		Owner:     owner,
	}
	feature.Region = &r
	return &r
}

func GetNumberOfTilesAround(b *Board, newTileCoords Coord) int {
	var ctr int
	for _, c := range newTileCoords.GetCoordsAround() {
		_, exists := b.GetTile(c)
		if exists {
			ctr++
		}
	}
	return ctr
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

func (r *Region) CompleteRegion(rs Regions) int {
	if r.Type == RegionMonastery {
		delete(rs, r)
		return 8
	}
	score := len(r.Districts)
	delete(rs, r)
	return score
}
