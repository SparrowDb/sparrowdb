package backup

import "testing"

func TestDo(T *testing.T) {
	CreateSnapshot("../data/teste3", "../data/teste3")
}
