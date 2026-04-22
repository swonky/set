package base

func GetCap(cap ...int) int {
	size := 0
	if len(cap) > 0 && cap[0] > 0 {
		size = cap[0]
	}
	return size
}

func MaxInt() int {
	return int(^uint(0) >> 1)
}
