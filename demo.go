package main

import (
	"context"
	"fmt"
	"golinq/gl"
	"strings"
	"time"
)

// Listens on a channel until the channel is closed or a timeout threshold is reached
// and return the number of elements received, or -1 if timed out.
func countOrTimeOut[T any](source chan T, timeoutSec int) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()
	ch := make(chan int)
	go func() { ch <- gl.Count(source) }()
	select {
	case res := <-ch:
		return res // counted successfully
	case <-ctx.Done():
		return -1 // timed out
	}
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

// Informal tests to check behavior of GoLinq methods
func main() {
	ints := []int{1, 2, 3, 6, 4, 1, 9, 5, 8}

	fmt.Println("Given ints:")
	fmt.Println(concatInts(", ", gl.From(ints))) // prints "1, 2, 3, 6, 4, 1, 9, 5, 8"

	isEven := func(i int) bool { return i%2 == 0 }
	square := func(i int) int { return i * i }

	fmt.Println("Squares of ints: ")
	squares := gl.Map(gl.From(ints), square)
	fmt.Println(concatInts(", ", squares)) // prints "1, 4, 9, 36, 16, 1, 81, 25, 64"

	fmt.Println("Even squares of given ints:")
	squaresOfEvens := concatInts(", ", gl.Filter(gl.Map(gl.From(ints), square), isEven))
	fmt.Println(squaresOfEvens) // prints "4, 36, 16, 64"

	fmt.Println("Max of given ints:")
	max := gl.Max(gl.From(ints))
	fmt.Println(max) // prints "9"

	fmt.Println("First int:")
	first := gl.First(gl.From(ints))
	fmt.Println(first) // prints "1"

	fmt.Println("Last int:")
	last := gl.Last(gl.From(ints))
	fmt.Println(last) // prints "8"

	fmt.Println("Count of ints:")
	count := gl.Count(gl.From(ints))
	fmt.Println(count) // prints "9"

	fmt.Println("Sum of ints:")
	sum := gl.Sum(gl.From(ints))
	fmt.Println(sum) // prints "39"

	fmt.Println("Sum of first three ints:")
	sum3 := gl.Sum(gl.Take(gl.From(ints), 3))
	fmt.Println(sum3) // prints "6"

	fmt.Println("Sum of final two ints:")
	sumFinal2 := gl.Sum(gl.Skip(gl.From(ints), count-2))
	fmt.Println(sumFinal2) // prints "13"

	fmt.Println("First ten Fibonacci numbers")
	first10Fibs := gl.Take(gl.Fibonaccis(), 10)
	fmt.Println(concatInts(", ", first10Fibs)) // prints "1, 1, 2, 3, 5, 8, 13, 21, 34, 55"

	firstFibAfterFour := gl.First(gl.Skip(gl.Fibonaccis(), 4))
	fmt.Printf("First ten Fibonacci numbers, ignoring the first four, ie starting with %d:\n", firstFibAfterFour) // "ie starting with 5"
	first10FibsIgnoreFirstFour := gl.Take(gl.Skip(gl.Fibonaccis(), 4), 10)
	fmt.Println(concatInts(", ", first10FibsIgnoreFirstFour)) // prints "5, 8, 13, 21, 34, 55, 89, 144, 233, 377"

	fmt.Println("Squares of first ten Fibonacci numbers")
	squareFirst10Fibs := gl.Take(gl.Map(gl.Fibonaccis(), square), 10) // Note: mapping before taking. Can we call map on an unending stream of data?
	fmt.Println(concatInts(", ", squareFirst10Fibs))                  // Indeed we can; this line prints "1, 1, 4, 9, 25, 64, 169, 441, 1156, 3025"

	fmt.Println("Total number of Fibonacci numbers:")
	fibCount := countOrTimeOut(gl.Fibonaccis(), 1)
	if fibCount >= 0 {
		fmt.Println(fibCount)
	} else {
		fmt.Println("Timed out, obviously") // Times out, obviously
	}

	fmt.Println("Multiply each integer in the test set by the subsequent integer:")
	product := func(a int, b int) int { return a * b }
	offsetProducts := gl.Zip(gl.From(ints), gl.Skip(gl.From(ints), 1), product)
	fmt.Println(concatInts(", ", offsetProducts)) // prints "2, 6, 18, 24, 4, 9, 45, 40"

	fmt.Println("Successive ratios of the five Fibonacci numbers after skipping the first five:")
	ratio := func(a int, b int) float64 { return float64(b) / float64(a) }
	fibs := gl.Fibonaccis()
	fibs2 := gl.Skip(gl.Fibonaccis(), 1)
	phiApproximations := gl.Take(gl.Skip(gl.Zip(fibs, fibs2, ratio), 5), 5)
	fmt.Println(concatFloats(", ", phiApproximations)) // prints "1.625000, 1.615385, 1.619048, 1.617647, 1.618182"
	close(fibs)
	close(fibs2)
}
