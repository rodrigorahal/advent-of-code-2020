package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type Tile struct {
	ID    int
	Image [][]string
}

func (t Tile) Row(i int) string {
	r := t.Image[i]
	return strings.Join(r, "")
}

func (t Tile) Column(i int) string {
	var c []string
	for _, row := range t.Image {
		c = append(c, row[i])
	}
	return strings.Join(c, "")
}

func (t Tile) FlipY() Tile {
	new := Tile{ID: t.ID, Image: nil}

	for i := len(t.Image) - 1; i >= 0; i-- {
		row := t.Image[i]
		new.Image = append(new.Image, make([]string, len(row)))
		j := len(t.Image) - 1 - i
		new.Image[j] = row
	}
	return new
}

func (t Tile) FlipX() Tile {
	new := Tile{ID: t.ID, Image: nil}

	for i, row := range t.Image {
		new.Image = append(new.Image, make([]string, len(row)))
		for j := len(row) - 1; j >= 0; j-- {
			new.Image[i][len(row)-1-j] = row[j]
		}
	}
	return new
}

func (t Tile) RotateRight() Tile {
	new := Tile{ID: t.ID, Image: nil}

	rows := len(t.Image)

	for _, row := range t.Image {
		new.Image = append(new.Image, make([]string, len(row)))
	}

	for i, row := range t.Image {
		for j, pixel := range row {
			new.Image[j][rows-1-i] = pixel
		}
	}
	return new
}

func (t Tile) Print() {
	fmt.Printf("Tile %d:\n", t.ID)
	for _, row := range t.Image {
		fmt.Println(strings.Join(row, ""))
	}
}

func (t Tile) Borders() []string {
	/* Returns borders in clockwise order
	Starting from the top
	*/
	var bs []string

	i := t.Image
	h := len(i)
	w := len(i[0])

	bs = append(bs, t.Row(0))
	bs = append(bs, t.Column(w-1))
	bs = append(bs, t.Row(h-1))
	bs = append(bs, t.Column(0))

	return bs
}

func (t Tile) Orientations() []Tile {
	var possible []Tile

	possible = append(possible, t)
	possible = append(possible, t.FlipX())
	possible = append(possible, t.FlipY())
	for _, p := range possible {
		prev := p
		for i := 0; i < 3; i++ {
			rot := prev.RotateRight()
			possible = append(possible, prev.RotateRight())
			prev = rot
		}
	}
	return possible
}

func (t Tile) Match(other Tile) (Tile, int, int) {
	for _, orientation := range other.Orientations() {
		for i, border := range t.Borders() {
			for j, orientationBorder := range orientation.Borders() {
				if border == orientationBorder && isAdjacent(i, j) {
					// It's a match!
					/*
						0: top
						1: right
						2: bottom
						3: left
					*/
					return orientation, i, j
				}
			}
		}
	}
	return t, -1, -1
}

func isAdjacent(i, j int) bool {
	return (i == 0 && j == 2) ||
		(i == 1 && j == 3) ||
		(i == 2 && j == 0) ||
		(i == 3 && j == 1)
}

func createImage(tileByTileID map[int]Tile, associationsByTileID map[int]map[int]int) Tile {
	type Pos struct {
		ID, X, Y int
	}

	var topLeftID int
	for tileID, directionByTileID := range associationsByTileID {
		_, hasRight := directionByTileID[1]
		_, hasBottom := directionByTileID[2]
		if len(directionByTileID) == 2 && hasRight && hasBottom {
			topLeftID = tileID
			break
		}
	}

	topLeft := tileByTileID[topLeftID]

	size := int(math.Sqrt(float64(len(tileByTileID)))) * len(topLeft.Image)

	image := make([][]string, size)
	for i := 0; i < size; i++ {
		image[i] = make([]string, size)
	}

	seen := make(map[int]bool)
	queue := []Pos{Pos{ID: topLeftID, X: 0, Y: 0}}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		tID := pos.ID
		x := pos.X
		y := pos.Y

		t := tileByTileID[tID]

		seen[tID] = true

		for i, row := range t.Image {
			for j, v := range row {
				image[y+i][x+j] = v
			}
		}

		for dir, adjID := range associationsByTileID[tID] {
			if _, ok := seen[adjID]; ok {
				continue
			}

			var dx int
			var dy int
			var mul int = len(t.Image)
			if dir == 1 {
				dx = 1 * mul
			} else if dir == 3 {
				dx = -1 * mul
			}

			if dir == 0 {
				dy = -1 * mul
			} else if dir == 2 {
				dy = 1 * mul
			}

			queue = append(queue, Pos{ID: adjID, X: x + dx, Y: y + dy})
		}
	}

	return Tile{ID: 0, Image: image}
}

func assemble(ts []Tile) Tile {
	/*
		Assemble a full image from Tiles
	*/
	tileByTileID := make(map[int]Tile)
	associationsByTileID := make(map[int]map[int]int)
	/*
		ID: {
			0: tileID-on-top
			1: tileID-on-right
			2: tileID-on-left
			3: tileID-on-bottom
		}
	*/

	for _, t := range ts {
		tileByTileID[t.ID] = t
	}

	matches := findMatches(ts)

	startID := ts[0].ID

	queue := []int{startID}
	seen := make(map[int]bool)

	for len(queue) > 0 {
		tID := queue[0]
		queue = queue[1:]
		seen[tID] = true

		for _, matchID := range matches[tID] {
			if _, ok := seen[matchID]; ok {
				continue
			}
			t := tileByTileID[tID]
			t2 := tileByTileID[matchID]

			o, i, j := t.Match(t2)

			queue = append(queue, matchID)

			if len(associationsByTileID[tID]) == 0 {
				associationsByTileID[tID] = make(map[int]int)
			}
			associationsByTileID[tID][i] = matchID

			tileByTileID[matchID] = o

			if len(associationsByTileID[matchID]) == 0 {
				associationsByTileID[matchID] = make(map[int]int)
			}
			associationsByTileID[matchID][j] = tID
		}
	}

	return createImage(tileByTileID, associationsByTileID)
}

func removeBorders(image Tile, tileSize int) Tile {
	var newImage [][]string

	for i, row := range image.Image {
		if i%tileSize == 0 || (i+1)%tileSize == 0 {
			continue
		}
		newImage = append(newImage, []string{})
		for j, v := range row {
			if j%tileSize == 0 || (j+1)%tileSize == 0 {
				continue
			}
			newImage[len(newImage)-1] = append(newImage[len(newImage)-1], v)
		}
	}
	return Tile{ID: 0, Image: newImage}
}

func hasMonsterAt(image Tile, i, j int) bool {
	img := image.Image
	return img[i][j+18] == "#" &&
		img[i+1][j+0] == "#" &&
		img[i+1][j+5] == "#" &&
		img[i+1][j+6] == "#" &&
		img[i+1][j+11] == "#" &&
		img[i+1][j+12] == "#" &&
		img[i+1][j+17] == "#" &&
		img[i+1][j+18] == "#" &&
		img[i+1][j+19] == "#" &&
		img[i+2][j+1] == "#" &&
		img[i+2][j+4] == "#" &&
		img[i+2][j+7] == "#" &&
		img[i+2][j+10] == "#" &&
		img[i+2][j+13] == "#" &&
		img[i+2][j+16] == "#"
}

func count(t Tile) int {
	/*
		Scan image looking for monsters
		Count # that are not part of a monster
	*/
	var monsters int
	var found bool
	var image Tile

	for _, o := range t.Orientations() {
		for i := 0; i < len(o.Image)-3; i++ {
			for j := 0; j < len(o.Image[0])-19; j++ {
				if hasMonsterAt(o, i, j) {
					monsters++
					found = true
				}
			}
		}
		if found {
			image = o
			break
		}
	}

	var sum int
	for _, row := range image.Image {
		for _, v := range row {
			if v == "#" {
				sum++
			}
		}
	}

	return sum - monsters*15
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func findMatches(ts []Tile) map[int][]int {
	matches := make(map[int][]int)

	for _, tile := range ts {
		matches[tile.ID] = make([]int, 0)
		for _, border := range tile.Borders() {
			matched := false
			for _, candidate := range ts {
				if tile.ID == candidate.ID {
					continue
				}
				for _, candidateBorder := range candidate.Borders() {
					if border == candidateBorder ||
						border == reverse(candidateBorder) {
						matches[tile.ID] = append(matches[tile.ID], candidate.ID)
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
		}
	}
	return matches
}

func findCorners(ts []Tile) int {
	var corners []Tile

	for _, tile := range ts {
		matches := 0
		for _, border := range tile.Borders() {
			matched := false
			for _, candidate := range ts {
				if tile.ID == candidate.ID {
					continue
				}
				for _, candidateBorder := range candidate.Borders() {
					if border == candidateBorder ||
						border == reverse(candidateBorder) {
						matched = true
						matches++
						break
					}
				}
				if matched {
					break
				}
			}
		}

		if matches == 2 {
			corners = append(corners, tile)
		}
	}

	product := 1
	for _, t := range corners {
		product *= t.ID
	}
	return product
}

func read(r io.Reader) []Tile {
	var tiles []Tile
	var tile Tile
	var row int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "Tile") {
			t = strings.ReplaceAll(t, ":", "")
			sid := strings.Split(t, " ")[1]
			id, _ := strconv.Atoi(sid)
			image := make([][]string, 0)
			new := Tile{id, image}
			tile = new
			row = 0
		} else if t == "" {
			tiles = append(tiles, tile)
		} else {
			r := make([]string, 0)
			tile.Image = append(tile.Image, r)
			for _, char := range t {
				tile.Image[row] = append(tile.Image[row], string(char))
			}
			row++
		}
	}
	tiles = append(tiles, tile)
	return tiles
}

func main() {
	file, _ := os.Open("input.txt")
	tiles := read(file)

	fmt.Println(findCorners(tiles))

	image := assemble(tiles)
	image = removeBorders(image, len(tiles[0].Image))
	fmt.Println(count(image))
}
