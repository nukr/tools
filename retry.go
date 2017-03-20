package tools

import (
	"fmt"
	"time"
)

// Retry ...
// Shameless copy from https://blog.abourget.net/en/2016/01/04/my-favorite-golang-retry-function/
func Retry(attempts int, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return nil
		}

		if i >= (attempts - 1) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
