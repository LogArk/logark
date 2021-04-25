package stdout

import "fmt"

func Send(log []byte) error {
	fmt.Println(string(log))
	return nil
}
