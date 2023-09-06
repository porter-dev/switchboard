package parser

import "testing"

func Test_ParseRawBytes(t *testing.T) {
	exampleYAML := []byte(`version: v2

variables:
- name: psqlPassword
  once: true
  type: string
  random: true`)

	_, err := ParseRawBytes(exampleYAML)

	if err != nil {
		t.Error(err)
	}
}
