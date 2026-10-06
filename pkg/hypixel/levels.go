package hypixel

import "github.com/guild-link/hypixel-grpc/pkg/common"

func CatacombsLevel(xp float64) float64 {
	return calcXPTable(xp, catacombsXPTable[:], 200_000_000)
}

func FarmingLevel(xp float64) float64 {
	return calcXPTable(xp, farmingXPTable[:], 0)
}

func GardenLevel(xp float64) float64 {
	return calcXPTable(xp, gardenXPTable[:], 10000)
}

func calcXPTable(xp float64, table []float64, finalLevelXP float64) float64 {
	level := 0
	for _, levelXP := range table {
		if xp < levelXP {
			return common.RoundToTwo(float64(level) + xp/levelXP)
		}

		level++
		xp -= levelXP
	}

	overflow := finalLevelXP

	if overflow == 0 {
		slope := 600_000.0
		overflow = slope
		if len(table) > 0 {
			overflow += table[len(table)-1]
		}

		for xp >= overflow {
			level++
			xp -= overflow
			overflow += slope
			if level%10 == 0 {
				slope *= 2
			}
		}
	}

	return common.RoundToTwo(float64(level) + xp/overflow)
}
