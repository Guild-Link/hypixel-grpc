package hypixel

import "github.com/guild-link/hypixel-grpc/pkg/common"

var catacombsXPTable = [...]float64{
	50, 75, 110, 160, 230, 330, 470, 670, 950, 1340,
	1890, 2665, 3760, 5260, 7380, 10300, 14400, 20000, 27600, 38000,
	52500, 71500, 97000, 132000, 180000, 243000, 328000, 445000, 600000, 800000,
	1065000, 1410000, 1900000, 2500000, 3300000, 4300000, 5600000, 7200000, 9200000, 12000000,
	15000000, 19000000, 24000000, 30000000, 38000000, 48000000, 60000000, 75000000, 93000000, 116250000,
}

var farmingXPTable = [...]float64{
	50, 125, 200, 300, 500, 750, 1000, 1500, 2000, 3500,
	5000, 7500, 10000, 15000, 20000, 30000, 50000, 75000, 100000, 200000,
	300000, 400000, 500000, 600000, 700000, 800000, 900000, 1000000, 1100000, 1200000,
	1300000, 1400000, 1500000, 1600000, 1700000, 1800000, 1900000, 2000000, 2100000, 2200000,
	2300000, 2400000, 2500000, 2600000, 2750000, 2900000, 3100000, 3400000, 3700000, 4000000,
	4300000, 4600000, 4900000, 5200000, 5500000, 5800000, 6100000, 6400000, 6700000, 7000000,
}

var gardenXPTable = [...]float64{
	0, 70, 70, 140, 240, 600, 1500, 2000, 2500, 3000, 10000, 10000, 10000, 10000, 10000,
}

func calcXPTable(xp float64, table []float64, finalLevelXP ...float64) float64 {
	level := 0
	for _, levelXP := range table {
		if xp < levelXP {
			return common.RoundToTwo(float64(level) + xp/levelXP)
		}

		level++
		xp -= levelXP
	}

	overflow := 0.0
	if len(finalLevelXP) > 0 {
		overflow = finalLevelXP[0]
	}

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
