package main

import (
  "sort"
)

type ByKey []PageDataType

func (a ByKey) Len() int           { return len(a) }
// Note that we sort in descending order.
func (a ByKey) Less(i, j int) bool { return a[j].Key < a[i].Key }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sortPageList() {
  sort.Sort(ByKey(pageList))
}
