package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	if len(memory) == 0 || len(pointers) == 0 {
		return
	}
	// Находим все занятые позиции в памяти
	occupiedPositions := findOccupiedPositions(memory, pointers)
	// Уплотняем память, перемещая занятые байты в начало
	moveMap := compactMemory(memory, occupiedPositions)
	// Очищаем неиспользуемую часть памяти
	clearUnusedMemory(memory, len(moveMap))
	// Обновляем указатели согласно новым позициям
	updatePointers(memory, pointers, moveMap)
}

func getBaseAddress(memory []byte) uintptr {
	if len(memory) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&memory[0]))
}

func getPointerOffset(ptr unsafe.Pointer, baseAddr uintptr) int {
	if ptr == nil {
		return -1
	}
	offset := uintptr(ptr) - baseAddr
	return int(offset)
}

func isValidOffset(offset int, memorySize int) bool {
	return offset >= 0 && offset < memorySize
}

func findOccupiedPositions(memory []byte, pointers []unsafe.Pointer) map[int]bool {
	occupiedPositions := make(map[int]bool)
	baseAddr := getBaseAddress(memory)

	for _, ptr := range pointers {
		offset := getPointerOffset(ptr, baseAddr)
		if isValidOffset(offset, len(memory)) {
			occupiedPositions[offset] = true
		}
	}

	return occupiedPositions
}

func compactMemory(memory []byte, occupiedPositions map[int]bool) map[int]int {
	moveMap := make(map[int]int)
	writePos := 0

	for readPos := 0; readPos < len(memory); readPos++ {
		if occupiedPositions[readPos] {
			if readPos != writePos {
				memory[writePos] = memory[readPos]
				memory[readPos] = 0x00
			}
			moveMap[readPos] = writePos
			writePos++
		}
	}

	return moveMap
}

func clearUnusedMemory(memory []byte, usedSize int) {
	for i := usedSize; i < len(memory); i++ {
		memory[i] = 0x00
	}
}

func updatePointers(memory []byte, pointers []unsafe.Pointer, moveMap map[int]int) {
	baseAddr := getBaseAddress(memory)

	for i, ptr := range pointers {
		if ptr == nil {
			continue
		}

		oldOffset := getPointerOffset(ptr, baseAddr)
		if !isValidOffset(oldOffset, len(memory)) {
			continue
		}

		if newPos, exists := moveMap[oldOffset]; exists {
			newPtr := unsafe.Pointer(baseAddr + uintptr(newPos))
			pointers[i] = newPtr
		}
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
