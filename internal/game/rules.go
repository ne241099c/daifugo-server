package game

func GetStrength(c *Card, isRev bool) int {
	if c.Suit == SuitJoker {
		return 99
	}

	strength := 0

	switch c.Rank {
	case 1: // Ace
		strength = 14
	case 2: // Two
		strength = 15
	default:
		strength = int(c.Rank)
	}

	if isRev {
		strength *= -1
	}

	return strength
}

func IsStronger(candidate, target *Card, isRev bool) bool {
	if target == nil {
		return true
	}

	myStrength := GetStrength(candidate, isRev)
	targetStrength := GetStrength(target, isRev)

	return myStrength > targetStrength
}
