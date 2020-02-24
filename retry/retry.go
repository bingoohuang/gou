package retry

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// RetryableFunc signature of retryable function
type RetryableFunc func() error

// Do ...
// nolint gomnd
func Do(fn RetryableFunc, opts ...Option) error {
	config := &Config{
		attempts:      10,
		delay:         100 * time.Millisecond,
		maxJitter:     100 * time.Millisecond,
		onRetry:       func(n int, err error) {},
		retryIf:       IsRecoverable,
		delayType:     CombineDelay(BackOffDelay, RandomDelay),
		lastErrorOnly: false,
	}

	for _, opt := range opts {
		opt(config)
	}

	errorsLen := config.attempts
	if config.lastErrorOnly {
		errorsLen = 1
	}

	errorLog := make(Error, errorsLen)
	lastErrIndex := 0

	for n := 0; n < config.attempts; n++ {
		if n > 0 {
			time.Sleep(config.delayType(n, config))
		}

		err := fn()
		if err == nil {
			return nil
		}

		if errorLog[lastErrIndex] = unpackUnrecoverable(err); !config.lastErrorOnly {
			lastErrIndex = n + 1
		}

		if !config.retryIf(err) {
			break
		}

		config.onRetry(n, err)
	}

	if config.lastErrorOnly {
		return errorLog[lastErrIndex]
	}

	return errorLog
}

// Error type represents list of errors in retry
type Error []error

// Error method return string representation of Error
// It is an implementation of error interface
func (e Error) Error() string {
	logWithNumber := make([]string, lenWithoutNil(e))

	for i, l := range e {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.Error()) // nolint gomnd
		}
	}

	return fmt.Sprintf("All attempts fail:\n%s", strings.Join(logWithNumber, "\n"))
}

func lenWithoutNil(e Error) (count int) {
	for _, v := range e {
		if v != nil {
			count++
		}
	}

	return
}

// WrappedErrors returns the list of errors that this Error is wrapping.
// It is an implementation of the `errwrap.Wrapper` interface
// in package [errwrap](https://github.com/hashicorp/errwrap) so that
// `retry.Error` can be used with that library.
func (e Error) WrappedErrors() []error {
	return e
}

type unrecoverableError struct {
	error
}

// Unrecoverable wraps an error in `unrecoverableError` struct
func Unrecoverable(err error) error {
	return unrecoverableError{err}
}

// IsRecoverable checks if error is an instance of `unrecoverableError`
func IsRecoverable(err error) bool {
	_, ok := err.(unrecoverableError)
	return !ok
}

func unpackUnrecoverable(err error) error {
	if ue, ok := err.(unrecoverableError); ok {
		return ue.error
	}

	return err
}

// IfFunc signature of retry if function
type IfFunc func(error) bool

// OnRetryFunc signature of OnRetry function
// n = count of attempts
type OnRetryFunc func(n int, err error)

// DelayTypeFunc ...
type DelayTypeFunc func(n int, config *Config) time.Duration

// Config config the retry options.
type Config struct {
	attempts      int
	delay         time.Duration
	maxJitter     time.Duration
	onRetry       OnRetryFunc
	retryIf       IfFunc
	delayType     DelayTypeFunc
	lastErrorOnly bool
}

// Option represents an option for retry.
type Option func(*Config)

// LastErrorOnly returns the direct last error that came from the retried function
// default is false (return wrapped errors with everything)
func LastErrorOnly(yes bool) Option { return func(c *Config) { c.lastErrorOnly = yes } }

// Attempts set count of retry, default is 10
func Attempts(attempts int) Option { return func(c *Config) { c.attempts = attempts } }

// Delay set delay between retry, default is 100ms
func Delay(delay time.Duration) Option { return func(c *Config) { c.delay = delay } }

// MaxJitter sets the maximum random Jitter between retries for RandomDelay
func MaxJitter(maxJitter time.Duration) Option { return func(c *Config) { c.maxJitter = maxJitter } }

// DelayType set type of the delay between retries
// default is BackOff
func DelayType(delayType DelayTypeFunc) Option { return func(c *Config) { c.delayType = delayType } }

// BackOffDelay is a DelayType which increases delay between consecutive retries
func BackOffDelay(n int, c *Config) time.Duration { return c.delay * (1 << n) }

// FixedDelay is a DelayType which keeps delay the same through all iterations
func FixedDelay(_ int, c *Config) time.Duration { return c.delay }

// RandomDelay is a DelayType which picks a random delay up to config.maxJitter
func RandomDelay(_ int, c *Config) time.Duration {
	return time.Duration(rand.Int63n(int64(c.maxJitter)))
}

// CombineDelay is a DelayType the combines all of the specified delays into a new DelayTypeFunc
func CombineDelay(delays ...DelayTypeFunc) DelayTypeFunc {
	return func(n int, c *Config) time.Duration {
		var total time.Duration
		for _, delay := range delays {
			total += delay(n, c)
		}

		return total
	}
}

// OnRetry function callback are called each retry
//
// log each retry example:
//
//	retry.Do(
//		func() error {
//			return errors.New("some error")
//		},
//		retry.OnRetry(func(n uint, err error) {
//			log.Printf("#%d: %s\n", n, err)
//		}),
//	)
func OnRetry(onRetry OnRetryFunc) Option { return func(c *Config) { c.onRetry = onRetry } }

// If controls whether a retry should be attempted after an error
// (assuming there are any retry attempts remaining)
//
// skip retry if special error example:
//
//	retry.Do(
//		func() error {
//			return errors.New("special error")
//		},
//		retry.If(func(err error) bool {
//			if err.Error() == "special error" {
//				return false
//			}
//			return true
//		})
//	)
//
// By default If stops execution if the error is wrapped using `retry.Unrecoverable`,
// so above example may also be shortened to:
//
//	retry.Do(
//		func() error {
//			return retry.Unrecoverable(errors.New("special error"))
//		}
//	)
func If(retryIf IfFunc) Option { return func(c *Config) { c.retryIf = retryIf } }
