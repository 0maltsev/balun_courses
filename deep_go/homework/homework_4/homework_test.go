package main

import (
    "reflect"
    "testing"

    "github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node struct {
    key   int
    value int
    left  *node
    right *node
}

type OrderedMap struct {
    root *node
    size int
}

func NewOrderedMap() OrderedMap {
    return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
    if m.root == nil {
        m.root = &node{key: key, value: value}
        m.size++
        return
    }
    
    if m.insertNode(m.root, key, value) {
        m.size++
    }
}

func (m *OrderedMap) insertNode(n *node, key, value int) bool {
    if key < n.key {
        if n.left == nil {
            n.left = &node{key: key, value: value}
            return true
        }
        return m.insertNode(n.left, key, value)
    } else if key > n.key {
        if n.right == nil {
            n.right = &node{key: key, value: value}
            return true
        }
        return m.insertNode(n.right, key, value)
    } else {
        n.value = value
        return false
    }
}

func (m *OrderedMap) Erase(key int) {
	var deleted bool
    if m.root == nil {
        return
    }

    m.root, deleted = m.deleteNode(m.root, key)
    if deleted {
        m.size--
    }
}

func (m *OrderedMap) deleteNode(n *node, key int) (*node, bool) {
    var deleted bool  
	if n == nil {
        return nil, false
    }  
    if key < n.key {
        n.left, deleted = m.deleteNode(n.left, key)
    } else if key > n.key {
        n.right, deleted = m.deleteNode(n.right, key)
    } else {
        deleted = true
        
        if n.left == nil && n.right == nil {
            return nil, true
        }
        
        if n.left == nil {
            return n.right, true
        }
        if n.right == nil {
            return n.left, true
        }

        successor := m.findMin(n.right)

        n.key = successor.key
        n.value = successor.value

        n.right, _ = m.deleteNode(n.right, successor.key)
    }
    
    return n, deleted
}

func (m *OrderedMap) findMin(n *node) *node {
    for n.left != nil {
        n = n.left
    }
    return n
}

func (m *OrderedMap) Contains(key int) bool {
    return m.findNode(m.root, key) != nil
}

func (m *OrderedMap) findNode(n *node, key int) *node {
    if n == nil {
        return nil
    }
    
    if key < n.key {
        return m.findNode(n.left, key)
    } else if key > n.key {
        return m.findNode(n.right, key)
    } else {
        return n
    }
}

func (m *OrderedMap) Size() int {
    return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
    m.inOrderTraversal(m.root, action)
}

func (m *OrderedMap) inOrderTraversal(n *node, action func(int, int)) {
    if n == nil {
        return
    }
    m.inOrderTraversal(n.left, action)
    action(n.key, n.value)
    m.inOrderTraversal(n.right, action)
}

func TestCircularQueue(t *testing.T) {
    data := NewOrderedMap()
    assert.Zero(t, data.Size())

    data.Insert(10, 10)
    data.Insert(5, 5)
    data.Insert(15, 15)
    data.Insert(2, 2)
    data.Insert(4, 4)
    data.Insert(12, 12)
    data.Insert(14, 14)

    assert.Equal(t, 7, data.Size())
    assert.True(t, data.Contains(4))
    assert.True(t, data.Contains(12))
    assert.False(t, data.Contains(3))
    assert.False(t, data.Contains(13))

    var keys []int
    expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
    data.ForEach(func(key, _ int) {
        keys = append(keys, key)
    })

    assert.True(t, reflect.DeepEqual(expectedKeys, keys))

    data.Erase(15)
    data.Erase(14)
    data.Erase(2)

    assert.Equal(t, 4, data.Size())
    assert.True(t, data.Contains(4))
    assert.True(t, data.Contains(12))
    assert.False(t, data.Contains(2))
    assert.False(t, data.Contains(14))

    keys = nil
    expectedKeys = []int{4, 5, 10, 12}
    data.ForEach(func(key, _ int) {
        keys = append(keys, key)
    })

    assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}