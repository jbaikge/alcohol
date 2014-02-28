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

type Price struct {
	IsSale bool
	Price  float64
	Retail float64
	Time   time.Time
}

type Indexes struct {
	Age, Name, Proof, Sku, Retail, Sale, Size, Type int
	IsSale                                          bool
}

type Listing struct {
	LastUpdate time.Time
	SaleOnly   bool
	Products   []*Product
}

type Product struct {
	Sku     Sku
	Name    string
	Type    string
	Volume  string
	History []Price
	delete  bool
}

type Sku int

var (
	months    = 13
	startDate = time.Date(time.Now().Year()-1, time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	idxs      = map[string]*Indexes{
		"discounts": &Indexes{Type: 0, Sku: 1, Name: 2, Retail: 3, Sale: 5, IsSale: true},
		"prices":    &Indexes{Type: 0, Sku: 1, Name: 2, Size: 3, Age: 4, Proof: 5, Retail: 6},
	}
	_ sort.Interface = new(Listing)
)

func (l *Listing) CountType(t string) (count int) {
	for i := range l.Products {
		if l.Products[i].Type == t {
			count++
		}
	}
	return
}

func (l *Listing) Len() int {
	return len(l.Products)
}

func (l *Listing) Less(i, j int) bool {
	var a, b = l.Products[i], l.Products[j]
	// Compare types first, if equal, process names
	if a.Type != b.Type {
		return a.Type < b.Type
	}
	// If names equal, sort by price
	if a.Name == b.Name {
		return a.Price() < b.Price()
	}
	// Compare names
	var ii, ji int // starting string indexes
	if a.Name[0:3] == "The" {
		ii = 4
	}
	if b.Name[0:3] == "The" {
		ji = 4
	}
	return a.Name[ii:] < b.Name[ji:]
}

func (l *Listing) Swap(i, j int) {
	l.Products[i], l.Products[j] = l.Products[j], l.Products[i]
}

func (l *Listing) Search(sku Sku) (p *Product, found bool) {
	for i := range l.Products {
		if l.Products[i].Sku == sku {
			return l.Products[i], true
		}
	}
	return
}

func NewProduct(sku Sku) *Product {
	return &Product{
		Sku:     sku,
		History: make([]Price, months),
	}
}

func (p *Product) Update(data []string, idx *Indexes, t time.Time) (err error) {
	p.Name = data[idx.Name]
	// Clean name up
	for _, suffix := range []string{"750ml", "1.75L"} {
		if len(p.Name) > len(suffix) && p.Name[len(p.Name)-len(suffix):] == suffix {
			p.Name = p.Name[:len(p.Name)-len(suffix)-1]
		}
	}

	p.Type = Type(data[idx.Type])

	if idx.Size > 0 {
		p.Volume = data[idx.Size]
	}

	p.delete = false
	h := Price{
		IsSale: idx.IsSale,
		Time:   t,
	}
	var str string
	if idx.IsSale {
		str = data[idx.Sale]
		h.Retail, err = strconv.ParseFloat(data[idx.Retail], 64)
	} else {
		str = data[idx.Retail][1:]
	}
	h.Price, err = strconv.ParseFloat(str, 64)
	p.History[historyIndex(t)] = h
	return
}

func (p *Product) OnSale() bool {
	return p.History[len(p.History)-1].IsSale
}

func (p *Product) Price() float64 {
	return p.History[len(p.History)-1].Price
}

func (p *Product) Retail() float64 {
	if r := p.History[len(p.History)-1].Retail; r > 0 {
		return r
	}
	return p.Price()
}

func (p *Product) FmtPrice() string {
	return fmt.Sprintf("$%0.2f", p.Price())
}

func (p *Product) FmtRetail() string {
	return fmt.Sprintf("$%0.2f", p.Retail())
}

func (s Sku) String() string {
	return fmt.Sprintf("%06d", s)
}

func initListing(dir string, l *Listing) (err error) {
	l.Products = make([]*Product, 0, 4096)

	// Read in the entire product list, then build from there
	prices, err := filepath.Glob(filepath.Join(dir, "prices", "????-??.csv"))
	sort.Strings(prices)
	discounts, err := filepath.Glob(filepath.Join(dir, "discounts", "????-??.csv"))
	sort.Strings(discounts)
	prices = append(prices, discounts...)
	for _, path := range prices {
		if err = readListing(path, l); err != nil {
			return
		}
	}

	// Remove discontinued products
	for i := 0; i < l.Len(); {
		if p := l.Products[i]; p.delete || p.Type == IgnoreType {
			l.Products = append(l.Products[0:i], l.Products[i+1:]...)
			continue
		}
		i++
	}

	sort.Sort(l)

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
	if t.Before(startDate) {
		// Bail if outside the date range
		return
	}
	if t.After(startDate.AddDate(0, months-1, 0)) {
		return
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

	// Add new products, update existing
	for row, err := CSV.Read(); err == nil; row, err = CSV.Read() {
		i, err = strconv.Atoi(row[idx.Sku])
		sku = Sku(i)
		p, ok := l.Search(sku)
		if !ok {
			p = NewProduct(sku)
			l.Products = append(l.Products, p)
		}
		if err := p.Update(row, idx, t); err != nil {
			return err
		}
	}

	// Set prices on the tail equal to the first price found when going back
	pIdx := historyIndex(t)
	for _, p := range l.Products {
		if p.History[pIdx].Price > 0.00 {
			continue
		}
		// Can't do anything about this situation..
		if pIdx == 0 {
			continue
		}
		p.History[pIdx].Price = p.History[pIdx-1].Price
	}
	return
}

func initListingSale(sale, all *Listing) {
	sale.LastUpdate = all.LastUpdate
	sale.SaleOnly = true
	sale.Products = make([]*Product, 0, 512)
	for _, p := range all.Products {
		if p.OnSale() {
			sale.Products = append(sale.Products, p)
		}
	}
}

func historyIndex(t time.Time) int {
	return (t.Year()-startDate.Year())*12 + int(t.Month()-startDate.Month())
}
