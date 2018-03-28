package main

import (
	"fmt"
	"time"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"encoding/json"
	"bytes"
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

type roomDescription struct {
	emptyList map[string]bool
	wallList  map[string]bool
	available bool
	connected map[int]bool
	id        int
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

var roomsDescription []roomDescription
var notConnectedRooms []roomDescription
var notConnectedRoomsCount int

func main() {
	var start = time.Now()
	var fieldMap = rooms(100, 100)
	var elapsed = time.Since(start)
	fmt.Printf("go execution %v", elapsed)

	fmt.Println("")
	fmt.Println("")
	var field, err = json.Marshal(fieldMap)

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
	var field = make([][]string, rows+1)
	var item string
	var subField []string
	var x, y int

	for y = 0; y <= rows; y++ {
		subField = make([]string, cells+1)
		field[y] = subField

		for x = 0; x <= cells; x++ {
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
	var newRoomCoordinates point

	for !generated {
		newRoomCoordinates = getRandomRoomSize()
		cursorRectangle = drawRoom(cursorRectangle, newRoomCoordinates)

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
			occupiedCells[concatYtoX(y, x)] = true
		}
	}

	for y := startY; y <= roomEndY; y++ {
		if !wallsIntersect(y, cursorRectangle.x) {
			fieldMap[y][cursorRectangle.x] = types["rock"]
			occupiedWalls[concatYtoX(y, cursorRectangle.x)] = true
		}

		if !wallsIntersect(y, roomEndX) {
			fieldMap[y][roomEndX] = types["rock"]
			occupiedWalls[concatYtoX(y, roomEndX)] = true
		}
	}

	for x := cursorRectangle.x; x <= roomEndX; x++ {
		if !wallsIntersect(startY, x) {
			fieldMap[startY][x] = types["rock"]
			occupiedWalls[concatYtoX(startY, roomEndX)] = true
		}

		if !wallsIntersect(roomEndY, x) {
			fieldMap[roomEndY][x] = types["rock"]
			occupiedWalls[concatYtoX(roomEndY, roomEndX)] = true
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
		if crossedRoom.occupiedCells[concatYtoX(y, x)] {
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

func findInt(array []int, value int) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

// tested
func getWallsInRange(room room, startX int, endX int) []string {
	var keys = getKeys(room.occupiedWalls)
	var allCells, filteredCells []string
	var xCoordinate, yCoordinate int

	for _, cell := range keys {
		xCoordinate, _ = strconv.Atoi(strings.Split(cell, ":")[1])

		if xCoordinate >= startX && xCoordinate <= endX {
			allCells = append(allCells, cell)
		}
	}

	for _, cell := range allCells {
		var coordinates = strings.Split(cell, ":")
		yCoordinate, _ = strconv.Atoi(coordinates[0])
		xCoordinate, _ = strconv.Atoi(coordinates[1])

		if find(allCells, concatYtoX(yCoordinate+1, xCoordinate)) == false {
			filteredCells = append(filteredCells, cell)
		}
	}

	return filteredCells
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
	var wallCoordinates []string
	var yCoordinate int
	var deleteIt bool
	var wallToCompareCoordinates []string
	var yCoordinateCompare int

	for _, wall := range wallsInRange {
		deleteIt = false
		wallCoordinates = strings.Split(wall, ":")
		yCoordinate, _ = strconv.Atoi(wallCoordinates[0])

		for _, wallToCompare := range wallsInRange {
			wallToCompareCoordinates = strings.Split(wallToCompare, ":")

			if wallCoordinates[1] == wallToCompareCoordinates[1] {
				yCoordinateCompare, _ = strconv.Atoi(wallToCompareCoordinates[0])

				if yCoordinate < yCoordinateCompare {
					deleteIt = true
				}
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

func clearFatWalls() {
	for y := 1; y < fieldSize["maxY"]; y++ {
	xFor:
		for x := 1; x < fieldSize["maxX"]; x++ {
			if fieldMap[y][x] != types["rock"] {
				continue xFor
			}

			if canDelete(y, x) {
				fieldMap[y][x] = types["empty"]
			}
		}
	}
}

func IsInListOfDeletion(mapP string) bool {
	var trueVariants = []string{
		// corner
		"001011000",
		"100110000",
		"000110100",
		"000011001",
		// half-rock
		"011011011",
		"111111000",
		"110110110",
		"000111111",
		// chair-rock
		"001011011",
		"100110110",
		"011011001",
		"110110100",
		"000011111",
		"000110111",
		"111011000",
		"111110000",
		// little tetris
		"001011001",
		"100110100",
		"000010111",
		"111010000",
		//boot
		"111011011",
		"111111100",
		"110110111",
		"001111111",
		"011011111",
		"111111001",
		"111110110",
		"100111111",
		// crakoz9bra
		"111011111",
		"111111101",
		"111110111",
		"101111111",
		// stairs
		"001011111",
		"100110111",
		"111110100",
		"111011001",
		// full
		"111111111",}

	return find(trueVariants, mapP) == true
}

func canDelete(y, x int) bool {
	return IsInListOfDeletion(getMapOfPointPosition(y, x))
}

func getMapOfPointPosition(y, x int) string {
	var buffer bytes.Buffer

	for yCheck := y - 1; yCheck <= y+1; yCheck++ {
		for xCheck := x - 1; xCheck <= x+1; xCheck++ {
			if fieldMap[yCheck][xCheck] == types["rock"] {
				buffer.WriteString("1")
			} else {
				buffer.WriteString("0")
			}
		}
	}

	var a = buffer.String()
	return a
}

func drawDoors() {
	var id = 0

	for y := 1; y < fieldSize["maxY"]; y++ {
		for x := 1; x < fieldSize["maxX"]; x++ {
			if fieldMap[y][x] == types["rock"] {
				continue
			}

			var hasDescription, roomDescription = getRoomDescription(y, x)

			if hasDescription {
				roomDescription.id = id
				id++
				roomsDescription = append(roomsDescription, roomDescription)
			}
		}
	}

	roomsDescription[0].available = true
	notConnectedRooms = roomsDescription
	notConnectedRoomsCount = len(roomsDescription) - 1

	var i = 0
	for notConnectedRoomsCount != 0 {
		connectNotConnectedRooms()
		i++
	}
}

func getRoomDescription(y, x int) (bool, roomDescription) {
	for _, room := range roomsDescription {
		if room.emptyList[concatYtoX(y, x)] {
			return false, roomDescription{}
		}
	}

	return true, generateRoomDescription(y, x)
}

func generateRoomDescription(y, x int) roomDescription {
	var roomDesc = roomDescription{
		emptyList: map[string]bool{},
		wallList:  map[string]bool{},
		connected: map[int]bool{},
	}

	roomDesc.emptyList[concatYtoX(y, x)] = true
	var toCheck = []point{{y: y, x: x}}
	var i = 0

	for len(toCheck) != 0 && i < 250 {
		var newToCheck []point

		for _, fieldToCheck := range toCheck {
			i++
			y = fieldToCheck.y
			x = fieldToCheck.x

			if fieldMap[y-1][x] == types["rock"] {
				roomDesc.wallList[concatYtoX(y-1, x)] = true
			}

			if fieldMap[y][x+1] == types["empty"] {
				if !roomDesc.emptyList[concatYtoX(y, x+1)] {
					newToCheck = append(newToCheck, point{y: y, x: x + 1})
				}

				roomDesc.emptyList[concatYtoX(y, x+1)] = true
			} else {
				roomDesc.wallList[concatYtoX(y, x+1)] = true
			}

			if fieldMap[y+1][x] == types["empty"] {
				if !roomDesc.emptyList[concatYtoX(y+1, x)] {
					newToCheck = append(newToCheck, point{y: y + 1, x: x})
				}

				roomDesc.emptyList[concatYtoX(y+1, x)] = true
			} else {
				roomDesc.wallList[concatYtoX(y+1, x)] = true
			}

			if fieldMap[y][x-1] == types["empty"] {
				if !roomDesc.emptyList[concatYtoX(y, x-1)] {
					newToCheck = append(newToCheck, point{y: y, x: x - 1})
				}

				roomDesc.emptyList[concatYtoX(y, x-1)] = true
			} else {
				roomDesc.wallList[concatYtoX(y, x-1)] = true
			}
		}

		toCheck = newToCheck
		newToCheck = []point{}
	}

	return roomDesc
}

func connectNotConnectedRooms() {
	for roomI := range notConnectedRooms {
		if len(notConnectedRooms)-1 < roomI {
			return
		}

		var room = &notConnectedRooms[roomI]
		var connectVariants = getConnectVariants(room)
		var connectNumber int

		if len(connectVariants) > 0 {
			connectNumber = getRandomInt(0, len(connectVariants))
		} else {
			continue
		}

		var toConnectI = connectVariants[connectNumber]
		var toConnect = &roomsDescription[toConnectI]

		if room.available || toConnect.available {
			room.available = true
			toConnect.available = true
		}

		connectTwoRooms(room, toConnect)
		recalculateConnectionDiff(room)
	}
}

func getConnectVariants(room *roomDescription) []int {
	var roomWalls = getKeys(room.wallList)
	var variants []int

	for _, roomWall := range roomWalls {
		for roomVariantI, roomVariant := range roomsDescription {
			if roomVariant.id != room.id && roomVariant.wallList[roomWall] && findInt(variants, roomVariantI) == false {
				variants = append(variants, roomVariantI)
			}
		}
	}

	return variants
}

func connectTwoRooms(from, to *roomDescription) {
	var jointWalls []string

	if from.connected[to.id] {
		return
	}

	from.connected[to.id] = true
	to.connected[from.id] = true

	for fromWall := range from.wallList {
		if to.wallList[fromWall] {
			jointWalls = append(jointWalls, fromWall)
		}
	}

	var randomWall = jointWalls[getRandomInt(0, len(jointWalls))]
	var coordinates = strings.Split(randomWall, ":")
	var y, _ = strconv.Atoi(coordinates[0])
	var x, _ = strconv.Atoi(coordinates[1])
	var coordinates2 = point{y: y, x: x}

	if coordinates2.y == fieldSize["maxY"] || coordinates2.y == 0 || coordinates2.x == fieldSize["maxX"] || coordinates2.x == 0 {
		return
	}

	fieldMap[coordinates2.y][coordinates2.x] = types["empty"]
}

func recalculateConnectionDiff(firstConnectedRoom *roomDescription) {
	var toCheck = []*roomDescription{firstConnectedRoom}

	for len(toCheck) != 0 {
		var newToCheck []*roomDescription

		//for _, connectedRoom := range toCheck {
		for roomId := range toCheck {
			var room = toCheck[roomId]

			if len(room.connected) == 0 || room.available {
				continue
			}

			room.available = true

		xFor:
			for roomIdRelation := range room.connected {
				var roomRelation = &roomsDescription[roomIdRelation]

				if roomRelation.available {
					continue xFor
				}
				//roomRelation.available = true

				newToCheck = append(newToCheck, roomRelation)
			}
		}
		//}

		toCheck = newToCheck
	}

	var filteredNotConnectedRooms []roomDescription
	for _, r := range roomsDescription {
		if !r.available {
			filteredNotConnectedRooms = append(filteredNotConnectedRooms, r)
		}
	}
	notConnectedRooms = filteredNotConnectedRooms
	notConnectedRoomsCount = len(notConnectedRooms)
}

func concatYtoX(y, x int) string {
	return strconv.Itoa(y) + ":" + strconv.Itoa(x)
}
