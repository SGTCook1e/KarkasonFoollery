package game

import (
	b "KarkasonFoollery/internal/board"
)

type Regions struct {
	nextID b.RegionID
	byID   map[b.RegionID]*Region
}

func NewRegions() *Regions {
	return &Regions{
		nextID: 1,
		byID:   make(map[b.RegionID]*Region),
	}
}

func FindNeighbourRegionID(neighbourTile b.Tile, neighbourDir b.Direction) (b.RegionID, bool) {
	feature, _ := neighbourTile.FeatureByDirection(neighbourDir.Opposite())
	if feature.RegionID == b.NoRegion {
		return b.NoRegion, false
	} else {
		return feature.RegionID, true
	}
}

func (rs *Regions) AppendRegion(r Region) b.RegionID {
	if r.ID != b.NoRegion {
		panic("region already has ID")
	}

	id := rs.nextID
	rs.nextID++

	r.ID = id
	rs.byID[id] = &r

	return id
}

func (rs *Regions) UniteRegions(coords b.Coord, newFeature int, RegionIds []b.RegionID) {
	targetId := RegionIds[0]
	for i := 1; i < len(RegionIds); i++ {
		rs.byID[targetId].Districts = append(rs.byID[targetId].Districts, rs.byID[RegionIds[i]].Districts...)
		rs.DeleteRegion(RegionIds[i])
	}
	rs.byID[targetId].ExpandRegion(coords, newFeature)
}

func (rs *Regions) DeleteRegion(id b.RegionID) {
	delete(rs.byID, id)
}
