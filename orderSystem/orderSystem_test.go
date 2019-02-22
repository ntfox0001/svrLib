package orderSystem_test

import "testing"

func TestMap1(t *testing.T) {
	m1 := make(map[string]string)
	m2 := m1
	print(m1, "\n", m2)
}
