package codes

import (
	"fmt"
	"sort"
)

/*
56. 合并区间：以数组 intervals 表示若干个区间的集合，
其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回一个不重叠的区间数组，
该数组需恰好覆盖输入中的所有区间。
可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，
遍历排序后的区间数组，将当前区间与切片中最后一个区间进行比较，如果有重叠，则合并区间；
如果没有重叠，则将当前区间添加到切片中。
*/

type itvS struct {
	starti int
	endi   int
}

func (t *itvS) SetNums(nums []int) bool {
	if len(nums) == 2 {
		t.starti = nums[0]
		t.endi = nums[1]

		if t.starti < t.endi {
			return true
		}
	}
	return false
}

func (t *itvS) SetItvS(a itvS) {
	t.starti = a.starti
	t.endi = a.endi
}

func (t *itvS) Marge(n *itvS) bool {
	if t.starti == n.starti && t.endi == n.endi { // 完全相等则无需合并
		return false
	}

	if t.isIn(n.starti) {
		if t.isIn(n.endi) {
			// 覆盖，不做处理
			return true
		} else {
			// endi 超出范围
			t.endi = n.endi
			return true
		}
	} else {
		if t.isIn(n.endi) {
			t.starti = n.starti
			return true
		}
	}
	return false
}

func (t itvS) isIn(n int) bool {
	if t.starti <= n && n <= t.endi {
		return true
	} else {
		return false
	}
}

func (t itvS) getList() []int {
	return []int{t.starti, t.endi}
}

var C7_t1 [][]int = [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
var C7_t2 [][]int = [][]int{{10, 30}, {50, 70}, {60, 80}, {90, 120}, {140, 160},
	{180, 200}, {220, 250}, {260, 280}, {300, 320}, {310, 330}, {25, 40}, {110, 150},
	{350, 380}, {390, 420}, {450, 470}, {460, 490}, {500, 530}, {230, 270}, {360, 400},
	{540, 560}, {600, 620}, {650, 680}, {720, 750}, {800, 820}, {670, 700}, {810, 850},
}
var C7_t3 [][]int = [][]int{{1, 4}, {3, 6}, {2, 9}, {10, 19}, {7, 14}}

func removeAtIndex(tx []itvS, index int) []itvS {
	resDa := []itvS{}
	for i := range tx {
		if i != index {
			resDa = append(resDa, tx[i])
		}
	}
	return resDa
}

func MergeNums(intervals [][]int) [][]int {

	fmt.Println("合并前的数组长度", len(intervals), "数据列表为：", intervals)

	var itvList []itvS = []itvS{}
	var a = itvS{}
	for _, nums := range intervals {
		if isOk := a.SetNums(nums); !isOk {
			continue
		}
		for {
			var isMarge bool = false
			for index := range itvList {
				pa := &(itvList[index])
				if isMarge = pa.Marge(&a); isMarge {
					a.SetItvS(*pa)
					itvList = removeAtIndex(itvList, index)
					break
				}
			}
			if !isMarge {
				break
			}
		}
		itvList = append(itvList, a)
	}
	var res [][]int = [][]int{}
	for _, a := range itvList {
		res = append(res, a.getList())
	}

	// 排个序
	// 按区间起始位置排序
	sort.Slice(res, func(i, j int) bool {
		return res[i][0] < res[j][0]
	})

	fmt.Println("合并后的数组长度", len(res), "数据列表为：", res)
	return res
}

func C7_ft1() {
	var a = itvS{}
	var b = itvS{}
	var nums = []int{90, 150} //[90 150] [140 160]
	var nums2 = []int{140, 160}
	a.SetNums(nums)
	fmt.Println("a is", a.getList())
	b.SetNums(nums2)
	fmt.Println("b is", b.getList())
	a.Marge(&b)
	fmt.Println("a is", a.getList())
}

// 优秀代码
func C7_marge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return intervals
	}

	// 按区间起始点排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	fmt.Println("排序后的数组长度", len(intervals), "数据列表为：", intervals)
	merged := make([][]int, 0)
	merged = append(merged, intervals[0])

	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		current := intervals[i]

		if last[1] < current[0] {
			// 无重叠，直接添加
			merged = append(merged, current)
		} else {
			// 合并区间
			if current[1] > last[1] {
				last[1] = current[1]
			}
		}
	}
	fmt.Println("合并后的数组长度", len(merged), "数据列表为：", merged)
	return merged
}
