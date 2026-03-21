package homework01

import (
	// "cmp"
	"sort"
	"strconv"
)

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {
	count := make(map[int]int)
	var singlev int
	for _, value := range nums {
		count[value] += 1
	}
	for k, v := range count {
		if v == 1 {
			singlev = k
		}

	}
	return singlev
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {
	arr := strconv.Itoa(x)
	arr_len := len(arr)
	isPalindrome := true

	for i := 0; i < arr_len/2; i++ {
		if arr[i] != arr[arr_len-i-1] {
			isPalindrome = false
			break
		}
	}
	return isPalindrome
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	// TODO: implement
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	var stack []rune

	for _, char := range s {
		left, ok := pairs[char]
		if ok {
			if len(stack) == 0 || stack[len(stack)-1] != left {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, char)
		}
	}
	return len(stack) == 0
}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]

	for i := 0; i < len(prefix); i++ {

		char := prefix[i]
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || strs[j][i] != char {
				return prefix[:i]
			}
		}

	}
	return prefix
}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i]+1 == 10 {

			digits[i] = 0
			if i == 0 {
				sliceDig := digits[:]
				sliceDig = append([]int{1}, sliceDig...)
				return sliceDig
				// digits = append(digits,)
			}
		} else {
			digits[i] += 1
			break
		}
	}
	return digits
}

// 6. 删除有序数组中的重复项
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	i := 0

	for j := 1; j < len(nums); j++ {
		if nums[i] != nums[j] {
			nums[i+1] = nums[j]
			i++
		}

	}

	return i + 1
}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func Merge(intervals [][]int) [][]int {
	merged := [][]int{}
	if len(intervals) == 0 {
		return merged
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	merged = append(merged, intervals[0])
	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		if last[1] < intervals[i][0] {
			merged = append(merged, intervals[i])

		} else {
			if last[1] < intervals[i][1] {
				last[1] = intervals[i][1]
			}
			// last[1] = cmp.Max(last[1], intervals[i][1])
		}
	}
	return merged

}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for idx, value := range nums {
		if idx2, ok := numMap[target-value]; ok {
			return []int{idx2, idx}
		}
		numMap[value] = idx

	}
	return nil
}
