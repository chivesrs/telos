package droprate

import (
	"fmt"
	"math"
)

// DropRate calculates the expected drop rate with the given enrage and streak.
// Also takes into account if Luck of the Dwarves is worn (lotd).
// Output is the denominator of the ratio (if output is 10, then drop rate is 1/10).
// See http://runescape.wikia.com/wiki/Telos,_the_Warden#Unique for calculation.
func DropRate(enrage int64, streak int64, lotd bool) (int64, error) {
	if streak <= 0 {
		return 0, fmt.Errorf("streak %v must be positive (>0)", streak)
	}
	if enrage < 0 {
		return 0, fmt.Errorf("enrage %v must be non-negative (>=0)", enrage)
	}
	if enrage > 4000 {
		return 0, fmt.Errorf("enrage %v cannot exceed 4000 (<=4000)", enrage)
	}

	// "Killing Telos within 25-99% enrage will grant the player Silver tier loot,
	// which has a 10 times reduced chance to access the unique drop table. If the
	// enrage is below 25%, the player will be awarded Bronze tier loot, which has
	// a further 3 times reduced chance to access the unique drop table, meaning that
	// the unique drops would be 30 times rarer than normally."
	divisor := 1
	if enrage >= 25 && enrage < 100 {
		divisor = 10
	} else if enrage < 25 {
		divisor = 30
	}

	// Luck of the dwarves adds on 25% to enrage for drop rate calculation (after determining loot tier).
	modifiedEnrage := enrage
	if lotd {
		modifiedEnrage += 25
	}

	equation := 10 + 0.25*float64(modifiedEnrage) + 3*float64(streak)
	equation = equation / float64(divisor)
	return round(10000 / equation), nil
}

func round(f float64) int64 {
	return int64(f + math.Copysign(0.5, f))
}
