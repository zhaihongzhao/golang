package sort

import (
	"math/rand"
)

// 选择排序 时间O(n^2) 空间O(1) 稳定
func SelectSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		// 找出最小值与arr[i]交换
		min := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[min] {
				min = j
			}
		}
		arr[i], arr[min] = arr[min], arr[i]
	}
}

// 冒泡排序 时间O(n^2) 空间O(1) 稳定
func BubbleSort(arr []int) {
	n := len(arr)
	for i := n - 1; i > 0; i-- {
		// 相邻交换，把最大值交换到i位置
		swap := false
		for j := 0; j < i; j++ {
			if arr[j] > arr[j+1] {
				swap = true
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
		if !swap {
			// 没有交换，说明列表已经有序
			break
		}
	}
}

// 插入排序 时间O(n^2) 空间O(1) 稳定
func InsertSort(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		// 把arr[i]插入到前面的有序区间
		val := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > val {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = val
	}
}

// 希尔排序 平均时间O(n^1.5) 最坏时间O(n^2) 空间O(1) 不稳定
func ShellSort(arr []int) {
	n := len(arr)
	for gap := n / 2; gap > 0; gap /= 2 {
		// 间隔为gap的插入排序
		for i := gap; i < n; i++ {
			// 把arr[i]插入到前面的有序区间
			val := arr[i]
			j := i - gap
			for j >= 0 && arr[j] > val {
				arr[j+gap] = arr[j]
				j -= gap
			}
			arr[j+gap] = val
		}
	}
}

// 堆排序 时间O(nlogn) 空间O(1) 不稳定
func HeapSort(arr []int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	// 构建最大堆
	for i := n/2 - 1; i >= 0; i-- {
		down(arr, i, n)
	}
	// 模拟出堆
	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		down(arr, 0, i)
	}
}

func down(arr []int, i, n int) {
	// 最大堆的向下调整
	for {
		max := i
		left, right := i*2+1, i*2+2
		if left < n && arr[left] > arr[max] {
			max = left
		}
		if right < n && arr[right] > arr[max] {
			max = right
		}
		if max == i {
			break
		}
		// 把最大值交换上来，迭代向下调整
		arr[i], arr[max] = arr[max], arr[i]
		i = max
	}
}

// 归并排序 时间O(nlogn) 空间O(n) 稳定
func MergeSort(arr []int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	// 分成两半各自排序
	mid := n / 2
	MergeSort(arr[:mid])
	MergeSort(arr[mid:])
	// 合并两个有序列表
	res := merge(arr[:mid], arr[mid:])
	// 把结果拷贝回来
	copy(arr, res)
}

func merge(a, b []int) []int {
	// 合并两个有序列表
	m, n := len(a), len(b)
	res := make([]int, 0, m+n)
	i, j := 0, 0
	for i < m && j < n {
		if a[i] <= b[j] {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
	res = append(res, a[i:]...)
	res = append(res, b[j:]...)
	return res
}

// 快速排序 平均时间O(nlogn) 最坏时间O(n^2) 空间O(logn) 不稳定
func QuickSort(arr []int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	// 随机选基准元素x
	pos := rand.Intn(n)
	arr[0], arr[pos] = arr[pos], arr[0]
	x := arr[0]
	// 划分区间。把小于x的放在左区间，大于x的放在右区间
	i, j := 0, n-1
	for i < j {
		for i < j && arr[j] >= x {
			j--
		}
		for i < j && arr[i] <= x {
			i++
		}
		arr[i], arr[j] = arr[j], arr[i]
	}
	arr[0], arr[i] = arr[i], arr[0]
	// 分别对左右区间排序
	QuickSort(arr[:i])
	QuickSort(arr[i+1:])
}

// 基数排序 时间O(dn) 空间O(n) 稳定
func RadixSort(arr []int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	// 找最大值
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	// 从低位到高位进行d趟桶式排序（d是最大值的位数）
	for x := 1; x <= max; x *= 10 {
		var buckets [10][]int
		// 入桶
		for _, v := range arr {
			i := v / x % 10
			buckets[i] = append(buckets[i], v)
		}
		// 出桶
		i := 0
		for _, bucket := range buckets {
			i += copy(arr[i:], bucket)
		}
	}
}
