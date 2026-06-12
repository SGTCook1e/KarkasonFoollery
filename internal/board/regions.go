package board

import (
	"fmt"
	"slices"
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

func (rs Regions) UniteRegions(neighbourRegions *[]*Region, newFeature *Feature) {
	unitedRegion := (*neighbourRegions)[0]
	nr := *neighbourRegions
	for i := 1; i < len(nr); i++ {
		for j := 0; j < len(nr[i].Districts); j++ {
			nr[i].Districts[j].Region = unitedRegion
		}
		unitedRegion.Districts = append(unitedRegion.Districts, nr[i].Districts...)
		delete(rs, nr[i])
	}
	newFeature.Region = unitedRegion
	unitedRegion.ExpandRegion(newFeature)
}

func (rs Regions) ManageRegions(b *Board, newTileCoord Coord, playerID int) {
	newTile, exists := b.GetTile(newTileCoord)
	if !exists {
		panic(fmt.Errorf("NewTile not present!"))
	}

	if newTile.Monastery {
		rs.AppendRegion(MakeRegion(newTile, nil, playerID, true))
	}
	// Map for each feature side of newTile with its neighbours regions slice
	neighbourRegions := make(map[*Feature][]*Region)

	for dir := Top; dir <= Left; dir++ {
		feature := newTile.FeatureByDirection(dir)
		if feature.Type != FeatureCity && feature.Type != FeatureRoad {
			continue
		}
		if _, ok := neighbourRegions[feature]; !ok {
			neighbourRegions[feature] = nil
		}
		neighbourTile, exists := b.GetTile(newTileCoord.CoordByDirection(dir))
		if exists {
			neighbourTile.CompleteFeaturesDirection(dir.Opposite())
		} else {
			continue
		}
		neighbourRegion, exists := FindNeighbourRegion(neighbourTile, dir)
		if !exists {
			continue
		}
		if !slices.Contains(neighbourRegions[feature], neighbourRegion) {
			neighbourRegions[feature] = append(neighbourRegions[feature], neighbourRegion)
		}
	}

	for feature, regions := range neighbourRegions {
		switch len(regions) {
		case 0: //If 0 regions found for feature, make a new region
			rs.AppendRegion(MakeRegion(newTile, feature, playerID, false))
		case 1: //If 1 regions found for feature, append this feature to existing neighbour region
			regions[0].ExpandRegion(feature)
		default: //If more than 1 regions found for feature, unite this regions and add feature to it
			rs.UniteRegions(&regions, feature)
		}
	}

	rs.ManageCompletion(b, neighbourRegions, newTileCoord, playerID)
}

func (rs Regions) ManageMonasteryRegions(b *Board, newTileCoord Coord) {
	for _, c := range newTileCoord.GetCoordsAround() {
		t, exists := b.GetTile(c)
		if !exists {
			continue
		}
		if !t.Monastery {
			continue
		}
		t.MonasteryRegion.MonasteryCtr++
		if t.MonasteryRegion.MonasteryCtr == 8 {
			fmt.Println(t.MonasteryRegion.CompleteRegion(rs))
		}
	}
}

func (rs Regions) ManageCompletion(b *Board, neighbourRegions map[*Feature][]*Region, newTileCoord Coord, playerID int) {
	rs.ManageMonasteryRegions(b, newTileCoord)

	for _, regions := range neighbourRegions {
		for _, r := range regions {
			for _, d := range r.Districts {
				d.UpdateCompletion()
			}
			r.UpdateCompletion()
			if r.Complete {
				fmt.Println(r.CompleteRegion(rs))
			}
		}
	}
}
