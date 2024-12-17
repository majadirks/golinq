package gl

import (
	"cmp"
)

// For each element in a channel, apply the given map function
// and send the result on a new channel.
// That new channel is returned.
func Map[T1 any, T2 any](source chan T1, mapper func(T1) T2) chan T2 {
	if source == nil {
		return nil
	}
	output := make(chan T2)
	go func() {
		for s := range source {
			output <- mapper(s)
		}
		close(output)
	}()
	return output
}

// Applies the given mapper to elements from the two channels until one of the channels is closed
func Zip[T1 any, T2 any, T3 any](xs chan T1, ys chan T2, mapper func(T1, T2) T3) chan T3 {
	if xs == nil || ys == nil {
		return nil
	}
	output := make(chan T3)
	go func() {
		for {
			x, hasX := <-xs
			y, hasY := <-ys
			if !hasX || !hasY {
				break
			}
			output <- mapper(x, y)
		}
		close(output)
	}()
	return output
}

// For each element in a channel,
// apply the given predicate and send any results
// where the predicate returns true on a new channel.
// That new channel is returned.
func Filter[T any](source chan T, predicate func(T) bool) chan T {
	if source == nil {
		return nil
	}
	output := make(chan T)
	go func() {
		for s := range source {
			if predicate(s) {
				output <- s
			}
		}
		close(output)
	}()
	return output
}

// Receives the first n = count values from a channel and sends them on a new channel.
// If the channel closes before n values are sent, all those values are sent.
func Take[T any](source chan T, count int) chan T {
	output := make(chan T)
	taken := 0
	go func() chan T {
		for s := range source {
			taken++
			if taken > count {
				break
			}
			output <- s
		}
		close(output)
		return output
	}()
	return output
}

// Ignores the first n = count values from a channel
// and sends the rest (if any) on a new channel.
func Skip[T any](source chan T, count int) chan T {
	output := make(chan T)
	skipped := 0
	go func() {
		for s := range source {
			skipped++
			if skipped > count {
				output <- s
			}
		}
		close(output)
	}()
	return output
}

// Aggregation functions

// Returns the maximum element received on the given channel
func Max[T cmp.Ordered](source chan T) T {
	var max T
	first := true
	for s := range source {
		if first || s > max {
			max = s
		}
		first = false
	}
	return max
}

// Returns the first element received on the given channel
func First[T any](source chan T) T {
	first := <-source
	return first
}

// Returns the Last element received on the given channel
func Last[T any](source chan T) T {
	var last T
	for s := range source {
		last = s
	}
	return last
}

// Listens on a channel until the channel is closed,
// and return the number of elements received
func Count[T any](source chan T) int {
	count := 0
	for {
		_, more := <-source
		if !more {
			return count
		}
		count++
	}
}

// Given a channel of numeric values, return their Sum
func Sum[T float32 | float64 | int | int32 | int64](source chan T) T {
	var ret T
	ret = 0
	for s := range source {
		ret += s
	}
	return ret
}

// Create a channel and send each element
// of the given array on that channel.
// After closing the channel, return it
func From[T any](source []T) chan T {
	output := make(chan T)
	go func() {
		for _, elem := range source {
			output <- elem
		}
		close(output)
	}()
	return output
}

// Output all the Fibonacci numbers onto a channel
func Fibonaccis() chan int {
	output := make(chan int)
	a := 1
	b := 1
	go func() {
		output <- a
		output <- b
		for {
			c := a + b
			output <- c
			a = b
			b = c
		}
	}()
	return output
}
