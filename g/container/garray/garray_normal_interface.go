// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/util/grand"
    "math"
    "sort"
    "strings"
)

type Array struct {
    mu    *rwmutex.RWMutex  // 互斥锁
    array []interface{}     // 底层数组
}

// Create an empty array.
// The param <unsafe> used to specify whether using array with un-concurrent-safety,
// which is false in default, means concurrent-safe in default.
//
// 创建一个空的数组对象，参数unsafe用于指定是否用于非并发安全场景，默认为false，表示并发安全。
func New(unsafe...bool) *Array {
    return NewArraySize(0, 0, unsafe...)
}

// See New.
//
// 同New方法。
func NewArray(unsafe...bool) *Array {
    return NewArraySize(0, 0, unsafe...)
}

// Create an array with given size and cap.
// The param <unsafe> used to specify whether using array with un-concurrent-safety,
// which is false in default, means concurrent-safe in default.
//
// 创建一个指定大小的数组对象，参数unsafe用于指定是否用于非并发安全场景，默认为false，表示并发安全。
func NewArraySize(size int, cap int, unsafe...bool) *Array {
    return &Array{
        mu    : rwmutex.New(unsafe...),
        array : make([]interface{}, size, cap),
    }
}

// Create an array with given slice <array>.
// The param <unsafe> used to specify whether using array with un-concurrent-safety,
// which is false in default, means concurrent-safe in default.
//
// 通过给定的slice变量创建数组对象，参数unsafe用于指定是否用于非并发安全场景，默认为false，表示并发安全。
func NewArrayFrom(array []interface{}, unsafe...bool) *Array {
    return &Array{
        mu    : rwmutex.New(unsafe...),
        array : array,
    }
}

// Get value by index.
//
// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *Array) Get(index int) interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    value := a.array[index]
    return value
}

// Set value by index.
//
// 设置指定索引的数据项, 调用方注意判断数组边界
func (a *Array) Set(index int, value interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array[index] = value
    return a
}

// Set the underlying slice array with the given <array> param.
//
// 设置底层数组变量.
func (a *Array) SetArray(array []interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array = array
    return a
}

// Replace the array items by given <array> from the beginning of array.
//
// 使用指定数组替换到对应的索引元素值.
func (a *Array) Replace(array []interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    max := len(array)
    if max > len(a.array) {
        max = len(a.array)
    }
    for i := 0; i < max; i++ {
        a.array[i] = array[i]
    }
    return a
}

// Calculate the sum of values in an array.
//
// 对数组中的元素项求和(将元素值转换为int类型后叠加)。
func (a *Array) Sum() (sum int) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, v := range a.array {
        sum += gconv.Int(v)
    }
    return
}

// Sort the array by custom function <less>.
//
// 使用自定义的排序函数将数组重新排序.
func (a *Array) SortFunc(less func(v1, v2 interface{}) bool) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    sort.Slice(a.array, func(i, j int) bool {
        return less(a.array[i], a.array[j])
    })
    return a
}

// Insert the <value> to the front of <index>.
//
// 在当前索引位置前插入一个数据项, 调用方注意判断数组边界。
func (a *Array) InsertBefore(index int, value interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    rear   := append([]interface{}{}, a.array[index : ]...)
    a.array = append(a.array[0 : index], value)
    a.array = append(a.array, rear...)
    return a
}

// Insert the <value> to the back of <index>.
//
// 在当前索引位置前插入一个数据项, 调用方注意判断数组边界。
func (a *Array) InsertAfter(index int, value interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    rear   := append([]interface{}{}, a.array[index + 1 : ]...)
    a.array = append(a.array[0 : index + 1], value)
    a.array = append(a.array, rear...)
    return a
}

// Remove an item by index.
//
// 删除指定索引的数据项, 调用方注意判断数组边界。
func (a *Array) Remove(index int) interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    // 边界删除判断，以提高删除效率
    if index == 0 {
        value  := a.array[0]
        a.array = a.array[1 : ]
        return value
    } else if index == len(a.array) - 1 {
        value  := a.array[index]
        a.array = a.array[: index]
        return value
    }
    // 如果非边界删除，会涉及到数组创建，那么删除的效率差一些
    value  := a.array[index]
    a.array = append(a.array[ : index], a.array[index + 1 : ]...)
    return value
}

// Push new items to the beginning of array.
//
// 将数据项添加到数组的最左端(索引为0)。
func (a *Array) PushLeft(value...interface{}) *Array {
    a.mu.Lock()
    a.array = append(value, a.array...)
    a.mu.Unlock()
    return a
}

// Push new items to the end of array.
//
// 将数据项添加到数组的最右端(索引为length - 1), 等于: Append。
func (a *Array) PushRight(value...interface{}) *Array {
    a.mu.Lock()
    a.array = append(a.array, value...)
    a.mu.Unlock()
    return a
}

// Pop an random item from array.
//
// 随机将一个数据项移出数组，并返回该数据项。
func (a *Array) PopRand() interface{} {
    return a.Remove(grand.Intn(len(a.array)))
}

// Pop an item from the beginning of array.
//
// 将最左端(索引为0)的数据项移出数组，并返回该数据项。
func (a *Array) PopLeft() interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    value  := a.array[0]
    a.array = a.array[1 : ]
    return value
}

// Pop an item from the end of array.
//
// 将最右端(索引为length - 1)的数据项移出数组，并返回该数据项。
func (a *Array) PopRight() interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    index  := len(a.array) - 1
    value  := a.array[index]
    a.array = a.array[: index]
    return value
}

// Pop <size> items from the beginning of array.
//
// 将最左端(首部)的size个数据项移出数组，并返回该数据项
func (a *Array) PopLefts(size int) []interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    length := len(a.array)
    if size > length {
        size = length
    }
    value  := a.array[0 : size]
    a.array = a.array[size : ]
    return value
}

// Pop <size> items from the end of array.
//
// 将最右端(尾部)的size个数据项移出数组，并返回该数据项
func (a *Array) PopRights(size int) []interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    index := len(a.array) - size
    if index < 0 {
        index = 0
    }
    value  := a.array[index :]
    a.array = a.array[ : index]
    return value
}

// Get items by range, returns array[start:end].
// Be aware that, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// 将最右端(尾部)的size个数据项移出数组，并返回该数据项
func (a *Array) Range(start, end int) []interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    length := len(a.array)
    if start > length || start > end {
        return nil
    }
    if start < 0 {
        start = 0
    }
    if end > length {
        end = length
    }
    array  := ([]interface{})(nil)
    if a.mu.IsSafe() {
        a.mu.RLock()
        defer a.mu.RUnlock()
        array = make([]interface{}, end - start)
        copy(array, a.array[start : end])
    } else {
        array = a.array[start : end]
    }
    return array
}

// See PushRight.
//
// 追加数据项, 等于: PushRight。
func (a *Array) Append(value...interface{}) *Array {
    a.PushRight(value...)
    return a
}

// Get the length of array.
//
// 数组长度。
func (a *Array) Len() int {
    a.mu.RLock()
    length := len(a.array)
    a.mu.RUnlock()
    return length
}

// Get the underlying data of array.
// Be aware that, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// 返回原始数据数组.
func (a *Array) Slice() []interface{} {
    array := ([]interface{})(nil)
    if a.mu.IsSafe() {
        a.mu.RLock()
        defer a.mu.RUnlock()
        array = make([]interface{}, len(a.array))
        copy(array, a.array)
    } else {
        array = a.array
    }
    return array
}

// Return a new array, which is a copy of current array.
//
// 克隆当前数组，返回当前数组的一个拷贝。
func (a *Array) Clone() (newArray *Array) {
    a.mu.RLock()
    array := make([]interface{}, len(a.array))
    copy(array, a.array)
    a.mu.RUnlock()
    return NewArrayFrom(array, !a.mu.IsSafe())
}

// Clear array.
//
// 清空数据数组
func (a *Array) Clear() *Array {
    a.mu.Lock()
    if len(a.array) > 0 {
        a.array = make([]interface{}, 0)
    }
    a.mu.Unlock()
    return a
}

// Check whether a value exists in the array.
//
// 查找指定数值是否存在
func (a *Array) Contains(value interface{}) bool {
    return a.Search(value) != -1
}

// Search array by <value>, returns the index of <value>, returns -1 if not exists.
//
// 查找指定数值的索引位置，返回索引位置，如果查找不到则返回-1
func (a *Array) Search(value interface{}) int {
    if len(a.array) == 0 {
        return -1
    }
    a.mu.RLock()
    result := -1
    for index, v := range a.array {
        if v == value {
            result = index
            break
        }
    }
    a.mu.RUnlock()

    return result
}

// Unique the array, clear repeated values.
//
// 清理数组中重复的元素项
func (a *Array) Unique() *Array {
    a.mu.Lock()
    for i := 0; i < len(a.array) - 1; i++ {
        for j := i + 1; j < len(a.array); j++ {
            if a.array[i] == a.array[j] {
                a.array = append(a.array[ : j], a.array[j + 1 : ]...)
            }
        }
    }
    a.mu.Unlock()
    return a
}

// Lock writing by callback function f.
//
// 使用自定义方法执行加锁修改操作
func (a *Array) LockFunc(f func(array []interface{})) *Array {
    a.mu.Lock(true)
    defer a.mu.Unlock(true)
    f(a.array)
    return a
}

// Lock reading by callback function f.
//
// 使用自定义方法执行加锁读取操作
func (a *Array) RLockFunc(f func(array []interface{})) *Array {
    a.mu.RLock(true)
    defer a.mu.RUnlock(true)
    f(a.array)
    return a
}

// Merge two arrays.
//
// 合并两个数组.
func (a *Array) Merge(array *Array) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    if a != array {
        array.mu.RLock()
        defer array.mu.RUnlock()
    }
    a.array = append(a.array, array.array...)
    return a
}

// Fills an array with num entries of the value of the value parameter,
// keys starting at the start_index parameter.
//
// 用value参数的值将数组填充num个条目，位置由startIndex参数指定的开始。
func (a *Array) Fill(startIndex int, num int, value interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    if startIndex < 0 {
        startIndex = 0
    }
    for i := startIndex; i < startIndex + num; i++ {
        if i > len(a.array) - 1 {
            a.array = append(a.array, value)
        } else {
            a.array[i] = value
        }
    }
    return a
}

// Chunks an array into arrays with size elements.
// The last chunk may contain less than size elements.
//
// 将一个数组分割成多个数组，其中每个数组的单元数目由size决定。最后一个数组的单元数目可能会少于size个。
func (a *Array) Chunk(size int) [][]interface{} {
    if size < 1 {
        return nil
    }
    a.mu.RLock()
    defer a.mu.RUnlock()
    length := len(a.array)
    chunks := int(math.Ceil(float64(length) / float64(size)))
    var n [][]interface{}
    for i, end := 0, 0; chunks > 0; chunks-- {
        end = (i + 1) * size
        if end > length {
            end = length
        }
        n = append(n, a.array[i*size : end])
        i++
    }
    return n
}

// Pad array to the specified length with a value.
// If size is positive then the array is padded on the right,
// if it's negative then on the left.
// If the absolute value of size is less than or equal to the length of the array
// then no padding takes place.
//
// 返回数组的一个拷贝，并用value将其填补到size指定的长度。
// 如果size为正数，则填补到数组的右侧，如果为负数则从左侧开始填补。
// 如果size的绝对值小于或等于数组的长度则没有任何填补。
func (a *Array) Pad(size int, val interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    if size == 0 || (size > 0 && size < len(a.array)) || (size < 0 && size > -len(a.array)) {
        return a
    }
    n := size
    if size < 0 {
        n = -size
    }
    n   -= len(a.array)
    tmp := make([]interface{}, n)
    for i := 0; i < n; i++ {
        tmp[i] = val
    }
    if size > 0 {
        a.array = append(a.array, tmp...)
    } else {
        a.array = append(tmp, a.array...)
    }
    return a
}

// Extract a slice of the array(If in concurrent safe usage, it returns a copy of the slice; else a pointer).
// It returns the sequence of elements from the array array as specified by the offset and length parameters.
//
// 返回根据offset和size参数所指定的数组中的一段序列。
func (a *Array) SubSlice(offset, size int) []interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if offset > len(a.array) {
        return nil
    }
    if offset + size > len(a.array) {
        size = len(a.array) - offset
    }
    if a.mu.IsSafe() {
        s := make([]interface{}, size)
        copy(s, a.array[offset:])
        return s
    } else {
        return a.array[offset:]
    }
}

// Picks one or more random entries out of an array(a copy),
// and returns the key (or keys) of the random entries.
//
// 从数组中随机取出size个元素项，构成slice返回。
func (a *Array) Rand(size int) []interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if size > len(a.array) {
        size = len(a.array)
    }
    n := make([]interface{}, size)
    for i, v := range grand.Perm(len(a.array)) {
        n[i] = a.array[v]
        if i == size - 1 {
            break
        }
    }
    return n
}

// Randomly shuffles the array.
//
// 随机打乱当前数组。
func (a *Array) Shuffle() *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    for i, v := range grand.Perm(len(a.array)) {
        a.array[i], a.array[v] = a.array[v], a.array[i]
    }
    return a
}

// Make array with elements in reverse order.
//
// 将当前数组反转。
func (a *Array) Reverse() *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    for i, j := 0, len(a.array) - 1; i < j; i, j = i + 1, j - 1 {
        a.array[i], a.array[j] = a.array[j], a.array[i]
    }
    return a
}

// Join array elements with a string.
//
// 使用glue字符串串连当前数组的元素项，构造成新的字符串返回。
func (a *Array) Join(glue string) string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return strings.Join(gconv.Strings(a.array), glue)
}

// Counts all the values of an array.
//
// 统计数组中所有的值出现的次数.
func (a *Array) CountValues() map[interface{}]int {
    m := make(map[interface{}]int)
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, v := range a.array {
        m[v]++
    }
    return m
}