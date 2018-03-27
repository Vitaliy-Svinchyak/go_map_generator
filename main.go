package main

import (
	"fmt"
	"time"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"encoding/json"
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
var fieldMap [][]string
var fieldSize = make(map[string]int)
var crossedRooms []room

var cursorRectangle = point{x: 0, y: 0}
var rawStart point

func main() {
	rand.Seed(time.Now().Unix())

	start := time.Now()
	fieldMap := rooms(100, 100)
	elapsed := time.Since(start)
	fmt.Printf("go execution %v", elapsed)

	fmt.Println("")
	fmt.Println("")
	field, err := json.Marshal(fieldMap)

	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println(string(field))
}

// tested
func getRandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// tested
func empty(rows, cells int) [][]string {

	var field [][]string
	var item string

	for y := 0; y <= rows; y++ {
		field = append(field, []string{})

		for x := 0; x <= cells; x++ {
			if y == 0 || y == rows || x == 0 || x == cells {
				item = types["rock"]
			} else {
				item = types["empty"]
			}

			field[y] = append(field[y], item)
		}
	}

	return field
}

func rooms(rows, cells int) [][]string {
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

func drawRooms() {
	rawStart = point{x: 0, y: 0}
	for !generated {
		var room = getRandomRoomSize()
		cursorRectangle = drawRoom(cursorRectangle, room)

		if cursorRectangle.x == fieldSize["maxX"] {
			moveCursorToTheNewRow()
			minYOnRow = -1
		}
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

func concatMaps(a, b []string) []string {
	for _, v := range b {
		a = append(a, v)
	}

	return a
}

// tested
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
		var wallsInRange []string

		// We must cling to the "lowest" room in our range
		for _, room := range createdRooms {
			if (room.start.x >= startX && room.start.x <= roomEndX) || (room.end.x >= startX && room.end.x <= roomEndX) {
				wallsInRange = concatMaps(wallsInRange, getWallsInRange(room, cursorRectangle.x, roomEndX))
			}
		}

		wallsInRange = deleteDuplications(wallsInRange)

		if len(wallsInRange) > 0 {
			startY = findMinY(wallsInRange)
		}

		for _, room := range createdRooms {
			if (room.end.x >= cursorRectangle.x || room.start.x <= roomEndX) && room.end.y >= startY {
				crossedRooms = append(crossedRooms, room)
			}
		}
	}

	return startY
}

// tested
func getKeys(array map[string]bool) []string {
	keys := reflect.ValueOf(array).MapKeys()

	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}

	return strkeys
}

// tested
func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}

	return ret
}

// tested
func find(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

// tested
func getWallsInRange(room room, startX int, endX int) []string {
	keys := getKeys(room.occupiedWalls)
	allCells := filter(keys, func(v string) bool {
		xCoordinate, _ := strconv.Atoi(strings.Split(v, ":")[1])

		return xCoordinate >= startX && xCoordinate <= endX
	})

	return filter(allCells, func(v string) bool {
		var coordinates = strings.Split(v, ":")
		yCoordinate, _ := strconv.Atoi(coordinates[0])
		xCoordinate, _ := strconv.Atoi(coordinates[1])

		return find(allCells, fmt.Sprintf("%d:%d", yCoordinate+1, xCoordinate)) == false
	})

}

// tested
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

// tested
func deleteDuplications(wallsInRange []string) []string {
	var filtered []string

	for _, wall := range wallsInRange {
		var deleteIt = false
		var wallCoordinates = strings.Split(wall, ":")
		yCoordinate, _ := strconv.Atoi(wallCoordinates[0])

		for _, wallToCompare := range wallsInRange {
			var wallToCompareCoordinates = strings.Split(wallToCompare, ":")
			yCoordinateCompare, _ := strconv.Atoi(wallToCompareCoordinates[0])

			if wallCoordinates[1] == wallToCompareCoordinates[1] && yCoordinate < yCoordinateCompare {
				deleteIt = true
			}
		}

		if !deleteIt {
			filtered = append(filtered, wall)
		}
	}

	return filtered
}

// tested
func findMinY(wallsInRange []string) int {
	var firstWallCoordinates = strings.Split(wallsInRange[0], ":")
	var minY, _ = strconv.Atoi(firstWallCoordinates[0])

	for _, wall := range wallsInRange {
		var wallCoordinates = strings.Split(wall, ":")
		var currentY, _ = strconv.Atoi(wallCoordinates[0])

		if currentY < minY {
			minY = currentY
		}
	}

	return minY
}
