package types

import (
	"gopkg.in/yaml.v3"
)

type YAMLNodeMetadata[T any] struct {
	value  *T
	line   *int
	column *int
}

func (y *YAMLNodeMetadata[T]) GetValue() T {
	if y == nil || y.value == nil {
		return *new(T)
	}

	return *y.value
}

func (y *YAMLNodeMetadata[T]) GetLine() int {
	if y == nil || y.line == nil {
		return -1
	}

	return *y.line
}

func (y *YAMLNodeMetadata[T]) GetColumn() int {
	if y == nil || y.column == nil {
		return -1
	}

	return *y.column
}

func (y *YAMLNodeMetadata[T]) UnmarshalYAML(node *yaml.Node) error {
	y.value = new(T)

	err := node.Decode(y.value)

	if err != nil {
		return err
	}

	y.line = &node.Line
	y.column = &node.Column

	return nil
}
