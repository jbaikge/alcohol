package state

import (
	"time"
)

type Sku int

type Diff struct {
	Put []*Product // Adds and updates
	Del []Sku      // Skus
}

type Product struct {
	Sku        Sku
	Name       string
	Type       string
	OnSale     bool
	Price      float64
	History    map[int64]float64
	LastUpdate time.Time
}

type State struct {
	Products map[Sku]*Product
}

func New() *State {
	return &State{
		Products: make(map[Sku]*Product),
	}
}

func (s *State) Add(p *Product) {
	
	s.Products[p.Sku] = p
}

func (s *State) Diff(old *State) (d *Diff) {
	d = &Diff{
		Put: make([]*Product, 0, len(s.Products)),
		Del: make([]Sku, 0, 8),
	}
	// Determine which SKUs are no longer sold
	for sku := range old.Products {
		if _, ok := s.Products[sku]; !ok {
			d.Del = append(d.Del, sku)
		}
	}
	// Determine which SKUs to update
	for sku, p := range s.Products {
		if op, ok := old.Products[sku]; !ok || p.LastUpdate.After(op.LastUpdate) {
			d.Put = append(d.Put, p)
		}
	}
	return
}
