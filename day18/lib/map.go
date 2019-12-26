package day18

type Map [][]string //uint8 ?

func (m Map) ToString() string {
	result := ""
	for y, line := range m {
		if y > 0 {
			result += "\n"
		}
		for _, s := range line {
			result += s
		}
	}
	return result
}

func (m Map) ToStringModified(transform func(Point) *string) string {
	result := ""
	for y, line := range m {
		if y > 0 {
			result += "\n"
		}
		for x, s := range line {
			t := transform(Point{x, y})
			if t != nil {
				result += *t
			} else {
				result += s
			}
		}
	}
	return result
}

func (m Map) Each(handler func(p Point, v string)) {
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			handler(Point{x, y}, m[y][x])
		}
	}
}

func (m Map) EachColumnCell(x int, handler func(p Point, v string)) {
	for y := 0; y < len(m); y++ {
		handler(Point{x, y}, m[y][x])
	}
}

func (m Map) EachRowColumn(y int, handler func(p Point, v string)) {
	for x := 0; x < len(m[y]); x++ {
		handler(Point{x, y}, m[y][x])
	}
}

func (m Map) EachValue(handler func(v string)) {
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			handler(m[y][x])
		}
	}
}

func (m Map) Filter(filter func(v string) bool) []Point {
	result := make([]Point, 0)
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			if filter(m[y][x]) {
				result = append(result, Point{x, y,})
			}
		}
	}
	return result
}

func (m Map) Count(filter func(v string) bool) int {
	total := 0
	m.EachValue(func(v string) {
		if filter(v) {
			total++
		}
	})
	return total
}

func (m Map) FindFirst(filter func(v string) bool) *Point {
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			if filter(m[y][x]) {
				return &Point{x, y}
			}
		}
	}
	return nil
}

func (m Map) Set(p Point, v string) *Map {
	m[p.Y][p.X] = v
	return &m
}

func (m Map) Contains(p Point) bool {
	if p.X < 0 || p.Y < 0 {
		return false
	}
	return len(m) > p.Y && len(m[p.Y]) > p.X
}

func (m Map) Get(p Point) *string {
	if m.Contains(p) {
		return &m[p.Y][p.X]
	} else {
		return nil
	}
}

func (m Map) Clone() *Map {
	clone := make(Map, len(m))
	for y := 0; y < len(m); y++ {
		clone[y ] = make([]string, len(m[y]))
		for x := 0; x < len(m[y]); x++ {
			clone[y][x] = m[y][x]
		}
	}
	return &clone
}

func (m Map) Height() int {
	return len(m)
}

func (m Map) Width() int {
	return len(m[0])
}

func NewMap(lines []string) Map {
	result := make(Map, 0)
	for _, line := range lines {
		r := make([]string, 0)
		for _, c := range line {
			r = append(r, string(c))
		}
		result = append(result, r)
	}
	return result
}
