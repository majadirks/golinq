# GoLinq

This project is an first experiment in learning Go. It provides a set of utility methods inspired by the C# LINQ methods for acting lazily on streams of data: filtering, mapping, aggregating, and so forth. For example, the following code uses ratios of Fibonacci numbers to approximate the golden ratio:

```
  fmt.Println("Successive ratios of the five Fibonacci numbers after skipping the first five:")
	ratio := func(a int, b int) float64 { return float64(b) / float64(a) }
	fibs := gl.Fibonaccis()
	fibs2 := gl.Skip(gl.Fibonaccis(), 1)
	phiApproximations := gl.Take(gl.Skip(gl.Zip(fibs, fibs2, ratio), 5), 5)
	fmt.Println(concatFloats(", ", phiApproximations)) // prints "1.625000, 1.615385, 1.619048, 1.617647, 1.618182"
```

Where C# uses objects implementing the `IEnumerable` interface, this Go analogue acts on channels and uses goroutines to provide lazy evaluation.

# Methods and Examples
The methods included are:
- `From`, which converts a slice to a channel, as in:
  ```
  ints := []int{1, 2, 3, 6, 4, 1, 9, 5, 8}
	fmt.Println("Given ints:")
	fmt.Println(concatInts(", ", gl.From(ints))) // prints "1, 2, 3, 6, 4, 1, 9, 5, 8"
  ```

- `Map`, which applies a given function to each successive value in a channel, as in:
  ```
	square := func(i int) int { return i * i }
	squares := gl.Map(gl.From(ints), square)
	fmt.Println(concatInts(", ", squares)) // prints "1, 4, 9, 36, 16, 1, 81, 25, 64"
  ```
- `Filter`, which returns a channel of values that match a given predicate, as in:
  ```
  isEven := func(i int) bool { return i%2 == 0 }
	squaresOfEvens := concatInts(", ", gl.Filter(gl.Map(gl.From(ints), square), isEven))
	fmt.Println(squaresOfEvens) // prints "4, 36, 16, 64"
  ```
- `Take`, which receives the first n values from a channel and sends them on a new channel, as in:
  ```
	sum3 := gl.Sum(gl.Take(gl.From(ints), 3))
	fmt.Println(sum3) // prints "6"
  ```
- `Skip`, which ignores the first n values from a channel and sends the rest (if any) on a new channel, as in:
  ```
  first10Fibs := gl.Take(gl.Fibonaccis(), 10)
	fmt.Println(concatInts(", ", first10Fibs)) // prints "1, 1, 2, 3, 5, 8, 13, 21, 34, 55"
  ```
- `Fibonaccis()`, which creates a channel on which all Fibonacci numbers are (lazily) sent
- The aggregation methods `First`, `Last`, `Max`, `Count`, and `Sum`, which do exactly what one would expect.

# Demo
One may see this code in action by running `go run demo.go` in the appropriate directory.
