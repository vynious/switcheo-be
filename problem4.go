package main

import (
	"fmt"
	"sync"
)

func sum_to_n_a(n int) int {
	/*
		using a mathematical formula to derive the sum of the integers to n

		time complexity : o(1)
		space complexity : o(1)
	*/

	result := n * (n + 1) / 2
	return result
}

func sum_to_n_b(n int) int {
	/*
		using goroutines and channel to make an iterative process to derive the sum of the integers to n

		time complexity : o(n) => iterating over n but concurrently
		space complexity : o(n) => uses goroutines and channels
	*/
	ch := make(chan int)
	go func() {
		total := 0
		for i := 1; i <= n; i++ {
			total += i
		}
		ch <- total
		close(ch)
	}()
	return <-ch
}

func sum_to_n_c(n int) int {
	/*
		using goroutines and channel to make a recursive process to derive the sum of the integers to n

		time complexity : o(n)
		space complexity : o(n)
	*/

	if n <= 1 {
		return n
	}

	ch := make(chan int)
	wg := new(sync.WaitGroup)

	// Start a separate goroutine to collect results
	result := 0
	done := make(chan bool)

	go func() {
		for i := range ch {
			result += i
		}
		done <- true
	}()

	wg.Add(1)

	// Start the recursive function in a goroutine
	go recursiveSumWithConcurrency(n, ch, wg)

	wg.Wait()
	close(ch) // Close the channel after all sends are done

	<-done // Wait for the result collecting goroutine to finish
	return result
}

func recursiveSumWithConcurrency(n int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	if n <= 1 {
		ch <- n
		return
	}

	wg.Add(1)
	go recursiveSumWithConcurrency(n-1, ch, wg)

	ch <- n
}

func main() {
	fmt.Printf("the sum of n = 4 is, %v", sum_to_n_a(4))
	fmt.Println()
	fmt.Printf("the sum of n = 4 is, %v", sum_to_n_b(4))
	fmt.Println()
	fmt.Printf("the sum of n = 4 is, %v", sum_to_n_c(4))

}
