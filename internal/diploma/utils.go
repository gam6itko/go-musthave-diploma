package diploma

func LuhnValidate(number uint64) bool {
	var sum byte
	for step := 1; number > 0; step++ {
		n := byte(number % 10)
		if step%2 == 0 && n != 9 {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		number = number / 10
		sum += n
	}

	return sum%10 == 0
}
