package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type History struct {
	IsSale bool
	Price  float64
	Time   time.Time
}

type Indexes struct {
	Age, Name, Proof, Sku, Retail, Sale, Size, Type int
	IsSale                                          bool
}

type Listing struct {
	LastUpdate time.Time
	Products   map[Sku]*Product
}

type Product struct {
	Sku    Sku
	Name   string
	Type   string
	Prices map[int64]History
	delete bool
}

type Sku int

var idxs = map[string]*Indexes{
	"discounts": &Indexes{Type: 0, Sku: 1, Name: 2, Retail: 3, Sale: 5, IsSale: true},
	"prices":    &Indexes{Type: 0, Sku: 1, Name: 2, Size: 3, Age: 4, Proof: 5, Retail: 6},
}

func NewProduct(sku Sku) *Product {
	return &Product{
		Sku:    sku,
		Prices: make(map[int64]History),
	}
}

func (p *Product) Update(data []string, idx *Indexes, t time.Time) (err error) {
	p.Name = data[idx.Name]
	p.Type = Type(data[idx.Type])
	p.delete = false
	h := History{
		IsSale: idx.IsSale,
		Time:   t,
	}
	var str string
	if idx.IsSale {
		str = data[idx.Sale]
	} else {
		str = data[idx.Retail][1:]
	}
	_, err = fmt.Sscanf(str, "%f", &h.Price)
	p.Prices[t.Unix()] = h
	return
}

// func (s Sku) String() string {
// 	return fmt.Sprintf("%06d", s)
// }

func initListing(dir string, l *Listing) (err error) {
	l.Products = make(map[Sku]*Product, 4096)

	// Read in the entire product list, then build from there
	prices, err := filepath.Glob(filepath.Join(dir, "prices", "*.csv"))
	sort.Strings(prices)
	discounts, err := filepath.Glob(filepath.Join(dir, "discounts", "*.csv"))
	sort.Strings(discounts)
	prices = append(prices, discounts...)
	for _, path := range prices {
		if err = readListing(path, l); err != nil {
			return
		}
	}

	// Remove discontinued products
	for s := range l.Products {
		if p := l.Products[s]; p.delete || p.Type != "Scotch" {
			delete(l.Products, s)
		}
	}

	for _, p := range l.Products {
		fmt.Printf("%06d ", p.Sku)
		h := p.Prices[l.LastUpdate.Unix()]
		//fmt.Printf("%6.2f", h.Price)
		for _, h := range p.Prices {
			fmt.Printf("%7.2f", h.Price)
		}
		if h.IsSale {
			fmt.Print("* ")
		}
		fmt.Println()
	}

	fmt.Printf("\n\n%d items\n", len(l.Products))

	return
}

func readListing(path string, l *Listing) (err error) {
	var (
		i   int
		sku Sku
		idx = idxs[filepath.Base(filepath.Dir(path))]
	)

	// Time Stuff
	t, err := time.Parse("2006-01.csv", filepath.Base(path))
	if err != nil {
		return
	}
	if t.After(l.LastUpdate) {
		l.LastUpdate = t
	}

	// Gain access to data
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	// Parse parse parse
	CSV := csv.NewReader(f)
	if _, err = CSV.Read(); err != nil {
		return
	}

	// Flip deletion flag when processing price lists
	if !idx.IsSale {
		for s := range l.Products {
			l.Products[s].delete = true
		}
	}

	for row, err := CSV.Read(); err == nil; row, err = CSV.Read() {
		i, err = strconv.Atoi(row[idx.Sku])
		sku = Sku(i)
		p, ok := l.Products[sku]
		if !ok {
			p = NewProduct(sku)
			l.Products[sku] = p
		}
		if err := p.Update(row, idx, t); err != nil {
			return err
		}
	}
	return
}
