package common

// HandleError handles the error
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
