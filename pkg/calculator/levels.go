package calculator

import "math"

func CatacombsLevel(xp float64) float64 {
	return calcXPTable(xp, catacombsXPTable, 200_000_000)
}

func GardenLevel(xp float64) float64 {
	return calcXPTable(xp, gardenXPTable, 10000)
}

func SkillLevel(skill string, xp float64, extraLevels int) float64 {
	maxLevel := min(skillMaxLevel[skill]+extraLevels, len(skillXPTable))
	for level, levelXP := range skillXPTable[:maxLevel] {
		if xp < levelXP {
			return round(float64(level) + xp/levelXP)
		}

		xp -= levelXP
	}

	return float64(maxLevel)
}

func SkyBlockLevel(xp float64) float64 {
	return round(xp / 100)
}

func Average(levels ...float64) float64 {
	total := 0.0
	for _, level := range levels {
		total += level
	}

	return round(total / float64(len(levels)))
}

func SlayerLevel(slayer string, xp float64) float64 {
	table := slayerXPTable[slayer]
	previous := 0.0
	for level, levelXP := range table {
		if xp < levelXP {
			return round(float64(level) + (xp-previous)/(levelXP-previous))
		}

		previous = levelXP
	}

	return float64(len(table))
}

func calcXPTable(xp float64, table []float64, finalLevelXP float64) float64 {
	level := 0
	for _, levelXP := range table {
		if xp < levelXP {
			return round(float64(level) + xp/levelXP)
		}

		level++
		xp -= levelXP
	}

	return round(float64(level) + xp/finalLevelXP)
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}
