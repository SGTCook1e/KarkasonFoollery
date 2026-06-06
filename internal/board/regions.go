package board

import (
	"fmt"
)

type Regions map[*Region]struct{}

func FindNeighbourRegion(neighbourTile *Tile, neighbourDir Direction) (*Region, bool) {
	feature := neighbourTile.FeatureByDirection(neighbourDir.Opposite())
	if feature.Region == nil {
		return nil, false
	} else {
		return feature.Region, true
	}
}

func (rs Regions) AppendRegion(newRegion *Region) {
	rs[newRegion] = struct{}{}
}

func ManageMonasteryRegions(b *Board, newTileCoord Coord) {

}

func (rs Regions) UniteRegions(neighbourRegions *[]*Region, newFeature *Feature) {
	unitedRegion := (*neighbourRegions)[0]
	nr := *neighbourRegions
	for i := 1; i <= len(nr); i++ {
		for j := 0; j <= len(nr[i].Districts); j++ {
			nr[i].Districts[j].Region = unitedRegion
		}
		unitedRegion.Districts = append(unitedRegion.Districts, nr[i].Districts...)
		delete(rs, nr[i])
	}
	newFeature.Region = unitedRegion
	unitedRegion.ExpandRegion(newFeature)
}

func (rs Regions) ManageRegions(b *Board, newTileCoord Coord) {
	newTile, exists := b.GetTile(newTileCoord)
	if !exists {
		panic(fmt.Errorf("NewTile not present!"))
	}

	defer rs.ManageCompletion(b, newTile, newTileCoord)

	if newTile.Monastery {
		rs.AppendRegion(MakeRegion(newTile, nil, 1, true))
	}
	// Map for each feature side of newTile with its neighbours regions slice
	neighbourRegions := make(map[*Feature][]*Region)

	for dir := Top; dir <= Left; dir++ {
		feature := newTile.FeatureByDirection(dir)
		if feature.Type == FeatureCity || feature.Type == FeatureRoad {
			neighbourTile, exists := b.GetTile(newTileCoord.CoordByDirection(dir))
			if exists {
				neighbourRegion, exists := FindNeighbourRegion(neighbourTile, dir)
				if exists {
					neighbourRegions[feature] = append(neighbourRegions[feature], neighbourRegion)
				} else {
					if _, exists := neighbourRegions[feature]; !exists {
						neighbourRegions[feature] = nil
					}
				}
			} else {
				if _, exists := neighbourRegions[feature]; !exists {
					neighbourRegions[feature] = nil
				}
			}
		}
	}

	for feature, regions := range neighbourRegions {
		switch len(regions) {
		case 0: //If 0 regions found for feature, make a new region
			rs.AppendRegion(MakeRegion(newTile, feature, 1, false))
		case 1: //If 1 regions found for feature, append this feature to existing neighbour region
			regions[0].ExpandRegion(feature)
		default: //If more than 1 regions found for feature, unite this regions and add feature to it
			rs.UniteRegions(&regions, feature)
		}
	}
}

func (rs Regions) ManageCompletion(b *Board, newTile *Tile, newTileCoord Coord) {

}
