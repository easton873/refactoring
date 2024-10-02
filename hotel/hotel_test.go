package hotel

import (
	"fmt"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestHotelFixture(t *testing.T) {
	gunit.Run(new(HotelFixture), t)
}

type HotelFixture struct {
	*gunit.Fixture
	Hotel
}

func (this *HotelFixture) Setup() {
	this.money = 10
	this.emptyRooms = []int{1, 2, 3, 4, 5}
	this.employees = []int{3, 3}
	this.occupiedRooms = make(map[int]string)
}

func (this *HotelFixture) TestBasics() {
	this.checkIn("joe", 1)
	this.So(len(this.emptyRooms), should.Equal, 4)
	this.So(len(this.peopleStaying), should.Equal, 1)
	this.So(this.money, should.Equal, 11)

	this.checkIn("joe", 1)
	this.So(len(this.emptyRooms), should.Equal, 4)
	this.So(len(this.peopleStaying), should.Equal, 1)
	this.So(this.money, should.Equal, 11)

	this.checkOut(1)
	this.So(len(this.emptyRooms), should.Equal, 4)
	this.So(len(this.dirtyRooms), should.Equal, 1)
	this.So(len(this.peopleStaying), should.Equal, 0)

	this.cleanRoom(1, 0)
	this.So(len(this.dirtyRooms), should.Equal, 0)
	this.So(this.employees[0], should.Equal, 2)
	this.So(len(this.emptyRooms), should.Equal, 5)

	this.payEmployee(1, 0)
	this.So(this.employees[0], should.Equal, 3)
	this.So(this.money, should.Equal, 10)

	this.payEmployee(11, 0)
	this.So(this.employees[0], should.Equal, 3)
	this.So(this.money, should.Equal, 10)
}

func (this *HotelFixture) TestGetEmptyRoom() {
	for room, found := this.getEmptyRoom(); found; room, found = this.getEmptyRoom() {
		this.checkIn("test", room)
	}

	room, found := this.getEmptyRoom()
	this.So(room, should.Equal, 0)
	this.So(found, should.BeFalse)
}

func (this *HotelFixture) TestMoney() {
	this.checkIn("test1", 1)
	this.checkIn("test2", 2)
	this.checkIn("test3", 3)
	this.checkIn("test4", 4)
	this.checkIn("test5", 5)

	this.So(this.money, should.Equal, 25)
	this.So(len(this.occupiedRooms), should.Equal, 5)
	this.So(len(this.emptyRooms), should.Equal, 0)
	this.So(len(this.dirtyRooms), should.Equal, 0)

	this.payEmployee(10, 0)
	this.So(this.employees[0], should.Equal, 3)

	this.checkOut(1)
	this.checkOut(2)
	this.checkOut(3)

	this.So(len(this.occupiedRooms), should.Equal, 2)
	this.So(len(this.emptyRooms), should.Equal, 0)
	this.So(len(this.dirtyRooms), should.Equal, 3)

	this.cleanRoom(1, 0)
	this.cleanRoom(2, 0)

	this.So(len(this.occupiedRooms), should.Equal, 2)
	this.So(len(this.emptyRooms), should.Equal, 2)
	this.So(len(this.dirtyRooms), should.Equal, 1)

	this.So(this.employees[0], should.Equal, 1)
	this.payEmployee(5, 0)
	this.So(this.employees[0], should.Equal, 3)

	this.cleanRoom(3, 0)
	this.So(this.employees[0], should.Equal, 2)
	this.payEmployee(5, 0)
	this.So(this.employees[0], should.Equal, 5)
}

func (this *HotelFixture) TestExpandHotel() {
	for room, found := this.getEmptyRoom(); found; room, found = this.getEmptyRoom() {
		this.checkIn("test", room)
	}
	err := this.buyRoom(1)
	this.So(err.Error(), should.Equal, "room already built")
	this.So(len(this.emptyRooms), should.Equal, 0)
	this.So(this.money, should.Equal, 25)

	err = this.buyRoom(10)
	this.So(err, should.BeNil)
	this.So(len(this.emptyRooms), should.Equal, 1)
	this.So(this.money, should.Equal, 20)

	err = this.buyRoom(9)
	this.So(err, should.BeNil)
	this.So(len(this.emptyRooms), should.Equal, 2)
	this.So(this.money, should.Equal, 14)

	this.checkIn("test", 10)
	this.checkIn("test", 9)
	this.So(this.money, should.Equal, 27)

	for i := 1; i <= 10; i++ {
		this.checkOut(i)
	}
	this.So(len(this.emptyRooms), should.Equal, 0)
	this.So(len(this.dirtyRooms), should.Equal, 7)
	this.So(len(this.occupiedRooms), should.Equal, 0)

	this.hireEmployee(2)
	happiestEmployeeIndex, found := this.getHappiestEmployeeIndex()
	this.So(found, should.BeTrue)
	this.cleanRoom(1, happiestEmployeeIndex)
	this.So(this.employees[0], should.Equal, 2)
	this.So(len(this.emptyRooms), should.Equal, 1)
	this.So(len(this.dirtyRooms), should.Equal, 6)
	this.So(len(this.occupiedRooms), should.Equal, 0)
	for i := 2; i <= 9; i++ {
		happiestEmployeeIndex, found = this.getHappiestEmployeeIndex()
		this.cleanRoom(i, happiestEmployeeIndex)
	}
	this.So(len(this.emptyRooms), should.Equal, 6)
	this.So(len(this.dirtyRooms), should.Equal, 1)
	this.So(len(this.occupiedRooms), should.Equal, 0)
	this.So(this.employees[0], should.Equal, 0)
	this.So(this.employees[1], should.Equal, 1)
	this.So(this.employees[2], should.Equal, 1)

	this.cleanRoom(10, 0)

	this.So(len(this.emptyRooms), should.Equal, 6)
	this.So(len(this.dirtyRooms), should.Equal, 1)
	this.So(len(this.occupiedRooms), should.Equal, 0)
	this.So(this.employees[0], should.Equal, 0)
	this.So(this.employees[1], should.Equal, 1)
	this.So(this.employees[2], should.Equal, 1)

	this.cleanRoom(10, 2)
	this.So(len(this.emptyRooms), should.Equal, 7)
	this.So(len(this.dirtyRooms), should.Equal, 0)
	this.So(len(this.occupiedRooms), should.Equal, 0)
	this.So(this.employees[0], should.Equal, 0)
	this.So(this.employees[1], should.Equal, 1)
	this.So(this.employees[2], should.Equal, 0)

	hotel.money = 0
	err = hotel.buyRoom(11)
	this.So(err.Error(), should.Equal, "not enough money")
}

func (this *HotelFixture) TestRoomLookup() {
	room, found := this.roomLookup("Joe")
	this.So(found, should.BeFalse)
	this.So(room, should.Equal, 0)

	this.checkIn("1", 1)
	room, found = this.roomLookup("1")
	this.So(found, should.BeTrue)
	this.So(room, should.Equal, 1)

	this.checkIn("2", 2)
	room, found = this.roomLookup("1")
	this.So(found, should.BeTrue)
	this.So(room, should.Equal, 1)

	room, found = this.roomLookup("2")
	this.So(found, should.BeTrue)
	this.So(room, should.Equal, 2)

	for room, found = this.getEmptyRoom(); found; room, found = this.getEmptyRoom() {
		this.checkIn(fmt.Sprintf("%d", room), room)
	}

	this.So(len(this.occupiedRooms), should.Equal, 5)
	for i := 1; i <= 5; i++ {
		room, found = this.roomLookup(fmt.Sprintf("%d", i))
		this.So(found, should.BeTrue)
		this.So(room, should.Equal, i)
	}
}

func (this *HotelFixture) TestNegativeWages() {
	this.So(len(this.employees), should.Equal, 2)
	this.hireEmployee(-1)
	this.So(len(this.employees), should.Equal, 2)

	this.So(this.employees[0], should.Equal, 3)
	this.payEmployee(-1, 0)
	this.So(this.employees[0], should.Equal, 3)
}
