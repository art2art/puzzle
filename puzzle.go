package puzzle

import (
	"fmt"
	"math/rand"
	"math"
)

const (
	side int = 9
	sbox int = 3
)

type Grid []int

func seq(s, e int) []int {
	seq := make([]int, e-s)
	for i, l := 0, e-s; i < l; s, i = s+1, i+1 {
		seq[i] = s
	}
	return seq
}

func (g Grid) copy() Grid {
	c := make([]int, len(g))
	copy(c, []int(g))
	return c
}

func adjacent(i int) []int {
	r, c, b := row(i), column(i), box(i)
	return union(r, c, b)
}

func union(r, c, b []int) (k []int) {
	m := make(map[int]bool, 20)
	for _, v := range append(r, append(b, c...)...) {
		m[v] = true
	}
	for key := range m { k = append(k, key) }
	return 
}

func column(e int) (c []int) {
	for i, t := 0, false; t != true && i < side; i++ {
		c = make([]int, 0)
		for j := 0; j < side; j++ {
			if e == j*side+i {
				t = true
				continue
			}
			c = append(c, j*side+i)
		}
	}
	return
}

func row(e int) (r []int) {
	for i, j := 0, side; i < side*side; i, j = j, j+side {
		if i <= e && e < j {
			for ; i < j; i++ {
				if i != e { r = append(r, i) }
			}
			return 
		}
	}
	return
}

func box(e int) (box []int) {
	for i := 0; i <= side*side; i = i+side*3 {
		for j := 0; j < 3; j++ {
			s1, e1 := i+3*j, i+3*j+3
			s2, e2 := s1+side, e1+side
			s3, e3 := s2+side, e2+side
			switch {
			case s1 <= e && e < e1, s2 <= e && e < e2, s3 <= e && e < e3:
				i1, i2, i3 := seq(s1, e1), seq(s2, e2), seq(s3, e3)
				box = append(i1, append(i2, i3...)...)
				for k, v := range box {
					if v == e { return append(box[:k], box[k+1:]...) }
				}
			}
		}
	}
	return 
}

func (grid Grid) hstar(adj []int) (len int, poss []int) {
	m := make(map[int]bool)
	for _, i := range adj {
		if val := grid[i]; val != 0 {
			m[val] = true
		}
	}
	for i := 1; i <= 9; i++ {
		if _, ok := m[i]; !ok {
			poss = append(poss, i)
			len++
		}
	}
	return
}

func (grid Grid) gstar(adj []int) (g int) {
	for _, i := range adj {
		if val := grid[i]; val == 0 { g++ }
	}
	return
}

func (grid Grid) fstar() (im int, iposs []int) {
	min := math.MaxInt32
	for i, v := range grid {
		if v != 0 {
			continue
		}
		adj := adjacent(i)
		hval, poss := grid.hstar(adj)
		if len(poss) == 0 {
			return 0, nil
		}
		fval := grid.gstar(adj) + hval
		if fval <= min {
			min, im, iposs = fval, i, poss
		}
	}
	return 
}

func next(in Grid) (out Grid, err error) {
	if in.Test() {
		return in, nil
	}
	i, poss := in.fstar()
	if poss == nil {
		return nil, fmt.Errorf(
			"Solution can't be found for this instance (%d)\n", i)
	}
	newg := in.copy()
	for _, val := range poss {
		newg[i] = val
		out, err = next(newg)
		if err == nil {
			break
		}
	}
	return
}

func (src Grid) Solve() (Grid, error) {
	return next(src)
}

func (g Grid) Test() bool {
	test := func (indices []int) bool {
		m := make(map[int]bool)
		for _, i := range indices {
			_, ok := m[i]
			switch {
			case g[i] == 0: return false
			case ok: return false
			default: m[i] = true
			}
		}
		return true
	}
	for i, j, k := 0, 0, 0; i < side; i, j = i+1, j+1 {
		if j > 2 {
			j = 0
			k = k+3*side
		}
		c := column(i)
		r := row(i*side)
		b := box(k+j*3)
		if !test(c) || !test(r) || !test(b) {
			return false
		}
	}
	return true
}

func (g Grid) verify(idx, val int) bool {
	adj := adjacent(idx)
	for _, a := range adj {
		if v := g[a]; v == val {
			return false
		}
	}
	return true
}

/* Don't guarantee sudoku puzzle will be right :) */
func RandomSudoku(closed int) Grid {
	grid := Grid(make([]int, side*side))
	for i := 0; i <= closed; i++ {
		t := true
		var idx, poss int
		for t {
			idx, poss = rand.Intn(side*side), rand.Intn(9)+1
			if val := grid[idx]; val != 0 { continue }
			t = !grid.verify(idx, poss)
		}
		grid[idx] = poss
	}
	return grid
}

func StaticSudoku() Grid {
	g := []int{
		1, 7, 9, 0, 0, 5, 4, 2, 0,
		0, 5, 3, 0, 4, 2, 0, 0, 1,
		0, 0, 6, 0, 9, 7, 0, 0, 8,
		0, 9, 5, 0, 0, 0, 1, 0, 4,
		0, 0, 2, 8, 6, 1, 3, 0, 0, 
		3, 0, 8, 0, 0, 0, 2, 6, 0,
		5, 0, 0, 9, 7, 0, 8, 0, 0,
		9, 0, 0, 2, 3, 0, 6, 1, 0, 
		0, 6, 4, 5, 0, 0, 7, 9, 3,
	}
	return Grid(g)
}

func (grid Grid) String() string {
	r := make([]string, side)
	g := []int(grid)
	for i := 0; i < side; i++ {
		r[i] = fmt.Sprintf("%v\n", g[(i*side):(i*side+side)])
	}
	return fmt.Sprint(r)
}
