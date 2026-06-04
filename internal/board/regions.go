package board

import (
	"fmt"
)

type Regions map[*Region]struct{}

func FindNeighbourRegion(neighbourTile *Tile, neighbourDir Direction) *Region {
	feature := neighbourTile.FeatureByDirection(neighbourDir.Opposite())
	return feature.Region
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
		panic(fmt.Errorf("Tile not present!"))
	}

	defer rs.ManageCompletion(b, newTile, newTileCoord)

	if newTile.Monastery {
		rs.AppendRegion(MakeRegion(newTile, nil, 1))
		return
	}
	// Map for each feature side of newTile with its neighbours regions slice
	neighbourRegions := make(map[*Feature][]*Region)

	for dir := Top; dir <= Left; dir++ {
		neighbourTile, exists := b.GetTile(newTileCoord.CoordByDirection(dir))
		if exists {
			feature := newTile.FeatureByDirection(dir)
			if feature.Type == FeatureCity || feature.Type == FeatureRoad {
				region := FindNeighbourRegion(neighbourTile, dir)
				neighbourRegions[feature] = append(neighbourRegions[feature], region)
			}
		}
	}

	for feature, regions := range neighbourRegions {
		switch len(regions) {
		case 0:
			rs.AppendRegion(MakeRegion(newTile, feature, 1))
		case 1:
			regions[0].ExpandRegion(feature)
		default:
			rs.UniteRegions(&regions, feature)
		}
	}
}

func (rs Regions) ManageCompletion(b *Board, newTile *Tile, newTileCoord Coord) {

}
