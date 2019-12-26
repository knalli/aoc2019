package dayless

import "strings"

type Direction string

func (d Direction) Return() Direction {
	switch d {
	case NORTH:
		return SOUTH
	case EAST:
		return WEST
	case SOUTH:
		return NORTH
	case WEST:
		return EAST
	}
	panic("invalid direction")
}

func (d Direction) Left() Direction {
	switch d {
	case NORTH:
		return WEST
	case EAST:
		return NORTH
	case SOUTH:
		return EAST
	case WEST:
		return SOUTH
	}
	panic("invalid direction")
}

func (d Direction) Right() Direction {
	switch d {
	case NORTH:
		return EAST
	case EAST:
		return SOUTH
	case SOUTH:
		return WEST
	case WEST:
		return NORTH
	}
	panic("invalid direction")
}

func (d Direction) ToString() string {
	switch d {
	case NORTH:
		return "north"
	case EAST:
		return "east"
	case SOUTH:
		return "south"
	case WEST:
		return "west"
	}
	panic("invalid direction")
}

const NORTH Direction = "north"
const EAST Direction = "east"
const SOUTH Direction = "south"
const WEST Direction = "west"

func NewDirection(str string) Direction {
	switch strings.ToLower(str) {
	case "north":
		return NORTH
	case "east":
		return EAST
	case "south":
		return SOUTH
	case "west":
		return WEST
	}
	panic("invalid direction")
}
