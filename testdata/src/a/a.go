package a

func f1() (int, int) {
	return 0, 0
}

func f2(x, y int) int {
	if x == 2 {
		if y == 2 {
			if x+y == 4 {
				return x + y
			}
		}
	}
	return 0
}

func f3(x, y, z int) {
	for i := 0; i < x; i++ {
		for j := 0; j < x; j++ {
			for k := 0; k < x; k++ {
				if i == 0 && j == 0 && k == 0 {
					println(i + j + k)
				} else {
					println(i, j, k)
				}
			}
		}
	}
}
