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

func (rs *Regions) Clone() Regions {
	clone := Regions{
		nextID: rs.nextID,
		byID:   make(map[b.RegionID]*Region, len(rs.byID)),
	}

	for id, region := range rs.byID {
		if region != nil {
			regionClone := region.Clone()
			clone.byID[id] = &regionClone
		}
	}

	return clone
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

func (rs *Regions) MergeRegions(coords b.Coord, newFeature int, regionIds []b.RegionID) {
	targetId := regionIds[0]
	for i := 1; i < len(regionIds); i++ {
		rs.byID[targetId].Districts = append(rs.byID[targetId].Districts, rs.byID[regionIds[i]].Districts...)
		rs.DeleteRegion(regionIds[i])
	}
	rs.byID[targetId].ExpandRegion(coords, newFeature)
}

func (rs *Regions) DeleteRegion(id b.RegionID) {
	delete(rs.byID, id)
}
