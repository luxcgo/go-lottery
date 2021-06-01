package main

import (
	"fmt"
	"strings"
	"time"
)

//组合算法(从nums中取出m个数)
func zuheResult(n int, m int) [][]int {
	if m < 1 || m > n {
		fmt.Println("Illegal argument. Param m must between 1 and len(nums).")
		return [][]int{}
	}
	//保存最终结果的数组，总数直接通过数学公式计算
	result := make([][]int, 0, mathZuhe(n, m))
	//保存每一个组合的索引的数组，1表示选中，0表示未选中
	indexs := make([]int, n)
	for i := 0; i < n; i++ {
		if i < m {
			indexs[i] = 1
		} else {
			indexs[i] = 0
		}
	}
	//第一个结果
	result = addTo(result, indexs)
	for {
		find := false
		//每次循环将第一次出现的 1 0 改为 0 1，同时将左侧的1移动到最左侧
		for i := 0; i < n-1; i++ {
			if indexs[i] == 1 && indexs[i+1] == 0 {
				find = true
				indexs[i], indexs[i+1] = 0, 1
				if i > 1 {
					moveOneToLeft(indexs[:i])
				}
				result = addTo(result, indexs)
				break
			}
		}
		//本次循环没有找到 1 0 ，说明已经取到了最后一种情况
		if !find {
			break
		}
	}
	return result
}

//将ele复制后添加到arr中，返回新的数组
func addTo(arr [][]int, ele []int) [][]int {
	newEle := make([]int, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)
	return arr
}
func moveOneToLeft(leftNums []int) {
	//计算有几个1
	sum := 0
	for i := 0; i < len(leftNums); i++ {
		if leftNums[i] == 1 {
			sum++
		}
	}
	//将前sum个改为1，之后的改为0
	for i := 0; i < len(leftNums); i++ {
		if i < sum {
			leftNums[i] = 1
		} else {
			leftNums[i] = 0
		}
	}
}

//根据索引号数组得到元素数组
func findNumsByIndexs(nums []int, indexs [][]int) [][]int {
	if len(indexs) == 0 {
		return [][]int{}
	}
	result := make([][]int, len(indexs))
	for i, v := range indexs {
		line := make([]int, 0)
		for j, v2 := range v {
			if v2 == 1 {
				line = append(line, nums[j])
			}
		}
		result[i] = line
	}
	return result
}

//数学方法计算排列数(从n中取m个数)
func mathPailie(n int, m int) int {
	return jieCheng(n) / jieCheng(n-m)
}

//数学方法计算组合数(从n中取m个数)
func mathZuhe(n int, m int) int {
	return jieCheng(n) / (jieCheng(n-m) * jieCheng(m))
}

//阶乘
func jieCheng(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func main1() {
	// nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	// m := 5
	// n := len(nums)
	// indexs := zuheResult(n, m)
	// for i := 0; i < len(indexs); i++ {
	// 	fmt.Println(indexs[i])
	// }
	Test10Base()
	fmt.Println(mmap)
	fmt.Println(len(mmap))
}

func Test10Base() {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	m := 5
	timeStart := time.Now()
	n := len(nums)
	indexs := zuheResult(n, m)
	result := findNumsByIndexs(nums, indexs)
	timeEnd := time.Now()
	fmt.Println("count:", len(result))
	fmt.Println("result:", result)
	fmt.Println("time consume:", timeEnd.Sub(timeStart))
	//结果是否正确
	rightCount := mathZuhe(n, m)
	if rightCount == len(result) {
		fmt.Println("结果正确")
	} else {
		fmt.Println("结果错误，正确结果是：", rightCount)
	}

	for _, value := range result {
		// fmt.Println(index, value)
		value := append(value, 10, 10, 10, 10, 10)
		arr := [10]int{}
		j := 0
		k := 5
		for i := 0; i < 10; i++ {

			if value[j] == i {
				arr[j] = i
				j++
			} else {
				arr[k] = i
				k++
			}
		}
		// fmt.Println(arr)
		s := strings.Replace(strings.Trim(fmt.Sprint(arr), "[]"), " ", "", -1)

		s1 := strings.Replace(strings.Trim(fmt.Sprint(arr[:5]), "[]"), " ", "", -1)
		s2 := strings.Replace(strings.Trim(fmt.Sprint(arr[5:]), "[]"), " ", "", -1)
		s3 := s2 + s1

		if nmap[s] == false {
			mmap[s] = arr
			nmap[s] = true
			nmap[s3] = true
		}
		// fmt.Println(s)
		// fmt.Println(s3)
	}

	//	for k, v := range mmap {
	//		fmt.Println(k, v)
	//	}
	//	fmt.Println(len(mmap))
}

var mmap = make(map[string][10]int)
var nmap = make(map[string]bool)

// 1. ten numbers
// 2. 排3 or 3D
// 3. what type? detailed or compressed 多四张图
// 4. method1 or method2

// 2. 排3 or 3D
// 1. ten numbers

/*
配置项
1. 排三原始数据（表头）
2. 排三组三和组六
3. 排三组三和组六压缩1
4. 排三组三和组六压缩2
5. 排三组三和组六压缩2延展
6. 排三组六（表头）
7. 排三组六压缩1
8. 排三组六压缩2
*/
