package ungo

import (
	"math/rand"
	"time"
)

func NumToFloat64[T Number](n T) float64 {
	switch v := any(n).(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case complex64:
		return float64(real(v))
	case complex128:
		return float64(real(v))
	default:
		return 0
	}
}

func StalinSort[T Number](arr []T) []T {
	result := make([]T, 0, len(arr))
	for i, v := range arr {
		if i == 0 || NumToFloat64(v) >= NumToFloat64(arr[i-1]) {
			result = append(result, v)
		}
	}
	return result
}

func QuickSort[T Number](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}
	pivot := arr[len(arr)/2]
	left := make([]T, 0, len(arr))
	right := make([]T, 0, len(arr))
	for _, v := range arr {
		if NumToFloat64(v) < NumToFloat64(pivot) {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	return append(append(QuickSort(left), pivot), QuickSort(right)...)
}

func HeapSort[T Number](arr []T) []T {
	heap := make([]T, len(arr))
	copy(heap, arr)
	for i := len(heap)/2 - 1; i >= 0; i-- {
		heapify(heap, len(heap), i)
	}
	for i := len(heap) - 1; i > 0; i-- {
		heap[0], heap[i] = heap[i], heap[0]
		heapify(heap, i, 0)
	}
	return heap
}

func heapify[T Number](arr []T, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2
	if left < n && NumToFloat64(arr[left]) > NumToFloat64(arr[largest]) {
		largest = left
	}
	if right < n && NumToFloat64(arr[right]) > NumToFloat64(arr[largest]) {
		largest = right
	}
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

func MergeSort[T Number](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])
	return merge(left, right)
}

func merge[T Number](left, right []T) []T {
	result := make([]T, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if NumToFloat64(left[i]) < NumToFloat64(right[j]) {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func InPlaceMergeSort[T Number](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := InPlaceMergeSort(arr[:mid])
	right := InPlaceMergeSort(arr[mid:])
	return inPlaceMerge(left, right)
}

func inPlaceMerge[T Number](left, right []T) []T {
	result := make([]T, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if NumToFloat64(left[i]) < NumToFloat64(right[j]) {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func TournamentSort[T Number](arr []T) []T {
	n := len(arr)
	if n <= 1 {
		return arr
	}
	return tournamentSort(arr, 0, n-1)
}

func tournamentSort[T Number](arr []T, start, end int) []T {
	if start == end {
		return []T{arr[start]}
	}
	mid := (start + end) / 2
	left := tournamentSort(arr, start, mid)
	right := tournamentSort(arr, mid+1, end)
	return merge(left, right)
}

func TreeSort[T Number](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}
	root := &treeNode[T]{Value: arr[0]}
	for _, v := range arr[1:] {
		insert(root, v)
	}
	return inOrder(root)
}

type treeNode[T Number] struct {
	Value T
	Left  *treeNode[T]
	Right *treeNode[T]
}

func insert[T Number](node *treeNode[T], value T) {
	if NumToFloat64(value) < NumToFloat64(node.Value) {
		if node.Left == nil {
			node.Left = &treeNode[T]{Value: value}
		} else {
			insert(node.Left, value)
		}
	} else {
		if node.Right == nil {
			node.Right = &treeNode[T]{Value: value}
		} else {
			insert(node.Right, value)
		}
	}
}

func inOrder[T Number](node *treeNode[T]) []T {
	if node == nil {
		return nil
	}
	result := inOrder(node.Left)
	result = append(result, node.Value)
	result = append(result, inOrder(node.Right)...)
	return result
}

func BlockSort[T Number](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}
	blockSize := 16
	blocks := make([][]T, 0, (len(arr)+blockSize-1)/blockSize)
	for i := 0; i < len(arr); i += blockSize {
		end := i + blockSize
		if end > len(arr) {
			end = len(arr)
		}
		blocks = append(blocks, arr[i:end])
	}
	for i := 0; i < len(blocks); i++ {
		blocks[i] = InsertionSort(blocks[i])
	}
	result := make([]T, 0, len(arr))
	for _, block := range blocks {
		result = append(result, block...)
	}
	return result
}

func InsertionSort[T Number](arr []T) []T {
	for i := 1; i < len(arr); i++ {
		j := i
		for j > 0 && NumToFloat64(arr[j]) < NumToFloat64(arr[j-1]) {
			arr[j], arr[j-1] = arr[j-1], arr[j]
			j--
		}
	}
	return arr
}

func PatienceSort[T Number](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}

	stacks := make([][]T, 0, len(arr))
	for _, v := range arr {
		insertPatience(&stacks, v)
	}
	result := make([]T, 0, len(arr))
	for _, stack := range stacks {
		result = append(result, stack...)
	}
	return result
}

func insertPatience[T Number](stacks *[][]T, value T) {
	for _, stack := range *stacks {
		if NumToFloat64(value) >= NumToFloat64(stack[len(stack)-1]) {
			stack = append(stack, value)
			return
		}
	}
	*stacks = append(*stacks, []T{value})
}

func BubbleSort[T Number](arr []T) []T {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr)-i-1; j++ {
			if NumToFloat64(arr[j]) > NumToFloat64(arr[j+1]) {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

func CocktailSort[T Number](arr []T) []T {
	for i := 0; i < len(arr)/2; i++ {
		for j := i; j < len(arr)-i-1; j++ {
			if NumToFloat64(arr[j]) > NumToFloat64(arr[j+1]) {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
		for j := len(arr) - i - 2; j > i; j-- {
			if NumToFloat64(arr[j]) < NumToFloat64(arr[j-1]) {
				arr[j], arr[j-1] = arr[j-1], arr[j]
			}
		}
	}
	return arr
}

func GnomeSort[T Number](arr []T) []T {
	i := 0
	for i < len(arr) {
		if i == 0 || NumToFloat64(arr[i]) >= NumToFloat64(arr[i-1]) {
			i++
		} else {
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
		}
	}
	return arr
}

func OddEvenSort[T Number](arr []T) []T {
	for i := 0; i < len(arr); i++ {
		if i%2 == 0 {
			for j := i + 1; j < len(arr); j++ {
				if NumToFloat64(arr[j]) < NumToFloat64(arr[i]) {
					arr[j], arr[i] = arr[i], arr[j]
				}
			}
		} else {
			for j := i + 1; j < len(arr); j++ {
				if NumToFloat64(arr[j]) > NumToFloat64(arr[i]) {
					arr[j], arr[i] = arr[i], arr[j]
				}
			}
		}
	}
	return arr
}

func Bogosort[T Number](arr []T) []T {
	for !isSorted(arr) {
		arr = shuffle(arr)
	}
	return arr
}

func isSorted[T Number](arr []T) bool {
	for i := 1; i < len(arr); i++ {
		if NumToFloat64(arr[i]) < NumToFloat64(arr[i-1]) {
			return false
		}
	}
	return true
}

func shuffle[T Number](arr []T) []T {
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func SleepSort[T Number](arr []T) []T {
	for i := range arr {
		go func(j int) {
			time.Sleep(time.Duration(NumToFloat64(arr[j])) * time.Millisecond)
			arr[i] = arr[j]
		}(i)
	}
	return arr
}

// VoidSort is a sorting algorithm that uses the void transformation to sort numbers.
func VoidSort[T Number](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}

	// 1. Convert the input array into a slice of UnusualNum
	// This lifts our certain numbers into the realm of possibilities.
	unusuals := make([]*UnusualNum, len(arr))
	for i, val := range arr {
		u := NewUnusualNum()
		AddUnknownPossibility(u, val)
		unusuals[i] = u
	}

	// 2. The Void Transformation
	// We divide every number by zero to expand its state into the
	// {inf, -inf, 0, real, imag} set you described.
	zero := NewUnusualNum()
	AddUnknownPossibility(zero, 0)

	for i := range unusuals {
		// x / 0 creates the 9-point set of nonsense
		unusuals[i] = unusuals[i].Divide(zero)

		// We "stabilize" it by adding its own absolute value
		// then multiplying by its sine to ensure maximum floating point drift.
		unusuals[i] = unusuals[i].Add(unusuals[i].Abs()).Sin()
	}

	// 3. Sorting by the "Forced First Real" component
	// Since we've expanded the possibilities, we sort by the first
	// arbitrary value that "collapses" out of the set.
	for i := 0; i < len(unusuals); i++ {
		for j := i + 1; j < len(unusuals); j++ {
			// Compare based on the real part of the first possibility in the set
			if unusuals[i].ForceFirstReal() > unusuals[j].ForceFirstReal() {
				unusuals[i], unusuals[j] = unusuals[j], unusuals[i]
				// Swap the original array to match the "void" order
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}

	return arr
}

func SlowSort[T Number](arr []T, i T, j T) []T {
	if NumToFloat64(i) > NumToFloat64(j) {
		return arr
	}

	mid := (i + j) / 2
	arr = SlowSort(arr, i, mid)
	arr = SlowSort(arr, mid+1, j)

	if NumToFloat64(arr[int(NumToFloat64(mid))]) > NumToFloat64(arr[int(NumToFloat64(mid))+1]) {
		arr[int(NumToFloat64(mid))], arr[int(NumToFloat64(mid))+1] = arr[int(NumToFloat64(mid))+1], arr[int(NumToFloat64(mid))]
	}

	arr = SlowSort(arr, i, j-1)

	return arr
}

// MiracleSort hopes the input is sorted.
func MiracleSort[T Number](arr []T) []T {
	return arr
}

// Wait until a cosmic ray flips the correct bits to sort the array.
func CosmicRaySorting[T Number](arr []T) []T {
	for !isSorted(arr) {
		time.Sleep(time.Second * 1)
	}
	return arr
}

func StoogeSort[T Number](arr []T, i T, j T) []T {
	if NumToFloat64(i) > NumToFloat64(j) {
		return arr
	}

	if NumToFloat64(arr[int(NumToFloat64(i))]) > NumToFloat64(arr[int(NumToFloat64(j))]) {
		arr[int(NumToFloat64(i))], arr[int(NumToFloat64(j))] = arr[int(NumToFloat64(j))], arr[int(NumToFloat64(i))]
	}

	if NumToFloat64(j)-NumToFloat64(i) > 1 {
		mid := (j - i) / 3
		arr = StoogeSort(arr, i, j-mid)
		arr = StoogeSort(arr, i+mid, j)
		arr = StoogeSort(arr, i, j-mid)
	}

	return arr
}
