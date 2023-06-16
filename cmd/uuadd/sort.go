package main

import (
  "sort"
)

type ByKey []PageDataType

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sortPageList() {
  sort.Sort(ByKey(pageList))
}
