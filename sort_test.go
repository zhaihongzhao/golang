package sort

import (
	"math/rand"
	"testing"
)

func TestSort(t *testing.T) {
	fs := []func([]int){
		SelectSort,
		BubbleSort,
		InsertSort,
		ShellSort,
		HeapSort,
		MergeSort,
		QuickSort,
		RadixSort,
	}

	arr := make([]int, 100)
	for _, f := range fs {
		for i := range arr {
			arr[i] = rand.Intn(100)
		}
		f(arr)
		t.Log(arr)
		if !isSorted(arr) {
			t.Fail()
		}
	}
}

func isSorted(arr []int) bool {
	n := len(arr)
	for i := 1; i < n; i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}
