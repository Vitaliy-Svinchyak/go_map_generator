package main

import (
	"testing"
	"reflect"
)

func TestFindMinY(t *testing.T) {
	v := findMinY([]string{"24:25", "24:26", "23:24", "30:27", "30:30", "30:28", "30:29", "30:23",})

	if v != 23 {
		t.Error("Failed: Expected 23, got ", v)
	}
}

func TestDeleteDuplications(t *testing.T) {
	v := deleteDuplications([]string{"9:21", "0:15", "9:15", "0:16", "9:16", "0:17", "9:17", "0:18", "9:18",
		"0:19", "9:19", "0:20", "9:20", "5:21", "0:22", "5:22", "0:23", "5:23", "15:18", "15:15", "15:16",
		"15:17", "13:23", "13:19", "13:20", "13:21", "13:22", "20:19", "20:15", "20:16", "20:17", "20:18",
		"22:19", "22:23", "22:20", "22:21", "22:22", "28:17", "28:15", "28:16", "28:18", "28:19",
		"28:20", "28:21", "28:22", "28:23", "30:18", "30:15", "30:16", "30:17", "30:23",})

	var expected = []string{"28:19", "28:20", "28:21", "28:22", "30:18",
		"30:15", "30:16", "30:17", "30:23",}
	if !reflect.DeepEqual(v, expected) {
		t.Error("Failed: Expected", expected, "\n got", v)
	}
}

func TestEmpty(t *testing.T) {
	v := empty(2, 2)
	var expected = [][]string{
		0: {
			0: types["rock"],
			1: types["rock"],
			2: types["rock"],
		},
		1: {
			0: types["rock"],
			1: types["empty"],
			2: types["rock"],
		},
		2: {
			0: types["rock"],
			1: types["rock"],
			2: types["rock"],
		},
	}

	if !reflect.DeepEqual(v, expected) {
		t.Error("Failed: Expected", expected, "\n got", v)
	}

}

func TestWallsIntersect(t *testing.T) {
	v := wallsIntersect(5, 6)

	if v != false {
		t.Error("Failed: Expected false, got ", v)
	}

	crossedRooms = append(crossedRooms,
		room{start: point{x: 0, y: 0},
			end: point{x: 6, y: 5},
			occupiedCells:
			map[string]bool{"0:0": true, "0:1": true, "0:2": true, "0:3": true, "0:4": true, "0:5": true, "0:6": true,
				"1:0": true, "1:1": true, "1:2": true, "1:3": true, "1:4": true, "1:5": true, "1:6": true, "2:0": true,
				"2:1": true, "2:2": true, "2:3": true, "2:4": true, "2:5": true, "2:6": true, "3:0": true, "3:1": true,
				"3:2": true, "3:3": true, "3:4": true, "3:5": true, "3:6": true, "4:0": true, "4:1": true, "4:2": true,
				"4:3": true, "4:4": true, "4:5": true, "4:6": true, "5:0": true, "5:1": true, "5:2": true, "5:3": true,
				"5:4": true, "5:5": true, "5:6": true,},
			occupiedWalls: map[string]bool{"0:0": true, "0:6": true, "1:0": true, "1:6": true, "2:0": true,
				"2:6": true, "3:0": true, "3:6": true, "4:0": true, "4:6": true, "5:0": true, "5:6": true, "0:1": true,
				"5:1": true, "0:2": true, "5:2": true, "0:3": true, "5:3": true, "0:4": true, "5:4": true, "0:5": true,
				"5:5": true,}})

	v = wallsIntersect(5, 6)

	if v != true {
		t.Error("Failed: Expected true, got ", v)
	}
}

func TestMoveCursorToTheNewRow(t *testing.T) {
	rawStart = point{x: 0, y: 0}
	fieldSize["maxY"] = 10
	fieldSize["maxX"] = 10

	createdRooms = append(createdRooms, room{start: point{x: 0, y: 0},
		end: point{x: 5, y: 5},
		occupiedCells: map[string]bool{"0:0": true, "0:1": true, "0:2": true, "0:3": true, "0:4": true, "0:5": true,
			"1:0": true, "1:1": true, "1:2": true, "1:3": true, "1:4": true, "1:5": true, "2:0": true, "2:1": true,
			"2:2": true, "2:3": true, "2:4": true, "2:5": true, "3:0": true, "3:1": true, "3:2": true, "3:3": true,
			"3:4": true, "3:5": true, "4:0": true, "4:1": true, "4:2": true, "4:3": true, "4:4": true, "4:5": true,
			"5:0": true, "5:1": true, "5:2": true, "5:3": true, "5:4": true, "5:5": true},
		occupiedWalls: map[string]bool{"0:0": true, "0:5": true, "1:0": true, "1:5": true, "2:0": true, "2:5": true,
			"3:0": true, "3:5": true, "4:0": true, "4:5": true, "5:0": true, "5:5": true, "0:1": true, "5:1": true,
			"0:2": true, "5:2": true, "0:3": true, "5:3": true, "0:4": true, "5:4": true}})

	createdRooms = append(createdRooms, room{start: point{x: 5, y: 0},
		end: point{x: 10, y: 5},
		occupiedCells: map[string]bool{"0:5": true, "0:6": true, "0:7": true, "0:8": true, "0:9": true, "0:10": true, "1:5": true,
			"1:6": true, "1:7": true, "1:8": true, "1:9": true, "1:10": true, "2:5": true, "2:6": true, "2:7": true,
			"2:8": true, "2:9": true, "2:10": true, "3:5": true, "3:6": true, "3:7": true, "3:8": true, "3:9": true,
			"3:10": true, "4:5": true, "4:6": true, "4:7": true, "4:8": true, "4:9": true, "4:10": true, "5:5": true,
			"5:6": true, "5:7": true, "5:8": true, "5:9": true, "5:10": true},
		occupiedWalls: map[string]bool{"0:5": true, "0:10": true, "1:5": true, "1:10": true, "2:5": true,
			"2:10": true, "3:5": true, "3:10": true, "4:5": true, "4:10": true, "5:5": true, "5:10": true, "0:6": true,
			"5:6": true, "0:7": true, "5:7": true, "0:8": true, "5:8": true, "0:9": true, "5:9": true}})

	moveCursorToTheNewRow()

	var expectedPoint = point{x: 0, y: 5}
	if !reflect.DeepEqual(cursorRectangle, expectedPoint) {
		t.Error("Failed: Expected", expectedPoint, "\n got", cursorRectangle)
	}
}

func TestGetWallsInRange(t *testing.T) {
	v := getWallsInRange(room{start: point{x: 0, y: 0,},
		end: point{x: 7, y: 6,},
		occupiedCells: map[string]bool{"0:0": true, "0:1": true, "0:2": true, "0:3": true, "0:4": true, "0:5": true,
			"0:6": true, "0:7": true, "1:0": true, "1:1": true, "1:2": true, "1:3": true, "1:4": true, "1:5": true,
			"1:6": true, "1:7": true, "2:0": true, "2:1": true, "2:2": true, "2:3": true, "2:4": true, "2:5": true,
			"2:6": true, "2:7": true, "3:0": true, "3:1": true, "3:2": true, "3:3": true, "3:4": true, "3:5": true,
			"3:6": true, "3:7": true, "4:0": true, "4:1": true, "4:2": true, "4:3": true, "4:4": true, "4:5": true,
			"4:6": true, "4:7": true, "5:0": true, "5:1": true, "5:2": true, "5:3": true, "5:4": true, "5:5": true,
			"5:6": true, "5:7": true, "6:0": true, "6:1": true, "6:2": true, "6:3": true, "6:4": true, "6:5": true,
			"6:6": true, "6:7": true,},
		occupiedWalls: map[string]bool{"0:0": true, "0:7": true, "1:0": true, "1:7": true, "2:0": true, "2:7": true,
			"3:0": true, "3:7": true, "4:0": true, "4:7": true, "5:0": true, "5:7": true, "6:0": true, "6:7": true,
			"0:1": true, "6:1": true, "0:2": true, "6:2": true, "0:3": true, "6:3": true, "0:4": true, "6:4": true,
			"0:5": true, "6:5": true, "0:6": true, "6:6": true,}},
		0,
		5)

	var expected = []string{"6:0", "0:1", "6:1", "0:2", "6:2", "0:3", "6:3", "0:4", "6:4", "0:5", "6:5",}

	if len(v) != len(expected) {
		t.Error("Failed length: Expected", expected, "\n got", v)
	}

	for _, element := range expected {
		if find(v, element) == false {
			t.Error("Failed: Expected", expected, "\n got", v)
		}
	}
}

func TestGetStartRowForRoom(t *testing.T) {
	cursorRectangle = point{y: 4, x: 8}
	createdRooms = append(createdRooms, room{
		start: point{x: 0, y: 0},
		end:   point{x: 6, y: 10},
		occupiedCells: map[string]bool{"0:0": true, "0:1": true, "0:2": true, "0:3": true, "0:4": true, "0:5": true,
			"0:6": true, "1:0": true, "1:1": true, "1:2": true, "1:3": true, "1:4": true, "1:5": true, "1:6": true,
			"2:0": true, "2:1": true, "2:2": true, "2:3": true, "2:4": true, "2:5": true, "2:6": true, "3:0": true,
			"3:1": true, "3:2": true, "3:3": true, "3:4": true, "3:5": true, "3:6": true, "4:0": true, "4:1": true,
			"4:2": true, "4:3": true, "4:4": true, "4:5": true, "4:6": true, "5:0": true, "5:1": true, "5:2": true,
			"5:3": true, "5:4": true, "5:5": true, "5:6": true, "6:0": true, "6:1": true, "6:2": true, "6:3": true,
			"6:4": true, "6:5": true, "6:6": true, "7:0": true, "7:1": true, "7:2": true, "7:3": true, "7:4": true,
			"7:5": true, "7:6": true, "8:0": true, "8:1": true, "8:2": true, "8:3": true, "8:4": true, "8:5": true,
			"8:6": true, "9:0": true, "9:1": true, "9:2": true, "9:3": true, "9:4": true, "9:5": true, "9:6": true,
			"10:0": true, "10:1": true, "10:2": true, "10:3": true, "10:4": true, "10:5": true, "10:6": true,},
		occupiedWalls: map[string]bool{"0:0": true, "0:6": true, "1:0": true, "1:6": true, "2:0": true, "2:6": true,
			"3:0": true, "3:6": true, "4:0": true, "4:6": true, "5:0": true, "5:6": true, "6:0": true, "6:6": true,
			"7:0": true, "7:6": true, "8:0": true, "8:6": true, "9:0": true, "9:6": true, "10:0": true, "10:6": true,
			"0:1": true, "10:1": true, "0:2": true, "10:2": true, "0:3": true, "10:3": true, "0:4": true, "10:4": true,
			"0:5": true, "10:5": true,},})

	createdRooms = append(createdRooms, room{
		start: point{x: 6, y: 0},
		end:   point{x: 10, y: 4},
		occupiedCells: map[string]bool{"0:6": true, "0:7": true, "0:8": true, "0:9": true, "0:10": true, "1:6": true,
			"1:7": true, "1:8": true, "1:9": true, "1:10": true, "2:6": true, "2:7": true, "2:8": true, "2:9": true,
			"2:10": true, "3:6": true, "3:7": true, "3:8": true, "3:9": true, "3:10": true, "4:6": true, "4:7": true,
			"4:8": true, "4:9": true, "4:10": true},
		occupiedWalls: map[string]bool{"0:6": true, "0:10": true, "1:6": true, "1:10": true, "2:6": true, "2:10": true,
			"3:6": true, "3:10": true, "4:6": true, "4:10": true, "0:7": true, "4:7": true, "0:8": true, "4:8": true,
			"0:9": true, "4:9": true}})

	createdRooms = append(createdRooms, room{
		start: point{x: 0, y: 4},
		end:   point{x: 8, y: 4},
		occupiedCells:
		map[string]bool{"4:0": true, "4:1": true, "4:2": true, "4:3": true, "4:4": true, "4:5": true, "4:6": true,
			"4:7": true, "4:8": true},
		occupiedWalls: map[string]bool{}})

	v := getStartRowForRoom(8, 10)

	// fix ?!
	if v != 5 {
		t.Error("Failed: Expected 4, got ", v)
	}
}

// Benchmarks

//func BenchmarkFindMinY(b *testing.B) {
//	b.ReportAllocs()
//
//	for n := 0; n < b.N; n++ {
//		findMinY([]string{"24:25", "24:26", "23:24", "30:27", "30:30", "30:28", "30:29", "30:23",})
//	}
//}

//func BenchmarkDeleteDuplications(b *testing.B) {
//	b.ReportAllocs()
//
//	for n := 0; n < b.N; n++ {
//		deleteDuplications([]string{"9:21", "0:15", "9:15", "0:16", "9:16", "0:17", "9:17", "0:18", "9:18",
//			"0:19", "9:19", "0:20", "9:20", "5:21", "0:22", "5:22", "0:23", "5:23", "15:18", "15:15", "15:16",
//			"15:17", "13:23", "13:19", "13:20", "13:21", "13:22", "20:19", "20:15", "20:16", "20:17", "20:18",
//			"22:19", "22:23", "22:20", "22:21", "22:22", "28:17", "28:15", "28:16", "28:18", "28:19",
//			"28:20", "28:21", "28:22", "28:23", "30:18", "30:15", "30:16", "30:17", "30:23",})
//	}
//}

func BenchmarkEmpty(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		empty(10, 10)
	}
}

//func BenchmarkGetWallsInRange(b *testing.B) {
//	b.ReportAllocs()
//
//	for n := 0; n < b.N; n++ {
//		getWallsInRange(room{start: point{x: 0, y: 0,},
//			end: point{x: 7, y: 6,},
//			occupiedCells: map[string]bool{"0:0": true, "0:1": true, "0:2": true, "0:3": true, "0:4": true, "0:5": true,
//				"0:6": true, "0:7": true, "1:0": true, "1:1": true, "1:2": true, "1:3": true, "1:4": true, "1:5": true,
//				"1:6": true, "1:7": true, "2:0": true, "2:1": true, "2:2": true, "2:3": true, "2:4": true, "2:5": true,
//				"2:6": true, "2:7": true, "3:0": true, "3:1": true, "3:2": true, "3:3": true, "3:4": true, "3:5": true,
//				"3:6": true, "3:7": true, "4:0": true, "4:1": true, "4:2": true, "4:3": true, "4:4": true, "4:5": true,
//				"4:6": true, "4:7": true, "5:0": true, "5:1": true, "5:2": true, "5:3": true, "5:4": true, "5:5": true,
//				"5:6": true, "5:7": true, "6:0": true, "6:1": true, "6:2": true, "6:3": true, "6:4": true, "6:5": true,
//				"6:6": true, "6:7": true,},
//			occupiedWalls: map[string]bool{"0:0": true, "0:7": true, "1:0": true, "1:7": true, "2:0": true, "2:7": true,
//				"3:0": true, "3:7": true, "4:0": true, "4:7": true, "5:0": true, "5:7": true, "6:0": true, "6:7": true,
//				"0:1": true, "6:1": true, "0:2": true, "6:2": true, "0:3": true, "6:3": true, "0:4": true, "6:4": true,
//				"0:5": true, "6:5": true, "0:6": true, "6:6": true,}},
//			0,
//			5)
//	}
//}

//func BenchmarkRooms(b *testing.B) {
//	b.ReportAllocs()
//
//	for n := 0; n < b.N; n++ {
//		rooms(20, 20)
//	}
//}

//func BenchmarkWallsIntersect(b *testing.B) {
//	b.ReportAllocs()
//
//	crossedRooms = append(crossedRooms,
//		room{start: point{x: 0, y: 0},
//			end: point{x: 6, y: 5},
//			occupiedCells:
//			map[string]bool{"0:0": true, "0:1": true, "0:2": true, "0:3": true, "0:4": true, "0:5": true, "0:6": true,
//				"1:0": true, "1:1": true, "1:2": true, "1:3": true, "1:4": true, "1:5": true, "1:6": true, "2:0": true,
//				"2:1": true, "2:2": true, "2:3": true, "2:4": true, "2:5": true, "2:6": true, "3:0": true, "3:1": true,
//				"3:2": true, "3:3": true, "3:4": true, "3:5": true, "3:6": true, "4:0": true, "4:1": true, "4:2": true,
//				"4:3": true, "4:4": true, "4:5": true, "4:6": true, "5:0": true, "5:1": true, "5:2": true, "5:3": true,
//				"5:4": true, "5:5": true, "5:6": true,},
//			occupiedWalls: map[string]bool{"0:0": true, "0:6": true, "1:0": true, "1:6": true, "2:0": true,
//				"2:6": true, "3:0": true, "3:6": true, "4:0": true, "4:6": true, "5:0": true, "5:6": true, "0:1": true,
//				"5:1": true, "0:2": true, "5:2": true, "0:3": true, "5:3": true, "0:4": true, "5:4": true, "0:5": true,
//				"5:5": true,}})
//
//	for n := 0; n < b.N; n++ {
//		wallsIntersect(5, 6)
//	}
//}
