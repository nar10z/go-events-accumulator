/*
 * Copyright (c) 2023.
 *
 * License MIT (https://raw.githubusercontent.com/nar10z/go-accumulator/main/LICENSE)
 *
 * Developed thanks to Nikita Terentyev (nar10z). Use it for good, and let your code work without problems!
 */

package go_accumulator

import "sync/atomic"

type eventExtended[T comparable] struct {
	done atomic.Bool
	// return error of flush operation
	fallback chan<- error
	// original data
	e T
}
