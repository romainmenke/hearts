package fakedb

func svg(hearts int) string {
	switch hearts {
	case 1:
		return oneHeart()
	case 2:
		return twoHearts()
	case 3:
		return threeHearts()
	default:
		return zeroHearts()
	}
}
