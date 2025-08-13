package main

import (
	"fmt"
	"math/rand"
	"time"
)

type runResult struct {
	run                         int
	arrayGen, sortDis, primeDis float64
	totalDis                    float64
	sortOrd, primeOrd           float64
	totalOrd                    float64
	primesDis, primesOrd        int
	comparisonsDisordered       int64
	comparisonsOrdered          int64
}

var comparisonCounter int64

func timeTrack(start time.Time) float64 {
	return toMilliseconds(time.Since(start))
}

func toMilliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

func main() {
	size := 1_000
	fmt.Printf("\n+------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-------------------+--------------+\n")
	fmt.Printf("| Array size: %d                                                                                                                                                      |\n", size)

	var results []runResult

	for i := range 10 {
		res := runResult{run: i + 1}

		start := time.Now()

		startGen := time.Now()
		arr := generateArrayWithRandomNumbers(size)
		res.arrayGen = timeTrack(startGen)

		comparisonCounter = 0
		startSort := time.Now()
		arrSorted := make([]int, len(arr))
		copy(arrSorted, arr)
		arrSorted = quickSortStart(arrSorted)
		res.sortDis = timeTrack(startSort)
		res.comparisonsDisordered = comparisonCounter

		startPrime := time.Now()
		primes := findPrimes(arrSorted)
		res.primesDis = len(primes)
		res.primeDis = timeTrack(startPrime)

		res.totalDis = timeTrack(start)

		startOrdered := time.Now()

		comparisonCounter = 0
		startSortOrdered := time.Now()
		arrSortedCopy := make([]int, len(arrSorted))
		copy(arrSortedCopy, arrSorted)
		arrSortedCopy = quickSortStart(arrSortedCopy)
		res.sortOrd = timeTrack(startSortOrdered)
		res.comparisonsOrdered = comparisonCounter

		startPrimeOrdered := time.Now()
		primesOrdered := findPrimes(arrSortedCopy)
		res.primesOrd = len(primesOrdered)
		res.primeOrd = timeTrack(startPrimeOrdered)

		res.totalOrd = timeTrack(startOrdered)

		results = append(results, res)
	}

	fmt.Printf("+------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-------------------+--------------+\n")
	fmt.Printf("| Run  | Array Gen ms    | Sort D ms       | Prime D ms      | Total D ms      | Sort O ms       | Prime O ms      | Total O ms      | Comparisons D/O   | Primes D/O   |\n")
	fmt.Printf("+------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-------------------+--------------+\n")
	for _, r := range results {
		fmt.Printf("| %-4d | %-15s | %-15s | %-15s | %-15s | %-15s | %-15s | %-15s | %-17s | %-12s |\n",
			r.run,
			fmt.Sprintf("%.6fms", r.arrayGen),
			fmt.Sprintf("%.6fms", r.sortDis),
			fmt.Sprintf("%.6fms", r.primeDis),
			fmt.Sprintf("%.6fms", r.totalDis),
			fmt.Sprintf("%.6fms", r.sortOrd),
			fmt.Sprintf("%.6fms", r.primeOrd),
			fmt.Sprintf("%.6fms", r.totalOrd),
			fmt.Sprintf("%d / %d", r.comparisonsDisordered, r.comparisonsOrdered),
			fmt.Sprintf("%d / %d", r.primesDis, r.primesOrd))
	}
	fmt.Printf("+------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-----------------+-------------------+--------------+\n")

	var totalDisord, totalOrd int64
	var totalPrimesDisord, totalPrimesOrd int
	for _, r := range results {
		totalDisord += r.comparisonsDisordered
		totalOrd += r.comparisonsOrdered
		totalPrimesDisord += r.primesDis
		totalPrimesOrd += r.primesOrd
	}
	avgDisord := float64(totalDisord) / float64(len(results))
	avgOrd := float64(totalOrd) / float64(len(results))
	avgPrimesDisord := float64(totalPrimesDisord) / float64(len(results))
	avgPrimesOrd := float64(totalPrimesOrd) / float64(len(results))

	fmt.Printf("\nMédia de comparações (Desordenado): %.2f\n", avgDisord)
	fmt.Printf("Média de comparações (Ordenado): %.2f\n", avgOrd)
	fmt.Printf("Diferença percentual: %.2f%%\n", ((avgDisord-avgOrd)/avgOrd)*100)
	fmt.Printf("Média de primos (Desordenado): %.2f\n", avgPrimesDisord)
	fmt.Printf("Média de primos (Ordenado): %.2f\n", avgPrimesOrd)
}

func generateArrayWithRandomNumbers(size int) []int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = r.Intn(999_999_999)
	}
	return arr
}

func medianOfThree(array []int, low, high int) int {
	mid := (low + high) / 2

	comparisonCounter++
	if array[low] > array[mid] {
		array[low], array[mid] = array[mid], array[low]
	}

	comparisonCounter++
	if array[mid] > array[high] {
		array[mid], array[high] = array[high], array[mid]
	}

	comparisonCounter++
	if array[low] > array[mid] {
		array[low], array[mid] = array[mid], array[low]
	}

	array[mid], array[high] = array[high], array[mid]
	return mid
}

func partition(array []int, low, high int) int {
	medianOfThree(array, low, high)
	pivot := array[high]

	i := low
	for j := low; j < high; j++ {
		comparisonCounter++
		if array[j] < pivot {
			array[i], array[j] = array[j], array[i]
			i++
		}
	}
	array[i], array[high] = array[high], array[i]
	return i
}

func quickSort(arr []int, low, high int) {
	comparisonCounter++
	if low < high {
		p := partition(arr, low, high)
		quickSort(arr, low, p-1)
		quickSort(arr, p+1, high)
	}
}

func quickSortStart(arr []int) []int {
	if len(arr) > 0 {
		quickSort(arr, 0, len(arr)-1)
	}
	return arr
}

func findPrimes(arr []int) []int {
	primes := make([]int, 0)
	for _, num := range arr {
		if isPrime(num) {
			primes = append(primes, num)
		}
	}
	return primes
}

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	if num <= 3 {
		return true
	}
	if num%2 == 0 || num%3 == 0 {
		return false
	}

	for i := 5; i*i <= num; i += 6 {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
	}
	return true
}
