package gl

import (
	"cmp"
	"context"
	"fmt"
	"strings"
	"time"
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

// Ignores the first n = count vales from a channel
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

// Listens on a channel until the channel is closed or a timeout threshold is reached
// and return the number of elements received, or -1 if timed out.
func countOrTimeOut[T any](source chan T, timeoutSec int) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()
	ch := make(chan int)
	go func() { ch <- Count(source) }()
	select {
	case res := <-ch:
		return res // counted successfully
	case <-ctx.Done():
		return -1 // timed out
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

// Given a channel of integers,
// return a string holding those integers, separated
// by the given separator
func concatInts(separator string, source chan int) string {
	if source == nil {
		return ""
	}
	var builder strings.Builder
	first := true
	for elem := range source {
		if !first {
			builder.WriteString(separator)
		}

		str := fmt.Sprintf("%d", elem)
		builder.WriteString(str)

		first = false
	}
	return builder.String()
}

// Given a channel of float64s,
// return a string holding those float64s, separated
// by the given separator
func concatFloats(separator string, source chan float64) string {
	if source == nil {
		return ""
	}
	var builder strings.Builder
	first := true
	for elem := range source {
		if !first {
			builder.WriteString(separator)
		}

		str := fmt.Sprintf("%.6f", elem)
		builder.WriteString(str)

		first = false
	}
	return builder.String()
}

func main() {
	ints := []int{1, 2, 3, 6, 4, 1, 9, 5, 8}

	fmt.Println("Given ints:")
	fmt.Println(concatInts(", ", From(ints))) // prints "1, 2, 3, 6, 4, 1, 9, 5, 8"

	isEven := func(i int) bool { return i%2 == 0 }
	square := func(i int) int { return i * i }

	fmt.Println("Even squares of given ints:")
	squaresOfEvens := concatInts(", ", Filter(Map(From(ints), square), isEven))
	fmt.Println(squaresOfEvens) // prints "4, 36, 16, 64"

	fmt.Println("Max of given ints:")
	max := Max(From(ints))
	fmt.Println(max) // prints "9"

	fmt.Println("First int:")
	first := First(From(ints))
	fmt.Println(first) // prints "1"

	fmt.Println("Last int:")
	last := Last(From(ints))
	fmt.Println(last) // prints "8"

	fmt.Println("Count of ints:")
	count := Count(From(ints))
	fmt.Println(count) // prints "9"

	fmt.Println("Sum of ints:")
	sum := Sum(From(ints))
	fmt.Println(sum) // prints "39"

	fmt.Println("Sum of first three ints:")
	sum3 := Sum(Take(From(ints), 3))
	fmt.Println(sum3) // prints "6"

	fmt.Println("Sum of final two ints:")
	sumFinal2 := Sum(Skip(From(ints), count-2))
	fmt.Println(sumFinal2) // prints "13"

	fmt.Println("First ten Fibonacci numbers")
	first10Fibs := Take(Fibonaccis(), 10)
	fmt.Println(concatInts(", ", first10Fibs)) // prints "1, 1, 2, 3, 5, 8, 13, 21, 34, 55"

	firstFibAfterFour := First(Skip(Fibonaccis(), 4))
	fmt.Printf("First ten Fibonacci numbers, ignoring the first four, ie starting with %d:\n", firstFibAfterFour) // "ie starting with 5"
	first10FibsIgnoreFirstFour := Take(Skip(Fibonaccis(), 4), 10)
	fmt.Println(concatInts(", ", first10FibsIgnoreFirstFour)) // prints "5, 8, 13, 21, 34, 55, 89, 144, 233, 377"

	fmt.Println("Squares of first ten Fibonacci numbers")
	squareFirst10Fibs := Take(Map(Fibonaccis(), square), 10) // Note: mapping before taking. Can we call map on an unending stream of data?
	fmt.Println(concatInts(", ", squareFirst10Fibs))         // Indeed we can; this line prints "1, 1, 4, 9, 25, 64, 169, 441, 1156, 3025"

	fmt.Println("Total number of Fibonacci numbers:")
	fibCount := countOrTimeOut(Fibonaccis(), 1)
	if fibCount >= 0 {
		fmt.Println(fibCount)
	} else {
		fmt.Println("Timed out, obviously") // Times out, obviously
	}

	fmt.Println("Multiply each integer in the test set by the subsequent integer:")
	product := func(a int, b int) int { return a * b }
	offsetProducts := Zip(From(ints), Skip(From(ints), 1), product)
	fmt.Println(concatInts(", ", offsetProducts)) // prints "2, 6, 18, 24, 4, 9, 45, 40"

	fmt.Println("Successive ratios of the five Fibonacci numbers after skipping the first five:")
	ratio := func(a int, b int) float64 { return float64(b) / float64(a) }
	fibs := Fibonaccis()
	fibs2 := Skip(Fibonaccis(), 1)
	phiApproximations := Take(Skip(Zip(fibs, fibs2, ratio), 5), 5)
	fmt.Println(concatFloats(", ", phiApproximations)) // prints "1.625000, 1.615385, 1.619048, 1.617647, 1.618182"
	close(fibs)
	close(fibs2)
}
