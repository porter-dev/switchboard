package types

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type YAMLNodeMetadata struct {
	Line   int
	Column int
}

// YAMLNode is intended to be used as a container data structure
// whenever the underlying data type's internal YAML node or its
// properties like the line and column need to be known
type YAMLNode[T any] struct {
	internalNode *yaml.Node
	value        *T
}

func (y *YAMLNode[T]) GetValue() T {
	if y == nil || y.value == nil {
		return *new(T)
	}

	return *y.value
}

func (y *YAMLNode[T]) GetYAMLNodeMetadata() *YAMLNodeMetadata {
	if y == nil || y.internalNode == nil {
		return nil
	}

	return &YAMLNodeMetadata{
		Line:   y.internalNode.Line,
		Column: y.internalNode.Column,
	}
}

func (y *YAMLNode[T]) GetLine() int {
	if y == nil || y.internalNode == nil {
		return -1
	}

	return y.internalNode.Line
}

func (y *YAMLNode[T]) GetColumn() int {
	if y == nil || y.internalNode == nil {
		return -1
	}

	return y.internalNode.Column
}

func (y *YAMLNode[T]) GetRawYAMLNode() *yaml.Node {
	if y == nil || y.internalNode == nil {
		return nil
	}

	return y.internalNode
}

func (y *YAMLNode[T]) UnmarshalYAML(node *yaml.Node) error {
	if y == nil {
		return fmt.Errorf("cannot unmarshal from nil pointer")
	}

	y.value = new(T)

	y.internalNode = node

	err := node.Decode(y.value)

	if err != nil {
		return err
	}

	return nil
}
