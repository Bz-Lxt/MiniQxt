package shuffle

import "testing"

func TestOrderReproducible(t *testing.T) {
	a := Order(99, 8)
	b := Order(99, 8)
	if len(a) != 8 {
		t.Fatal(len(a))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("not stable %v %v", a, b)
		}
	}
	c := Order(100, 8)
	same := true
	for i := range a {
		if a[i] != c[i] {
			same = false
		}
	}
	if same {
		t.Fatal("different seeds should differ")
	}
}

func TestApply(t *testing.T) {
	items := []string{"a", "b", "c"}
	out := Apply(items, []int{2, 0, 1})
	if out[0] != "c" || out[1] != "a" || out[2] != "b" {
		t.Fatalf("%v", out)
	}
}
