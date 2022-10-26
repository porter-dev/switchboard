package types

import (
	"gopkg.in/yaml.v3"
)

type YAMLNode[T any] struct {
	value  *T
	line   *int
	column *int
}

func (y *YAMLNode[T]) GetValue() T {
	if y == nil || y.value == nil {
		return *new(T)
	}

	return *y.value
}

func (y *YAMLNode[T]) GetLine() int {
	if y == nil || y.line == nil {
		return -1
	}

	return *y.line
}

func (y *YAMLNode[T]) GetColumn() int {
	if y == nil || y.column == nil {
		return -1
	}

	return *y.column
}

func (y *YAMLNode[T]) UnmarshalYAML(node *yaml.Node) error {
	y.value = new(T)

	err := node.Decode(y.value)

	if err != nil {
		return err
	}

	y.line = &node.Line
	y.column = &node.Column

	return nil
}

type YAMLMapInterfaceNode struct {
	compositeData map[*YAMLNode[string]]*YAMLNode[interface{}]
	raw           map[string]interface{}
}

func (y *YAMLMapInterfaceNode) UnmarshalYAML(node *yaml.Node) error {
	y.compositeData = make(map[*YAMLNode[string]]*YAMLNode[any])
	y.raw = make(map[string]interface{})

	err := node.Decode(&y.raw)

	if err != nil {
		return err
	}

	err = node.Decode(&y.compositeData)

	if err != nil {
		return err
	}

	return nil
}
