package hypixel

func CatacombsLevel(xp float64) float64 {
	return calcXPTable(xp, catacombsXPTable[:], 200_000_000)
}

func FarmingLevel(xp float64) float64 {
	return calcXPTable(xp, farmingXPTable[:])
}

func GardenLevel(xp float64) float64 {
	return calcXPTable(xp, gardenXPTable[:], 10000)
}
