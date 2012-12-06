// Copyright 2012 Sonia Keys
// License MIT: http://www.opensource.org/licenses/MIT

package main

// Knuth's data object
type x struct {
	c          *y
	u, d, l, r *x
	// except x0 is not Knuth's.  it's pointer to first constraint in row,
	// so that the sudoku string can be constructed from the dlx solution.
	x0 *x
}

// Knuth's column object
type y struct {
	x
	s int // size
	n int // name
}

// an object to hold the matrix and solution
type DLX struct {
	ch []y  // all column headers
	h  *y   // ch[0], the root node
	o  []*x // solution
}

// constructor creates the column headers but no rows.
func New(nCols int) *DLX {
	ch := make([]y, nCols+1)
	h := &ch[0]
	d := &DLX{ch, h, nil}
	h.c = h
	h.l = &ch[nCols].x
	ch[nCols].r = &h.x
	nh := ch[1:]
	for i := range ch[1:] {
		hi := &nh[i]
		ix := &hi.x
		hi.n = i
		hi.c = hi
		hi.u = ix
		hi.d = ix
		hi.l = &h.x
		h.r = ix
		h = hi
	}
	return d
}

// rows define constraints
func (d *DLX) AddRow(nr []int) {
	if len(nr) == 0 {
		return
	}
	r := make([]x, len(nr))
	x0 := &r[0]
	for x, j := range nr {
		ch := &d.ch[j+1]
		ch.s++
		np := &r[x]
		np.c = ch
		np.u = ch.u
		np.d = &ch.x
		np.l = &r[(x+len(r)-1)%len(r)]
		np.r = &r[(x+1)%len(r)]
		np.u.d, np.d.u, np.l.r, np.r.l = np, np, np, np
		np.x0 = x0
	}
}

// the dlx algorithm
func (d *DLX) Search() bool {
	h := d.h
	j := h.r.c
	if j == h {
		return true
	}
	c := j
	for minS := j.s; ; {
		j = j.r.c
		if j == h {
			break
		}
		if j.s < minS {
			c, minS = j, j.s
		}
	}

	cover(c)
	k := len(d.o)
	d.o = append(d.o, nil)
	for r := c.d; r != &c.x; r = r.d {
		d.o[k] = r
		for j := r.r; j != r; j = j.r {
			cover(j.c)
		}
		if d.Search() {
			return true
		}
		r = d.o[k]
		c = r.c
		for j := r.l; j != r; j = j.l {
			uncover(j.c)
		}
	}
	d.o = d.o[:len(d.o)-1]
	uncover(c)
	return false
}

func cover(c *y) {
	c.r.l, c.l.r = c.l, c.r
	for i := c.d; i != &c.x; i = i.d {
		for j := i.r; j != i; j = j.r {
			j.d.u, j.u.d = j.u, j.d
			j.c.s--
		}
	}
}

func uncover(c *y) {
	for i := c.u; i != &c.x; i = i.u {
		for j := i.l; j != i; j = j.l {
			j.c.s++
			j.d.u, j.u.d = j, j
		}
	}
	c.r.l, c.l.r = &c.x, &c.x
}
