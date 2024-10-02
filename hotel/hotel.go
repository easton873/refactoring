package hotel

import "errors"

type Hotel struct {
	peopleStaying []string
	emptyRooms    []int
	occupiedRooms map[int]string
	dirtyRooms    []int
	employees     []int

	money int
}

func (h *Hotel) getEmptyRoom() (int, bool) {
	if len(h.emptyRooms) == 0 {
		return 0, false
	}
	return h.emptyRooms[0], true
}

func (h *Hotel) checkIn(name string, room int) {
	for i := range h.emptyRooms {
		if h.emptyRooms[i] == room {
			h.emptyRooms = append(h.emptyRooms[:i], h.emptyRooms[i+1:]...)
			h.occupiedRooms[room] = name
			h.peopleStaying = append(h.peopleStaying, name)
			h.money += len(h.occupiedRooms)
			break
		}
	}
}

func (h *Hotel) roomLookup(name string) (int, bool) {
	for k, v := range h.occupiedRooms {
		if v == name {
			return k, true
		}
	}
	return 0, false
}

func (h *Hotel) checkOut(room int) {
	for i, peopleStaying := range h.peopleStaying {
		if peopleStaying == h.occupiedRooms[room] {
			h.peopleStaying = append(h.peopleStaying[:i], h.peopleStaying[i+1:]...)
			h.dirtyRooms = append(h.dirtyRooms, room)
			delete(h.occupiedRooms, room)
			break
		}
	}
}

func (h *Hotel) cleanRoom(room int, employeeIndex int) {
	if h.employees[employeeIndex] <= 0 {
		return
	}
	for i, r := range h.dirtyRooms {
		if r == room {
			h.dirtyRooms = append(h.dirtyRooms[:i], h.dirtyRooms[i+1:]...)
			h.emptyRooms = append(h.emptyRooms, room)
			h.employees[employeeIndex]--
			break
		}
	}
}

func (h *Hotel) hireEmployee(moneyIn int) {
	if h.money < moneyIn || moneyIn < 0 {
		return
	}
	h.money -= moneyIn
	h.employees = append(h.employees, moneyIn)
}

func (h *Hotel) payEmployee(money int, index int) {
	if h.money < money || money < 0 {
		return
	}
	h.money -= money
	hotelWorkRequired := float64(len(h.emptyRooms)) / float64(len(h.emptyRooms)+len(h.dirtyRooms)+len(h.occupiedRooms))
	h.employees[index] += int(float64(money) * hotelWorkRequired)
}

func (h *Hotel) getHappiestEmployeeIndex() (int, bool) {
	if len(h.employees) == 0 {
		return 0, false
	}
	happiest := 0
	for i := 1; i < len(h.employees); i++ {
		if h.employees[i] > h.employees[happiest] {
			happiest = i
		}
	}
	return happiest, true
}

func (h *Hotel) buyRoom(room int) error {
	for _, r := range h.emptyRooms {
		if r == room {
			return errors.New("room already built")
		}
	}
	for _, r := range h.dirtyRooms {
		if r == room {
			return errors.New("room already built")
		}
	}
	for r := range h.occupiedRooms {
		if r == room {
			return errors.New("room already built")
		}
	}
	total := len(h.emptyRooms) + len(h.dirtyRooms) + len(h.occupiedRooms)
	if h.money < total {
		return errors.New("not enough money")
	}
	h.money -= total
	h.emptyRooms = append(h.emptyRooms, room)
	return nil
}
