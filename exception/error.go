package exception

import "fmt"

func PanicLogging(err interface{}) {
	fmt.Printf("%+v\n", err)
	if err != nil {
		panic(err)
	}
}
