package state

import (
	"testing"
	"time"
)

var oldState = &State{
	Products: map[Sku]*Product{
		1: &Product{Sku: 1, Name: "Product 1", LastUpdate: time.Date(2013, time.October, 1, 0, 0, 0, 0, time.Local)},
		2: &Product{Sku: 2, Name: "Product 2", LastUpdate: time.Date(2013, time.October, 1, 0, 0, 0, 0, time.Local)},
		3: &Product{Sku: 3, Name: "Product 3", LastUpdate: time.Date(2013, time.October, 1, 0, 0, 0, 0, time.Local)},
	},
}
var updateState = &State{
	Products: map[Sku]*Product{
		1: &Product{Sku: 1, Name: "Product 1", LastUpdate: time.Date(2013, time.October, 1, 0, 0, 0, 0, time.Local)},
		2: &Product{Sku: 2, Name: "Updated 2", LastUpdate: time.Date(2013, time.December, 1, 0, 0, 0, 0, time.Local)},
		3: &Product{Sku: 3, Name: "Updated 3", LastUpdate: time.Date(2013, time.November, 1, 0, 0, 0, 0, time.Local)},
	},
}
var addDelState = &State{
	Products: map[Sku]*Product{
		1: &Product{Sku: 1, Name: "Product 1", LastUpdate: time.Date(2013, time.October, 1, 0, 0, 0, 0, time.Local)},
		3: &Product{Sku: 3, Name: "Product 3", LastUpdate: time.Date(2013, time.October, 1, 0, 0, 0, 0, time.Local)},
		4: &Product{Sku: 4, Name: "Product 4", LastUpdate: time.Date(2013, time.November, 1, 0, 0, 0, 0, time.Local)},
	},
}

func TestAdd(t *testing.T) {
	s := New()
	for _, p := range oldState.Products {
		s.Add(p)
	}
	if len(s.Products) != len(oldState.Products) {
		t.Fatal("Lengths did not match")
	}
}

func TestDiffUpdate(t *testing.T) {
	diff := updateState.Diff(oldState)
	t.Logf("%+v", diff)
}

func TestDiffAddDel(t *testing.T) {
	diff := addDelState.Diff(oldState)
	t.Logf("%+v", diff)
}
