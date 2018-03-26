package main

import (
	"fmt"
	"time"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type room struct {
	start         point
	end           point
	occupiedCells map[string]bool
	occupiedWalls map[string]bool
}

var types = map[string]string{
	"rock":  "■",
	"human": "o",
	"tree":  "֏",
	"empty": " ",
	"track": "x",
	"route": "w",
}

var createdRooms []room
var generated = false
var minYOnRow = -1

var minRoomSize = map[string]int{"x": 4, "y": 4}
var maxRoomSize = map[string]int{"x": 9, "y": 9}
var fieldMap map[int]map[int]string
var fieldSize = make(map[string]int)
var crossedRooms = map[int]room{}

func main() {
	start := time.Now()
	fieldMap := rooms(10, 10)
	elapsed := time.Since(start)

	//field, _ := json.Marshal(fieldMap)
	//fmt.Println(string(field))
	fmt.Println("")
	fmt.Println("")
	fmt.Println(fieldMap)

	fmt.Println("go execution", elapsed)
}

// tested
func getRandomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// tested
func empty(rows, cells int) map[int]map[int]string {

	var field = map[int]map[int]string{}
	var item string

	for y := 0; y <= rows; y++ {
		field[y] = map[int]string{}

		for x := 0; x <= cells; x++ {
			if y == 0 || y == rows || x == 0 || x == cells {
				item = types["rock"]
			} else {
				item = types["empty"]
			}

			field[y][x] = item
		}
	}

	return field
}

func rooms(rows, cells int) map[int]map[int]string {
	fieldMap = empty(rows, cells)

	fieldSize["maxY"] = rows
	fieldSize["maxX"] = cells

	drawRooms()
	//clearFatWalls()
	//drawDoors()
	fieldMap[1][1] = types["human"]
	fieldMap[1][len(fieldMap[0])-2] = types["human"]
	fieldMap[len(fieldMap)-2][len(fieldMap[0])-2] = types["human"]
	fieldMap[len(fieldMap)-2][1] = types["human"]
	return fieldMap
}

var cursorRectangle point
var rawStart point

func drawRooms() {
	//cursorRectangle := point{x: 0, y: 0}
	//rawStart := point{x: 0, y: 0}
	var a = 10
	//for !generated {
	for a > 0 {
		var room = getRandomRoomSize()
		cursorRectangle = drawRoom(cursorRectangle, room)

		//if cursorRectangle["x"] == fieldSize["maxX"] {
		moveCursorToTheNewRow()
		//	minYOnRow = -1
		//}
		a--
	}
}

func getRandomRoomSize() point {
	var x = getRandomInt(minRoomSize["x"], maxRoomSize["x"])
	var y = getRandomInt(minRoomSize["y"], maxRoomSize["y"])
	if cursorRectangle.x+x > fieldSize["maxX"] {
		x = fieldSize["maxX"] - cursorRectangle.x
	}

	if cursorRectangle.y+y > fieldSize["maxY"] {
		y = fieldSize["maxY"] - cursorRectangle.y
	}

	if cursorRectangle.x+x == fieldSize["maxX"]-1 {
		x++
	}

	return point{
		x: x,
		y: y,
	}
}

func drawRoom(cursorRectangle point, roomToDraw point) point {
	var roomEndX = cursorRectangle.x + roomToDraw.x
	var startY = getStartRowForRoom(cursorRectangle.x, roomEndX)
	var roomEndY = startY + roomToDraw.y

	if roomEndY > fieldSize["maxY"] {
		roomEndY = fieldSize["maxY"]
	}

	if roomEndY+1 == fieldSize["maxY"] {
		roomEndY = fieldSize["maxY"]
	}

	var occupiedCells = map[string]bool{}
	var occupiedWalls = map[string]bool{}

	for y := startY; y <= roomEndY; y++ {
		for x := cursorRectangle.x; x <= roomEndX; x++ {
			occupiedCells[ fmt.Sprintf("%d:%d", y, x)] = true
		}
	}

	for y := startY; y <= roomEndY; y++ {
		if !wallsIntersect(y, cursorRectangle.x) {
			fieldMap[y][cursorRectangle.x] = types["rock"]
			occupiedWalls[fmt.Sprintf("%d:%d", y, cursorRectangle.x)] = true
		}

		if !wallsIntersect(y, roomEndX) {
			fieldMap[y][roomEndX] = types["rock"]
			occupiedWalls[fmt.Sprintf("%d:%d", y, roomEndX)] = true
		}
	}

	for x := cursorRectangle.x; x <= roomEndX; x++ {
		if !wallsIntersect(startY, x) {
			fieldMap[startY][x] = types["rock"]
			occupiedWalls[fmt.Sprintf("%d:%d", startY, x)] = true
		}

		if !wallsIntersect(roomEndY, x) {
			fieldMap[roomEndY][x] = types["rock"]
			occupiedWalls[fmt.Sprintf("%d:%d", roomEndY, x)] = true
		}
	}

	createdRooms = append(createdRooms, room{
		start:         point{x: cursorRectangle.x, y: startY},
		end:           point{x: roomEndX, y: roomEndY},
		occupiedCells: occupiedCells,
		occupiedWalls: occupiedWalls,
	})

	if minYOnRow > roomEndY || minYOnRow == -1 {
		minYOnRow = roomEndY
	}

	return point{y: startY, x: roomEndX}
}

func concatMaps(a, b map[string]bool) map[string]bool {
	for k, v := range b {
		a[k] = v
	}

	return a
}

// to test
func wallsIntersect(y, x int) bool {
	for _, crossedRoom := range crossedRooms {
		if crossedRoom.occupiedCells[fmt.Sprintf("%d:%d", y, x)] {
			return true
		}
	}

	return false
}

func getStartRowForRoom(startX, roomEndX int) int {
	var startY = cursorRectangle.y

	if rawStart.y != 0 {
		var wallsInRange = map[string]bool{}

		// We must cling to the "lowest" room in our range
		for _, room := range createdRooms {
			if (room.start.x >= startX && room.start.x <= roomEndX) || (room.end.x >= startX && room.end.x <= roomEndX) {
				wallsInRange = concatMaps(wallsInRange, getWallsInRange(room, cursorRectangle.x, roomEndX))
			}
		}

		//wallsInRange = deleteDuplications(wallsInRange)
		//
		//if len(wallsInRange) > 0 {
		//	startY = findMinY(wallsInRange)
		//}
		//
		//for _, room := range createdRooms {
		//	if (room.end.x >= cursorRectangle.x || room.start.x <= roomEndX) && room.end.y >= startY {
		//		crossedRooms = append(crossedRooms, room)
		//	}
		//}
	}

	return startY
}

func getKeys(array map[string]bool) []string {
	keys := reflect.ValueOf(array).MapKeys()

	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}

	return strkeys
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func find(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func getWallsInRange(room room, startX int, endX int) map[string]bool {
	allCells := filter(getKeys(room.occupiedWalls), func(v string) bool {
		xCoordinate, _ := strconv.Atoi(strings.SplitAfter(v, ":")[1])

		return xCoordinate >= startX && xCoordinate <= endX
	})

	filter(allCells, func(v string) bool {
		var coordinates = strings.SplitAfter(v, ":")
		yCoordinate, _ := strconv.Atoi(coordinates[0])
		xCoordinate, _ := strconv.Atoi(coordinates[1])

		return find(allCells, fmt.Sprintf("%d:%d", yCoordinate+1, xCoordinate)) == false
	})

	return map[string]bool{}

}

func moveCursorToTheNewRow() {
	if rawStart.y == fieldSize["maxY"] {
		generated = true
	}

	var roomWithMaxY = createdRooms[0]

	for _, room := range createdRooms {
		if room.start.x == 0 && room.end.y >= roomWithMaxY.end.y {
			roomWithMaxY = room
		}
	}

	cursorRectangle = point{x: 0, y: roomWithMaxY.end.y}
	rawStart = cursorRectangle
}
