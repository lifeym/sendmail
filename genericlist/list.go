package genericlist

import (
	"errors"
	"slices"
)

var ErrListValueNotFound = errors.New("value not found")

type GenericList[T comparable] struct {
	Data []T
}

func NewGenericList[T comparable]() *GenericList[T] {
	return &GenericList[T]{
		Data: []T{},
	}
}

func (l *GenericList[T]) Append(value T) {
	l.Data = append(l.Data, value)
}

func (l *GenericList[T]) ValueByIndex(index int) (T, error) {
	var value T
	if index >= len(l.Data) || index < 0 {
		return value, errors.New("index out of Range")
	}

	return l.Data[index], nil
}

func (l *GenericList[T]) RemoveByIndex(index int) (T, error) {
	var value T
	if index >= len(l.Data) || index < 0 {
		return value, errors.New("index out of Range")
	}

	for i, data := range l.Data {
		if i == index {
			value = data
			l.Data = slices.Concat(l.Data[:i], l.Data[i+1:])
		}
	}

	return value, nil
}

func (l *GenericList[T]) RemoveByValue(data T) (T, error) {
	var value T
	for i, info := range l.Data {
		if data == info {
			value = info
			l.Data = slices.Concat(l.Data[:i], l.Data[i+1:])
			return value, nil
		}
	}

	return value, ErrListValueNotFound
}

func (l *GenericList[T]) Len() int {
	return len(l.Data)
}
