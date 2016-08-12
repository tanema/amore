package mat

import (
	"errors"
)

// A Stack is an OpenGL-style matrix stack,
// usually used for things like scenegraphs. This allows you
// to easily maintain matrix state per call level.
type Stack []*Mat4

func NewStack() *Stack {
	return &Stack{New4()}
}

// Copies the top element and pushes it on the stack.
func (ms *Stack) Push(top *Mat4) {
	(*ms) = append(*ms, top)
}

// Removes the first element of the matrix from the stack, if there is only one element left
// there is an error.
func (ms *Stack) Pop() (*Mat4, error) {
	if len(*ms) == 1 {
		return nil, errors.New("Cannot pop from mat stack, at minimum stack length of 1")
	}
	popped := ms.Peek()
	(*ms) = (*ms)[:len(*ms)-1]
	return popped, nil
}

// Returns the top element.
func (ms *Stack) Peek() *Mat4 {
	return (*ms)[len(*ms)-1]
}
