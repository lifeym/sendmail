package genericlist

import (
	"errors"
	"slices"
)

var ErrListValueNotFound = errors.New("value not found")

type GenericList[T comparable] []T

func NewGenericList[T comparable]() GenericList[T] {
	return []T{}
}

func (l GenericList[T]) Append(value T) {
	l = append(l, value)
}

func (l GenericList[T]) ValueByIndex(index int) (T, error) {
	var value T
	if index >= len(l) || index < 0 {
		return value, errors.New("index out of Range")
	}

	return l[index], nil
}

func (l GenericList[T]) RemoveByIndex(index int) (T, error) {
	var value T
	if index >= len(l) || index < 0 {
		return value, errors.New("index out of Range")
	}

	for i, data := range l {
		if i == index {
			value = data
			l = slices.Concat(l[:i], l[i+1:])
		}
	}

	return value, nil
}

func (l GenericList[T]) RemoveByValue(data T) (T, error) {
	var value T
	for i, info := range l {
		if data == info {
			value = info
			l = slices.Concat(l[:i], l[i+1:])
			return value, nil
		}
	}

	return value, ErrListValueNotFound
}
