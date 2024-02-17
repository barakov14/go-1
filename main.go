package main

import "fmt"

func SortSlice(slice []int) {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[j] < slice[i] {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

func IncrementOdd(slice []int) {
	for i := 0; i < len(slice); i++ {
		if i%2 != 0 {
			slice[i]++
		}
	}
}

func PrintSlice(slice []int) {
	fmt.Println(slice)
}

func ReverseSlice(slice []int) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func appendFunc(dst func([]int), src ...func([]int)) func([]int) {
	return func(slice []int) {
		dst(slice)
		for _, f := range src {
			f(slice)
		}
	}
}

func main() {
	slice := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}

	SortSlice(slice)
	PrintSlice(slice)

	IncrementOdd(slice)
	PrintSlice(slice)

	ReverseSlice(slice)
	PrintSlice(slice)

	dstFunc := func(slice []int) {
		fmt.Println("Destination function:", slice)
	}
	srcFunc1 := func(slice []int) {
		fmt.Println("Source function 1:", slice)
	}
	srcFunc2 := func(slice []int) {
		fmt.Println("Source function 2:", slice)
	}
	newFunc := appendFunc(dstFunc, srcFunc1, srcFunc2)

	newFunc(slice)
}
