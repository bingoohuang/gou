// nolint gomnd
package retry

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestCustomRetryFunction(t *testing.T) {
	url := "http://example.com"

	var body []byte

	err := Do(
		func() error {
			resp, err := http.Get(url)

			if err == nil {
				defer func() {
					if err := resp.Body.Close(); err != nil {
						panic(err)
					}
				}()
				body, err = ioutil.ReadAll(resp.Body)
			}

			return err
		},
		DelayType(func(n int, config *Config) time.Duration { return 0 }),
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}

func TestGet(t *testing.T) {
	url := "http://example.com"

	var body []byte

	err := Do(
		func() error {
			resp, err := http.Get(url)

			if err == nil {
				defer func() {
					if err := resp.Body.Close(); err != nil {
						panic(err)
					}
				}()
				body, err = ioutil.ReadAll(resp.Body)
			}

			return err
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}

func TestDoAllFailed(t *testing.T) {
	var retrySum int
	err := Do(
		func() error { return errors.New("test") },
		OnRetry(func(n int, err error) { retrySum += n }),
		Delay(time.Nanosecond),
	)
	assert.Error(t, err)

	expectedErrorFormat := `All attempts fail:
#1: test
#2: test
#3: test
#4: test
#5: test
#6: test
#7: test
#8: test
#9: test
#10: test`
	assert.Equal(t, expectedErrorFormat, err.Error(), "retry error format")
	assert.Equal(t, 45, retrySum, "right count of retry")
}

func TestDoFirstOk(t *testing.T) {
	var retrySum int
	err := Do(
		func() error { return nil },
		OnRetry(func(n int, err error) { retrySum += n }),
	)
	assert.NoError(t, err)
	assert.Equal(t, 0, retrySum, "no retry")

}

func TestRetryIf(t *testing.T) {
	var retryCount int
	err := Do(
		func() error {
			if retryCount >= 2 {
				return errors.New("special")
			}

			return errors.New("test")
		},
		OnRetry(func(n int, err error) { retryCount++ }),
		If(func(err error) bool {
			return err.Error() != "special"
		}),
		Delay(time.Nanosecond),
	)
	assert.Error(t, err)

	expectedErrorFormat := `All attempts fail:
#1: test
#2: test
#3: special`
	assert.Equal(t, expectedErrorFormat, err.Error(), "retry error format")
	assert.Equal(t, 2, retryCount, "right count of retry")

}

func TestDefaultSleep(t *testing.T) {
	start := time.Now()
	err := Do(
		func() error { return errors.New("test") },
		Attempts(3),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur > 300*time.Millisecond, "3 times default retry is longer then 300ms")
}

func TestFixedSleep(t *testing.T) {
	start := time.Now()
	err := Do(
		func() error { return errors.New("test") },
		Attempts(3),
		DelayType(FixedDelay),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur < 500*time.Millisecond, "3 times default retry is shorter then 500ms")
}

func TestLastErrorOnly(t *testing.T) {
	var retrySum int
	err := Do(
		func() error { return fmt.Errorf("%d", retrySum) },
		OnRetry(func(n int, err error) { retrySum++ }),
		Delay(time.Nanosecond),
		LastErrorOnly(true),
	)
	assert.Error(t, err)
	assert.Equal(t, "9", err.Error())
}

func TestUnrecoverableError(t *testing.T) {
	attempts := 0
	expectedErr := errors.New("error")
	err := Do(
		func() error {
			attempts++
			return Unrecoverable(expectedErr)
		},
		Attempts(2),
		LastErrorOnly(true),
	)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 1, attempts, "unrecoverable error broke the loop")
}

func TestCombineFixedDelays(t *testing.T) {
	start := time.Now()
	err := Do(
		func() error { return errors.New("test") },
		Attempts(3),
		DelayType(CombineDelay(FixedDelay, FixedDelay)),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur > 400*time.Millisecond, "3 times combined, fixed retry is longer then 400ms")
	assert.True(t, dur < 500*time.Millisecond, "3 times combined, fixed retry is shorter then 500ms")
}

func TestRandomDelay(t *testing.T) {
	start := time.Now()
	err := Do(
		func() error { return errors.New("test") },
		Attempts(3),
		DelayType(RandomDelay),
		MaxJitter(50*time.Millisecond),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur > 2*time.Millisecond, "3 times random retry is longer then 2ms")
	assert.True(t, dur < 100*time.Millisecond, "3 times random retry is shorter then 100ms")
}
