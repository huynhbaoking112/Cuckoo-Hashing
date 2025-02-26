package main

import (
	"fmt"
	"hash/fnv"
)

// CuckooHash đại diện cho bảng băm Cuckoo
type CuckooHash struct {
	table1 []int // Bảng thứ nhất (dùng h1)
	table2 []int // Bảng thứ hai (dùng h2)
	size   int   // Kích thước mỗi bảng
	count  int   // Số lượng phần tử hiện tại
}

// NewCuckooHash tạo một bảng Cuckoo Hash mới với kích thước cho trước
func NewCuckooHash(size int) *CuckooHash {
	return &CuckooHash{
		table1: make([]int, size),
		table2: make([]int, size),
		size:   size,
		count:  0,
	}
}

// hash1 là hàm băm thứ nhất
func (ch *CuckooHash) hash1(key int) int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d", key)))
	return int(h.Sum32()) % ch.size
}

// hash2 là hàm băm thứ hai
func (ch *CuckooHash) hash2(key int) int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d%d", key, key)))
	return int(h.Sum32()) % ch.size
}

// Insert chèn một khóa vào bảng
func (ch *CuckooHash) Insert(key int) bool {
	if ch.contains(key) {
		return true // Nếu khóa đã tồn tại, không chèn
	}

	if ch.count >= ch.size/2 { // Load factor > 0.5, rehash
		ch.rehash()
	}

	return ch.insertHelper(key, 0)
}

// insertHelper xử lý chèn với giới hạn đẩy tối đa
func (ch *CuckooHash) insertHelper(key int, depth int) bool {
	const maxDepth = 10 // Giới hạn số lần đẩy để tránh vòng lặp vô hạn

	if depth >= maxDepth { // Nếu vượt quá giới hạn, rehash
		ch.rehash()
		return ch.insertHelper(key, 0)
	}

	// Thử đặt vào table1
	pos1 := ch.hash1(key)
	if ch.table1[pos1] == 0 { // Ô trống
		ch.table1[pos1] = key
		ch.count++
		return true
	}

	// Đẩy giá trị hiện tại ra
	oldKey := ch.table1[pos1]
	ch.table1[pos1] = key

	// Thử đặt oldKey vào table2
	pos2 := ch.hash2(oldKey)
	if ch.table2[pos2] == 0 { // Ô trống
		ch.table2[pos2] = oldKey
		ch.count++
		return true
	}

	// Đẩy giá trị trong table2 ra
	olderKey := ch.table2[pos2]
	ch.table2[pos2] = oldKey

	// Tiếp tục đẩy olderKey
	return ch.insertHelper(olderKey, depth+1)
}

// Lookup kiểm tra xem khóa có trong bảng không
func (ch *CuckooHash) Lookup(key int) bool {
	pos1 := ch.hash1(key)
	if ch.table1[pos1] == key {
		return true
	}

	pos2 := ch.hash2(key)
	if ch.table2[pos2] == key {
		return true
	}

	return false
}

// contains kiểm tra xem khóa đã tồn tại chưa (dùng trong Insert)
func (ch *CuckooHash) contains(key int) bool {
	return ch.Lookup(key)
}

// rehash tăng kích thước bảng và chèn lại tất cả khóa
func (ch *CuckooHash) rehash() {
	oldTable1 := ch.table1
	oldTable2 := ch.table2
	ch.size *= 2 // Tăng gấp đôi kích thước
	ch.table1 = make([]int, ch.size)
	ch.table2 = make([]int, ch.size)
	ch.count = 0

	// Chèn lại các giá trị từ bảng cũ
	for _, key := range oldTable1 {
		if key != 0 {
			ch.Insert(key)
		}
	}
	for _, key := range oldTable2 {
		if key != 0 {
			ch.Insert(key)
		}
	}
}

// printTable in trạng thái bảng để kiểm tra
func (ch *CuckooHash) printTable() {
	fmt.Printf("Table 1: %v\n", ch.table1)
	fmt.Printf("Table 2: %v\n", ch.table2)
}

func main() {
	// Tạo bảng Cuckoo Hash với kích thước ban đầu là 6
	ch := NewCuckooHash(6)

	// Chèn các khóa
	fmt.Println("Chèn 7:")
	ch.Insert(7)
	ch.printTable()

	fmt.Println("Chèn 13:")
	ch.Insert(13)
	ch.printTable()

	fmt.Println("Chèn 9:")
	ch.Insert(9)
	ch.printTable()

	fmt.Println("Chèn 19:")
	ch.Insert(19)
	ch.printTable()

	// Tra cứu
	fmt.Println("Tìm 13:", ch.Lookup(13)) // true
	fmt.Println("Tìm 20:", ch.Lookup(20)) // false
}
