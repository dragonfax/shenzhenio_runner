package main

import (
	"bufio"
	"fmt"
	"github.com/yuin/gopher-lua"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type TerminalType int

const (
	SIMPLE TerminalType = 0
	XBUS   TerminalType = 1
)

type TerminalDirection int

const (
	INPUT  TerminalDirection = 0
	OUTPUT TerminalDirection = 1
)

type Terminal struct {
	Name      string
	Type      TerminalType
	Direction TerminalDirection
	Index     int
	X         int
	Y         int
}

type TerminalList []*Terminal

func NewTerminalList() TerminalList {
	return TerminalList(make([]*Terminal, 0, 1))
}

func (tl TerminalList) findTerminal(index int) *Terminal {
	for _, t := range tl {
		if t.Index == index {
			return t
		}
	}
	panic("couldn't find terminal by index")
}

func ParseProblemFile(file string) TerminalList {

	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("library.lua"); err != nil {
		panic(err)
	}

	if err := L.DoFile(file); err != nil {
		panic(err)
	}

	terminals := parseTerminals(L)

	sockets := parseBoard(L)

	addPositionToTerminals(terminals, sockets)

	return terminals
}

type Socket struct {
	Index int
	X     int
	Y     int
}

func parseBoard(L *lua.LState) []Socket {

	err := L.CallByParam(
		lua.P{
			Fn:      L.GetGlobal("get_board"),
			NRet:    1,
			Protect: false,
		},
	)
	if err != nil {
		panic(err)
	}

	ret := L.Get(-1)
	L.Pop(1)

	board_string := string(ret.(lua.LString))

	sockets := make([]Socket, 0, 1)

	scanner := bufio.NewScanner(strings.NewReader(board_string))
	y := 0
	x := 0
LINE:
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break LINE
		}
		for _, r := range line {
			switch {
			case r == '\n':
			case r == ' ':
				// skip
			case r == '.':
				// empty cell
				x += 1
			case r == '#':
				// board cell, skipo
				x += 1
			case unicode.IsDigit(r):
				// add a new socket to the board
				i, err := strconv.Atoi(string(r))
				if err != nil {
					panic("failed")
				}
				sockets = append(sockets, Socket{i, x, y})
				x += 1
			default:
				panic("unknown cell in board definition '" + string(r) + "'")
			}
		}
		y += 1
		x = 0
	}

	return sockets
}

func addPositionToTerminals(terminals TerminalList, sockets []Socket) {

	for _, socket := range sockets {
		terminal := terminals.findTerminal(socket.Index)
		terminal.X = socket.X
		terminal.Y = socket.Y
	}
}

func parseTerminals(L *lua.LState) TerminalList {

	err := L.CallByParam(
		lua.P{
			Fn:      L.GetGlobal("get_data"),
			NRet:    1,
			Protect: false,
		},
	)
	if err != nil {
		panic(err)
	}

	terminals_lua := L.GetGlobal("TERMINALS").(*lua.LTable)

	terminals := NewTerminalList()
	terminals_lua.ForEach(func(name lua.LValue, tdata lua.LValue) {

		table := tdata.(*lua.LTable)

		direction_s := string(table.RawGetString("direction").(lua.LString))

		var direction TerminalDirection
		switch direction_s {
		case "input":
			direction = INPUT
		case "output":
			direction = OUTPUT
		default:
			panic("unknown terminal direction " + direction_s)

		}

		type_s := string(table.RawGetString("type").(lua.LString))

		var tt TerminalType
		switch type_s {
		case "simple":
			tt = SIMPLE
		case "xbus":
			tt = XBUS
		default:
			panic("unknown terminal type " + type_s)

		}

		index_s := string(table.RawGetString("index").(lua.LString))
		index, err := strconv.Atoi(index_s)
		if err != nil {
			panic("bad number in string '" + index_s + "'")
		}

		terminal := &Terminal{
			Name:      name.String(),
			Type:      tt,
			Direction: direction,
			Index:     index,
		}
		terminals = append(terminals, terminal)
	})

	return terminals
}
