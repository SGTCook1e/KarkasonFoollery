package game

import (
	b "KarkasonFoollery/internal/board"
)

type Regions struct {
	nextID b.RegionID
	ByID   map[b.RegionID]*Region
}

func NewRegions() *Regions {
	return &Regions{
		nextID: 1,
		ByID:   make(map[b.RegionID]*Region),
	}
}

func (rs *Regions) Clone() Regions {
	clone := Regions{
		nextID: rs.nextID,
		ByID:   make(map[b.RegionID]*Region, len(rs.ByID)),
	}

	for id, region := range rs.ByID {
		if region != nil {
			regionClone := region.Clone()
			clone.ByID[id] = &regionClone
		}
	}

	return clone
}

func findNeighbourRegionID(neighbourTile b.Tile, neighbourDir b.Direction) (b.RegionID, bool) {
	feature, _ := neighbourTile.FeatureByDirection(neighbourDir.Opposite())
	if feature.RegionID == b.NoRegion {
		return b.NoRegion, false
	} else {
		return feature.RegionID, true
	}
}

func (rs *Regions) addRegion(r Region) b.RegionID {
	if r.ID != b.NoRegion {
		panic("region already has ID")
	}

	id := rs.nextID
	rs.nextID++

	r.ID = id
	rs.ByID[id] = &r

	return id
}

func (rs *Regions) mergeRegions(coords b.Coord, newFeature int, regionIds []b.RegionID) {
	targetReg := rs.ByID[regionIds[0]]

	for i := 1; i < len(regionIds); i++ {
		reg := rs.ByID[regionIds[i]]
		targetReg.Districts = append(targetReg.Districts, reg.Districts...)
		rs.deleteRegion(regionIds[i])
	}

	targetReg.expandRegion(coords, newFeature)
}

func (rs *Regions) deleteRegion(id b.RegionID) {
	delete(rs.ByID, id)
}
