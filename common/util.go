package common

func checkErrorNotNil(err error) bool {
	if err != nil {
		return true
	}
	return false
}
