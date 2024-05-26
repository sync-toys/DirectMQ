package dmqspecagent

import (
	"fmt"
	"strconv"
	"time"
)

func retry[TResult interface{}](maxRetries int, retryInterval time.Duration, funcToExec func() TResult) TResult {
	panics := []interface{}{}

	for i := 0; i < maxRetries; i++ {
		ok, result := func() (ok bool, result TResult) {
			defer func() {
				if r := recover(); r != nil {
					panics = append(panics, r)
					ok = false
				}
			}()

			return true, funcToExec()
		}()

		if ok {
			return result
		}

		time.Sleep(retryInterval)
	}

	panic(fmt.Sprintf("failed to execute function after %s retries: %v"+strconv.Itoa(maxRetries), panics))
}
