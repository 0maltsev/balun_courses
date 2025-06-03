package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	visited := make(map[uintptr]struct{})
	
	for _, stack := range stacks {
		processStack(stack, visited)
	}

	return convertMapToSlice(visited)
}

func processStack(stack []uintptr, visited map[uintptr]struct{}) {
	for _, ptr := range stack {
		if isValidPointer(ptr) {
			dfs(ptr, visited)
		}
	}
}

func dfs(ptr uintptr, visited map[uintptr]struct{}) {
	if !isValidPointer(ptr) {
		return
	}
	if _, ok := visited[ptr]; ok {
		return
	}
	visited[ptr] = struct{}{}
	value := readPointerValue(ptr)
	if isValidPointer(value) {
		dfs(value, visited)
	}
}

func isValidPointer(ptr uintptr) bool {
	return ptr != 0
}

func readPointerValue(ptr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr))
}

func convertMapToSlice(visited map[uintptr]struct{}) []uintptr {
	result := make([]uintptr, len(visited))
	i := 0
	for ptr := range visited {
		result[i] = ptr
		i++
	}
	return result
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}

	pointers := Trace(stacks)
	assert.ElementsMatch(t, pointers, expectedPointers)
}