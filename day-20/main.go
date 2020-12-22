package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	TOP    = 0
	RIGHT  = 1
	BOTTOM = 2
	LEFT   = 3
)

const (
	Monster = `                  # 
#    ##    ##    ###
 #  #  #  #  #  #   `
)

type Tile struct {
	body [][]int
	size int
}

func NewTile(body [][]int) *Tile {
	return &Tile{
		body: body,
		size: len(body),
	}
}

func (t *Tile) String() string {
	chs := make([]string, 0, len(t.body))
	chs = append(chs, fmt.Sprintf("\n\n%+v", t.Edges()))
	for _, bb := range t.body {
		var b bytes.Buffer
		for _, v := range bb {
			if v == 1 {
				b.WriteRune('#')
			} else if v == 2 {
				b.WriteRune('O')
			} else {
				b.WriteRune('.')
			}
		}
		chs = append(chs, b.String())
	}
	return strings.Join(chs, "\n")
}

func (t *Tile) Edges() [4]int {
	var res [4]int
	for i := 0; i < t.size; i++ {
		if t.body[0][i] > 0 {
			res[0] |= 1 << i
		}
		if t.body[i][t.size-1] > 0 {
			res[1] |= 1 << i
		}
		if t.body[t.size-1][i] > 0 {
			res[2] |= 1 << i
		}
		if t.body[i][0] > 0 {
			res[3] |= 1 << i
		}
	}
	return res
}

func (t *Tile) Side(dir int) int {
	switch dir {
	case TOP:
		return t.Top()
	case RIGHT:
		return t.Right()
	case BOTTOM:
		return t.Bottom()
	case LEFT:
		return t.Left()
	default:
		panic(fmt.Sprintf("unknown direction: %d", dir))
	}
}

func (t *Tile) Top() int {
	return t.Edges()[TOP]
}

func (t *Tile) Right() int {
	return t.Edges()[RIGHT]
}

func (t *Tile) Bottom() int {
	return t.Edges()[BOTTOM]
}

func (t *Tile) Left() int {
	return t.Edges()[LEFT]
}

func (t *Tile) Rotate() *Tile {
	nb := make([][]int, t.size)
	for i := 0; i < t.size; i++ {
		nb[i] = make([]int, t.size)
		for j := 0; j < t.size; j++ {
			nb[i][j] = t.body[t.size-1-j][i]
		}
	}
	return NewTile(nb)
}

func (t *Tile) FlipVert() *Tile {
	nb := make([][]int, t.size)
	for i := 0; i < t.size; i++ {
		nb[i] = make([]int, t.size)
		for j := 0; j < t.size; j++ {
			nb[i][j] = t.body[t.size-1-i][j]
		}
	}
	return NewTile(nb)
}

func (t *Tile) FlipHor() *Tile {
	nb := make([][]int, t.size)
	for i := 0; i < t.size; i++ {
		nb[i] = make([]int, t.size)
		for j := 0; j < t.size; j++ {
			nb[i][j] = t.body[i][t.size-1-j]
		}
	}
	return NewTile(nb)
}

func (t *Tile) Align(dir, want int) *Tile {
	switch dir {
	case TOP:
		return t.AlignTop(want)
	case RIGHT:
		return t.AlignRight(want)
	case BOTTOM:
		return t.AlignBottom(want)
	case LEFT:
		return t.AlignLeft(want)
	default:
		panic(fmt.Sprintf("unknown direction: %d", dir))
	}
}

func (t *Tile) AlignTop(want int) *Tile {
	ix := -1
	for i, edge := range t.Edges() {
		if edge == want {
			ix = i
			break
		}
	}
	switch ix {
	case TOP:
		return t
	case RIGHT:
		return t.Rotate().Rotate().Rotate()
	case BOTTOM:
		return t.FlipVert()
	case LEFT:
		return t.FlipHor().Rotate().Rotate().Rotate()
	}
	return nil
}

func (t *Tile) AlignRight(want int) *Tile {
	ix := -1
	for i, edge := range t.Edges() {
		if edge == want {
			ix = i
			break
		}
	}
	switch ix {
	case TOP:
		return t.Rotate()
	case RIGHT:
		return t
	case BOTTOM:
		return t.FlipVert().Rotate()
	case LEFT:
		return t.FlipHor()
	}
	return nil
}

func (t *Tile) AlignBottom(want int) *Tile {
	ix := -1
	for i, edge := range t.Edges() {
		if edge == want {
			ix = i
			break
		}
	}
	switch ix {
	case TOP:
		return t.FlipVert()
	case RIGHT:
		return t.Rotate().Rotate().Rotate().FlipVert()
	case BOTTOM:
		return t
	case LEFT:
		return t.Rotate().Rotate().Rotate()
	}
	return nil
}

func (t *Tile) AlignLeft(want int) *Tile {
	ix := -1
	for i, edge := range t.Edges() {
		if edge == want {
			ix = i
			break
		}
	}
	switch ix {
	case TOP:
		return t.Rotate().FlipHor()
	case RIGHT:
		return t.FlipHor()
	case BOTTOM:
		return t.Rotate()
	case LEFT:
		return t
	}
	return nil
}

func (t *Tile) ConnectsTo(want int) bool {
	for _, edge := range t.Edges() {
		if edge == want {
			return true
		}
	}
	return false
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("unexpected error: %s", err.Error())
	}
}

func readfile(path string) string {
	f, err := ioutil.ReadFile(path)
	noerr(err)
	return strings.TrimRight(string(f), "\r\n")
}

func parseTileId(s string) int {
	chs := strings.SplitN(s, " ", 2)
	n, err := strconv.Atoi(chs[1][:len(chs[1])-1])
	noerr(err)
	return n
}

func parseTile(s string) *Tile {
	lines := strings.Split(s, "\n")
	size := len(lines)
	body := make([][]int, size)
	for i, line := range lines {
		body[i] = make([]int, size)
		for j, r := range line {
			if r == '#' {
				body[i][j] = 1
			}
		}
	}
	return NewTile(body)
}

func makeMap(tiles map[int]*Tile) [][]*Tile {
	sz := int(math.Sqrt(float64(len(tiles))))

	tileMap := make([][]*Tile, sz)
	for i := 0; i < sz; i++ {
		tileMap[i] = make([]*Tile, sz)
	}

	rem := make(map[int]struct{})
	anyId := -1
	for tid := range tiles {
		if anyId < 0 {
			anyId = tid
		}
		rem[tid] = struct{}{}
	}

	field := make(map[[2]int]*Tile)
	steps := [][3]int{{TOP, -1, 0}, {RIGHT, 0, 1}, {BOTTOM, 1, 0}, {LEFT, 0, -1}}

	var head [3]int
	q := make([][3]int, 0, len(tiles))
	q = append(q, [3]int{anyId, 0, 0})
	mini, minj := 0, 0
	for len(q) > 0 {
		head, q = q[0], q[1:]
		tid, i, j := head[0], head[1], head[2]
		if i < mini {
			mini = i
		}
		if j < minj {
			minj = j
		}
		tile := tiles[tid]
		field[[2]int{i, j}] = tiles[tid]
		delete(rem, tid)
		for _, step := range steps {
			newpos := [2]int{i + step[1], j + step[2]}
			if _, ok := field[newpos]; ok {
				// already placed a tile there
				continue
			}
			dir := step[0]
			side := tile.Side(dir)
		Neighbors:
			for nbid := range rem {
				nbtile := tiles[nbid]
				connects := false
				if nbtile.ConnectsTo(side) {
					connects = true
				} else if vert := nbtile.FlipVert(); vert.ConnectsTo(side) {
					nbtile = vert
					connects = true
				} else if hor := nbtile.FlipHor(); hor.ConnectsTo(side) {
					nbtile = hor
					connects = true
				}
				if connects {
					tiles[nbid] = nbtile.Align((dir+2)%4, side)
					q = append(q, [3]int{nbid, newpos[0], newpos[1]})
					break Neighbors
				}
			}
		}
	}

	for pos, tile := range field {
		tileMap[pos[0]-mini][pos[1]-minj] = tile
	}

	return tileMap
}

func trimTile(tile *Tile) *Tile {
	nb := make([][]int, tile.size-2)
	for i := 0; i < len(nb); i++ {
		nb[i] = make([]int, tile.size-2)
		for j := 0; j < len(nb); j++ {
			nb[i][j] = tile.body[i+1][j+1]
		}
	}
	return NewTile(nb)
}

func makeMegaTile(tileMap [][]*Tile) *Tile {
	nsz := len(tileMap) * tileMap[0][0].size
	nb := make([][]int, nsz)
	for i := 0; i < nsz; i++ {
		nb[i] = make([]int, nsz)
	}
	for x := 0; x < len(tileMap); x++ {
		for y := 0; y < len(tileMap); y++ {
			tile := tileMap[x][y]
			for i := 0; i < tile.size; i++ {
				for j := 0; j < tile.size; j++ {
					nb[x*tile.size+i][y*tile.size+j] = tile.body[i][j]
				}
			}
		}
	}
	return NewTile(nb)
}

func parseMonster(pattern string) [][2]int {
	res := [][2]int{}
	for i, line := range strings.Split(pattern, "\n") {
		for j := 0; j < len(line); j++ {
			if line[j] == '#' {
				res = append(res, [2]int{i, j})
			}
		}
	}
	return res
}

func findAllMonsters(tile *Tile, monster [][2]int) *Tile {
	nb := make([][]int, tile.size)
	for i := 0; i < tile.size; i++ {
		nb[i] = make([]int, tile.size)
		copy(nb[i], tile.body[i])
	}
	anymatch := false
Row:
	for i := 0; i < tile.size; i++ {
	Column:
		for j := 0; j < tile.size; j++ {
			for _, pos := range monster {
				ni, nj := i+pos[0], j+pos[1]
				if ni >= tile.size {
					break Row
				}
				if nj >= tile.size {
					continue Row
				}
				if tile.body[ni][nj] != 1 {
					continue Column
				}
			}
			anymatch = true
			for _, pos := range monster {
				nb[i+pos[0]][j+pos[1]] = 2
			}
		}
	}
	if anymatch {
		return NewTile(nb)
	}
	return nil
}

func countRoughness(tile *Tile) int {
	cnt := 0
	for i := 0; i < tile.size; i++ {
		for j := 0; j < tile.size; j++ {
			if tile.body[i][j] == 1 {
				cnt++
			}
		}
	}
	return cnt
}

func main() {
	tiles := make(map[int]*Tile)
	chunks := strings.Split(readfile("input"), "\n\n")
	for _, ch := range chunks {
		chch := strings.SplitN(ch, "\n", 2)
		header := chch[0]
		body := chch[1]
		tid := parseTileId(header)
		tile := parseTile(body)
		tiles[tid] = tile
	}

	tileMap := makeMap(tiles)

	for i := 0; i < len(tileMap); i++ {
		for j := 0; j < len(tileMap); j++ {
			tileMap[i][j] = trimTile(tileMap[i][j])
		}
	}

	megatile := makeMegaTile(tileMap)
	monster := parseMonster(Monster)

	for _, tile := range []*Tile{
		megatile, megatile.FlipVert(), megatile.FlipHor(),
		megatile.Rotate(), megatile.Rotate().FlipVert(), megatile.Rotate().FlipHor(),
		megatile.Rotate().Rotate(), megatile.Rotate().Rotate().FlipVert(), megatile.Rotate().Rotate().FlipHor(),
		megatile.Rotate().Rotate().Rotate(), megatile.Rotate().Rotate().Rotate().FlipVert(), megatile.Rotate().Rotate().Rotate().FlipHor(),
	} {
		if found := findAllMonsters(tile, monster); found != nil {
			log.Printf("moster:\n%s", found)
			cnt := countRoughness(found)
			log.Printf("the answer is: %d", cnt)
			break
		}
	}
}
