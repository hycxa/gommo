package ext

func Assert(condition bool) {
	if !condition {
		panic("assert failed")
	}
}

func AssertM(condition bool, m string) {
	if !condition {
		panic(m)
	}
}

func AssertE(err error) {
	if err != nil {
		panic(err.Error())
	}
}
