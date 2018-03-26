"use strict";
const type = {
    rock: '■',
    human: 'o',
    tree: '֏',
    empty: ' ',
    track: 'x',
    route: 'w'
};

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}


class MapGenerator {

    constructor() {
        this.createdRooms = [];
        this.generated = false;
        this.minYOnRow = -1;
    }

    empty(rows, cells) {
        let field = [];
        let item;

        for (let y = 0; y <= rows; y++) {
            field[y] = [];

            for (let x = 0; x <= cells; x++) {
                if (y === 0 || y === rows || x === 0 || x === cells) {
                    item = type.rock;
                } else {
                    item = type.empty;
                }

                field[y][x] = item;
            }
        }


        return field;
    }

    rooms(rows, cells) {
        this.minRoomSize = {x: 4, y: 4};
        this.maxRoomSize = {x: 9, y: 9};
        this.fieldMap = this.empty(rows, cells);
        this.fieldSize = {maxY: rows, maxX: cells};
        console.time('rooms');
        this.drawRooms();
        this.clearFatWalls();
        this.drawDoors();
        console.timeEnd('rooms');
        this.fieldMap[1][1] = type.human;
        this.fieldMap[1][this.fieldMap[0].length - 2] = type.human;
        this.fieldMap[this.fieldMap.length - 2][this.fieldMap[0].length - 2] = type.human;
        this.fieldMap[this.fieldMap.length - 2][1] = type.human;

        this.clearMemory();
        return this.fieldMap;
    }

    clearMemory() {
        this.createdRooms = null;
        this.roomsDescription = null;
        this.notConnectedRooms = null;
    }

    drawDoors() {
        this.roomsDescription = [];

        let id = 0;
        for (let y = 1; y < this.fieldSize.maxY; y++) {
            for (let x = 1; x < this.fieldSize.maxX; x++) {
                if (this.fieldMap[y][x] === type.rock) {
                    continue;
                }

                const roomDescription = this.getRoomDescription(y, x);
                if (roomDescription) {
                    roomDescription.id = id;
                    id++;
                    this.roomsDescription.push(roomDescription);
                }
            }
        }

        this.roomsDescription[0].available = true;
        this.notConnectedRooms = this.roomsDescription;
        this.notConnectedRoomsCount = this.roomsDescription.length - 1;

        let i = 0;
        while (this.notConnectedRoomsCount !== 0) {
            this.connectNotConnectedRooms(i);
            i++;
        }
    }

    connectNotConnectedRooms(iteration) {
        for (const room of this.notConnectedRooms) {
            const connectVariants = this.getConnectVariants(room);
            const connectNumber = getRandomInt(1, connectVariants.length - 1);
            let toConnect = connectVariants[connectNumber];

            if (room.available || toConnect.available) {
                room.available = true;
                toConnect.available = true;
            }

            this.connectTwoRooms(room, toConnect);
            this.recalCulateConnectionDiff(room);
        }
    }

    recalCulateConnectionDiff(connectedRoom) {
        let toCheck = [connectedRoom];

        while (toCheck.length !== 0) {
            const newToCheck = [];

            for (const roomId in connectedRoom.connected) {
                const room = this.roomsDescription[roomId];

                if (!room.connected || room.available) {
                    continue;
                }

                room.available = true;

                for (const roomIdRelation in room.connected) {
                    const roomRelation = this.roomsDescription[roomIdRelation];

                    if (roomRelation.available) {
                        continue;
                    }

                    newToCheck.push(roomRelation);
                }
            }

            toCheck = newToCheck;
        }

        this.notConnectedRooms = this.roomsDescription.filter(r => !r.available);
        this.notConnectedRoomsCount = this.notConnectedRooms.length;
    }

    connectTwoRooms(from, to) {
        const jointWalls = [];
        from.connected = from.connected || [];
        to.connected = to.connected || [];

        if (from.connected[to.id]) {
            return;
        }

        from.connected[to.id] = true;
        to.connected[from.id] = true;

        for (const fromWall in from.wallList) {
            if (to.wallList[fromWall]) {
                jointWalls.push(fromWall);
            }
        }

        const randomWall = jointWalls[getRandomInt(0, jointWalls.length - 1)];
        let coordinates = randomWall.split(':');
        coordinates = {y: +coordinates[0], x: +coordinates[1]};

        if (coordinates.y === this.fieldSize.maxY || coordinates.y === 0 || coordinates.x === this.fieldSize.maxX || coordinates.x === 0) {
            return;
        }

        this.fieldMap[coordinates.y][coordinates.x] = type.empty;
    }

    getConnectVariants(room) {
        const roomWalls = Object.keys(room.wallList);
        const variants = [];

        for (const roomWall of roomWalls) {
            for (const roomVariant of this.roomsDescription) {
                if (roomVariant.id !== room.id && roomVariant.wallList[roomWall]) {
                    variants.push(roomVariant);
                }
            }
        }

        return variants;
    }

    getRoomDescription(y, x) {
        for (const room of this.roomsDescription) {
            if (room.emptyList[`${y}:${x}`]) {
                return false;
            }
        }

        return this.generateRoomDescription(y, x);
    }

    generateRoomDescription(y, x) {
        const room = {
            emptyList: {},
            wallList: {}
        };
        room.emptyList[`${y}:${x}`] = true;
        let toCheck = [{y: y, x: x}];
        let i = 0;

        while (toCheck.length !== 0 && i < 250) {
            let newToCheck = [];

            for (const fieldToCheck of toCheck) {
                i++;
                y = fieldToCheck.y;
                x = fieldToCheck.x;

                if (this.fieldMap[y - 1][x] === type.rock) {
                    room.wallList[`${y - 1}:${x}`] = true;
                }

                if (this.fieldMap[y][x + 1] === type.empty) {
                    if (!room.emptyList[`${y}:${x + 1}`]) {
                        newToCheck.push({y: y, x: x + 1});
                    }

                    room.emptyList[`${y}:${x + 1}`] = true;
                } else {
                    room.wallList[`${y}:${x + 1}`] = true;
                }


                if (this.fieldMap[y + 1][x] === type.empty) {
                    if (!room.emptyList[`${y + 1}:${x}`]) {
                        newToCheck.push({y: y + 1, x: x});
                    }
                    room.emptyList[`${y + 1}:${x}`] = true;
                } else {
                    room.wallList[`${y + 1}:${x}`] = true;
                }

                if (this.fieldMap[y][x - 1] === type.empty) {
                    if (!room.emptyList[`${y}:${x - 1}`]) {
                        newToCheck.push({y: y, x: x - 1});
                    }
                    room.emptyList[`${y}:${x - 1}`] = true;
                } else {
                    room.wallList[`${y}:${x - 1}`] = true;
                }
            }

            toCheck = newToCheck;
            newToCheck = [];
        }

        return room;
    }

    clearFatWalls() {
        for (let y = 1; y < this.fieldSize.maxY; y++) {
            xFor: for (let x = 1; x < this.fieldSize.maxX; x++) {
                if (this.fieldMap[y][x] !== type.rock) {
                    continue xFor;
                }

                if (this.canDelete(y, x)) {
                    this.fieldMap[y][x] = type.empty;
                }
            }
        }
    }

    IsInListOfDeletion(map) {
        const trueVariants = [
            // corner
            `001011000`,
            `100110000`,
            `000110100`,
            `000011001`,
            // half-rock
            `011011011`,
            `111111000`,
            `110110110`,
            `000111111`,
            // chair-rock
            `001011011`,
            `100110110`,
            `011011001`,
            `110110100`,
            `000011111`,
            `000110111`,
            `111011000`,
            `111110000`,
            // little tetris
            `001011001`,
            `100110100`,
            `000010111`,
            `111010000`,
            //boot
            `111011011`,
            `111111100`,
            `110110111`,
            `001111111`,
            `011011111`,
            `111111001`,
            `111110110`,
            `100111111`,
            // crakoz9bra
            `111011111`,
            `111111101`,
            `111110111`,
            `101111111`,
            // stairs
            `001011111`,
            `100110111`,
            `111110100`,
            `111011001`,
            // full
            `111111111`,
        ];

        return ~trueVariants.indexOf(map);
    }

    canDelete(y, x) {
        return this.IsInListOfDeletion(this.getMapOfPointPosition(y, x));
    }

    getMapOfPointPosition(y, x) {
        let map = '';

        for (let yCheck = y - 1; yCheck <= y + 1; yCheck++) {
            for (let xCheck = x - 1; xCheck <= x + 1; xCheck++) {
                if (this.fieldMap[yCheck][xCheck] === type.rock) {
                    map += '1';
                } else {
                    map += '0';
                }
            }
        }

        return map;
    }

    drawRooms() {
        this.cursorRectangle = {x: 0, y: 0};
        this.rawStart = {x: 0, y: 0};

        while (!this.generated) {
            const room = this.getRandomRoomSize();
            this.cursorRectangle = this.drawRoom(this.cursorRectangle, room);

            if (this.cursorRectangle.x === this.fieldSize.maxX) {
                this.moveCursorToTheNewRow();
                this.minYOnRow = -1;
            }
        }
    }

    moveCursorToTheNewRow() {
        if (this.rawStart.y === this.fieldSize.maxY) {
            this.generated = true;
        }

        let roomWithMaxY = this.createdRooms[0];

        for (const room of this.createdRooms) {
            if (room.start.x === 0 && room.end.y >= roomWithMaxY.end.y) {
                roomWithMaxY = room;
            }
        }

        this.cursorRectangle = {x: 0, y: roomWithMaxY.end.y};
        this.rawStart = JSON.parse(JSON.stringify(this.cursorRectangle));
    }

    drawRoom(cursorRectangle, room) {
        const roomEndX = cursorRectangle.x + room.x;
        const startY = this.getStartRowForRoom(this.cursorRectangle.x, roomEndX);
        let roomEndY = startY + room.y;

        if (roomEndY > this.fieldSize.maxY) {
            roomEndY = this.fieldSize.maxY;
        }

        if (roomEndY + 1 === this.fieldSize.maxY) {
            roomEndY = this.fieldSize.maxY;
        }

        const occupiedCells = {};
        const occupiedWalls = {};

        for (let y = startY; y <= roomEndY; y++) {
            for (let x = cursorRectangle.x; x <= roomEndX; x++) {
                occupiedCells[`${y}:${x}`] = true;
            }
        }

        for (let y = startY; y <= roomEndY; y++) {
            if (!this.wallsIntersect(y, cursorRectangle.x)) {
                this.fieldMap[y][cursorRectangle.x] = type.rock;
                occupiedWalls[`${y}:${cursorRectangle.x}`] = true;
            }

            if (!this.wallsIntersect(y, roomEndX)) {
                this.fieldMap[y][roomEndX] = type.rock;
                occupiedWalls[`${y}:${roomEndX}`] = true;
            }
        }

        for (let x = cursorRectangle.x; x <= roomEndX; x++) {
            if (!this.wallsIntersect(startY, x)) {
                this.fieldMap[startY][x] = type.rock;
                occupiedWalls[`${startY}:${x}`] = true;
            }

            if (!this.wallsIntersect(roomEndY, x)) {
                this.fieldMap[roomEndY][x] = type.rock;
                occupiedWalls[`${roomEndY}:${x}`] = true;
            }
        }

        this.createdRooms.push({
            start: {x: cursorRectangle.x, y: startY},
            end: {x: roomEndX, y: roomEndY},
            occupiedCells: occupiedCells,
            occupiedWalls: occupiedWalls
        });

        if (this.minYOnRow > roomEndY || this.minYOnRow === -1) {
            this.minYOnRow = roomEndY;
        }


        return {y: startY, x: roomEndX};
    }

    wallsIntersect(y, x) {
        for (const crossedRoom of this.crossedRooms) {
            if (crossedRoom.occupiedCells[`${y}:${x}`]) {
                return true;
            }
        }

        return false;
    }

    getStartRowForRoom(startX, roomEndX) {
        let startY = this.cursorRectangle.y;
        this.crossedRooms = [];

        if (this.rawStart.y !== 0) {
            let wallsInRange = [];

            // We must cling to the "lowest" room in our range
            for (const room of this.createdRooms) {
                if ((room.start.x >= startX && room.start.x <= roomEndX) || (room.end.x >= startX && room.end.x <= roomEndX)) {
                    wallsInRange = wallsInRange.concat(this.getWallsInRange(room, this.cursorRectangle.x, roomEndX));
                }
            }

            wallsInRange = this.deleteDuplications(wallsInRange);

            if (wallsInRange.length) {
                startY = this.findMinY(wallsInRange);
            }

            for (const room of this.createdRooms) {
                if ((room.end.x >= this.cursorRectangle.x || room.start.x <= roomEndX) && room.end.y >= startY) {
                    this.crossedRooms.push(room);
                }
            }
        }

        return startY;
    }

    deleteDuplications(wallsInRange) {
        const filtered = [];

        for (const wall of wallsInRange) {
            let deleteIt = false;
            const wallCoordinates = wall.split(':');

            for (const wallToCompare of wallsInRange) {
                const wallToCompareCoordinates = wallToCompare.split(':');
                if (+wallCoordinates[1] === +wallToCompareCoordinates[1] && +wallCoordinates[0] < +wallToCompareCoordinates[0]) {
                    deleteIt = true;
                }
            }

            if (!deleteIt) {
                filtered.push(wall);
            }
        }

        return filtered;
    }

    findMinY(wallsInRange) {
        let minY = +wallsInRange[0].split(':')[0];

        for (const wall of wallsInRange) {
            const wallCoordinates = wall.split(':');

            if (+wallCoordinates[0] < minY) {
                minY = +wallCoordinates[0];
            }
        }

        return minY;
    }

    getWallsInRange(room, startX, endX) {
        const keys = Object.keys(room.occupiedWalls);

        const allCells = keys.filter(v => {
            const xCoordinate = +v.split(':')[1];
            return xCoordinate >= startX && xCoordinate <= endX;
        });

        return allCells.filter(v => {
            const coordinates = v.split(':');
            return allCells.indexOf((+coordinates[0] + 1) + ':' + coordinates[1]) === -1;
        });
    }

    getRandomRoomSize() {
        let x = getRandomInt(this.minRoomSize.x, this.maxRoomSize.x);
        let y = getRandomInt(this.minRoomSize.y, this.maxRoomSize.y);

        if (this.cursorRectangle.x + x > this.fieldSize.maxX) {
            x = this.fieldSize.maxX - this.cursorRectangle.x;
        }

        if (this.cursorRectangle.y + y > this.fieldSize.maxY) {
            y = this.fieldSize.maxY - this.cursorRectangle.y;
        }

        if (this.cursorRectangle.x + x === this.fieldSize.maxX - 1) {
            x++;
        }

        return {
            x: x,
            y: y
        }
    }
}

const mapGenerator = new MapGenerator();

console.time('node execution');
mapGenerator.empty(10, 10);
console.timeEnd('node execution');