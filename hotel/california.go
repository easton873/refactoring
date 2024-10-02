package hotel

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	stringScanner = &bufio.Scanner{}
	fileScanner   = &bufio.Scanner{}

	inputString = "checkin\nJimothy\ncheckin\nTerry\ncheckin\nSandy\nview\ncheckout\nTerry\nclean\n101\n0\ncheckin\nSally\nview\npay\n0\n3\nview\nbuild\n104\ncheckin\nGeorge\ncheckin\nMegan\npay\n1\n2\nview\nhire\n1\ncheckout\nSally\ncheckout\nSandy\ncheckout\nGeorge\ncheckout\nMegan\nhire\n10\nview\nquit"
	hotel       = Hotel{
		money:         10,
		emptyRooms:    []int{100, 101, 102, 201, 202, 200, 301, 302, 300},
		employees:     []int{3, 3},
		occupiedRooms: map[int]string{},
	}
	startingHow = "user"
)

func Run() {
	if Printer == nil {
		Printer = ConsolePrinter{}
	}
	how := startingHow
	if how == "user" {

	} else if how == "string" {
		initStringScanner()
	} else if how == "file" {
		initFileScanner()
	}

	for {
		Printer.Printf("Enter command (help if you don't know any): ")
		userInput := getUserInput(how)
		if userInput == "quit" {
			return
		} else if userInput == "checkin" {
			checkIn(how)
		} else if userInput == "checkout" {
			checkOut(how)
		} else if userInput == "view" {
			view()
		} else if userInput == "clean" {
			cleanRoom(how)
		} else if userInput == "build" {
			buildRoom(how)
		} else if userInput == "pay" {
			payEmployee(how)
		} else if userInput == "hire" {
			hireEmployee(how)
		} else if userInput == "help" {
			printHelp()
		} else {
			fmt.Println("Unknown command")
		}

		totalHappiness := 0
		for _, employee := range hotel.employees {
			totalHappiness += employee
		}
		if hotel.money == 0 && totalHappiness == 0 && len(hotel.emptyRooms) == 0 {
			log.Fatal("All out of ways to make money")
		}
	}
}

func initFileScanner() {
	fp, err := os.Open("hotel/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fileScanner = bufio.NewScanner(fp)
}

func initStringScanner() {
	stringsReader := strings.NewReader(inputString)
	stringScanner = bufio.NewScanner(stringsReader)
}

func getUserInput(how string) string {
	switch how {
	case "user":
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		return strings.TrimSpace(scanner.Text())
	case "string":
		if stringScanner.Scan() {
			s := strings.TrimSpace(stringScanner.Text())
			Printer.Println(s)
			return s
		}
		log.Fatal("reached end of string but no quit")
		return ""
	case "file":
		if fileScanner.Scan() {
			s := strings.TrimSpace(fileScanner.Text())
			Printer.Println(s)
			return s
		}
		log.Fatal("reached end of file but no quit")
		return ""
	default:
		return ""
	}
}

func checkIn(how string) {
	Printer.Printf("Enter the name of the person checking in: ")
	name := getUserInput(how)
	room, found := hotel.getEmptyRoom()
	if !found {
		Printer.Println("No room in the hotel")
		return
	}
	hotel.checkIn(name, room)
	Printer.Printf("%s checked in to room %d\n", name, room)
}

func checkOut(how string) {
	Printer.Printf("Enter the name of the person checking out: ")
	name := getUserInput(how)
	room, found := hotel.roomLookup(name)
	if !found {
		Printer.Printf("Nobody staying at the hotel named %s\n", name)
		return
	}
	hotel.checkOut(room)
	Printer.Printf("%s checked out of room %d\n", name, room)
}

func cleanRoom(how string) {
	Printer.Printf("Enter the room number you'd like to clean: ")
	roomStr := getUserInput(how)
	room64, err := strconv.ParseInt(roomStr, 10, 32)
	room := int(room64)
	if err != nil {
		Printer.Println("not a number")
		return
	}
	for _, dirtyRoom := range hotel.dirtyRooms {
		if dirtyRoom == room {
			goto Found
		}
	}
	Printer.Printf("Room %d isn't dirty\n", room)
	return
Found:
	for i, employee := range hotel.employees {
		Printer.Printf("Employee %d: Happiness level = %d\n", i, employee)
	}
	Printer.Printf("Enter the number of the employee you wish to clean the room (blank for the happiest): ")
	employeeStr := getUserInput(how)
	var employee int
	var employee64 int64
	if employeeStr == "" {
		goto Happiest
	}
	employee64, err = strconv.ParseInt(employeeStr, 10, 32)
	if err != nil {
		Printer.Println("not a number")
		return
	}
	employee = int(employee64)
	if employee >= len(hotel.employees) || employee < 0 {
		Printer.Println("Not a real employee")
		return
	}
	if hotel.employees[employee] == 0 {
		Printer.Println("That employee is not happy enough")
		return
	}
	hotel.cleanRoom(room, employee)

Happiest:
	var found bool
	employee, found = hotel.getHappiestEmployeeIndex()
	if found == false {
		Printer.Println("No employee has the happiness to clean this")
		return
	}
	hotel.cleanRoom(room, employee)
	Printer.Printf("Room number %d cleaned!\n", room)
}

func buildRoom(how string) {
	Printer.Printf("Enter the room number you'd like to build: ")
	roomStr := getUserInput(how)
	room64, err := strconv.ParseInt(roomStr, 10, 32)
	if err != nil {
		Printer.Println("not a number")
		return
	}
	room := int(room64)
	err = hotel.buyRoom(room)
	if err != nil {
		Printer.Println(err)
		return
	}
	Printer.Printf("Room number %d built!\n", room)
}

func payEmployee(how string) {
	for i, employee := range hotel.employees {
		Printer.Printf("Employee %d: Happiness level = %d\n", i, employee)
	}
	Printer.Printf("Which employee do you want to pay? ")
	employeeIndex, err := getIntFromUser(how)
	if err != nil {
		Printer.Println(err)
		return
	}
	if employeeIndex >= len(hotel.employees) || employeeIndex < 0 {
		Printer.Println("no such employee")
		return
	}
	Printer.Println("Money: ", hotel.money)
	Printer.Printf("How much do you want to pay? ")
	payment, err := getIntFromUser(how)
	if err != nil {
		Printer.Println(err)
		return
	}
	if payment > hotel.money {
		Printer.Println("you don't have that much money")
		return
	} else if payment < 0 {
		Printer.Println("you must pay your employee a positive amount of money")
		return
	}
	hotel.payEmployee(payment, employeeIndex)
	Printer.Println("They've been paid")
}

func hireEmployee(how string) {
	Printer.Println("Money: ", hotel.money)
	Printer.Printf("How much do you want to pay the new employee? ")
	payment, err := getIntFromUser(how)
	if err != nil {
		Printer.Println(err)
		return
	}
	if payment > hotel.money {
		Printer.Println("you don't have that much money")
		return
	}
	if payment < 0 {
		Printer.Println("you must pay your employee a positive amount of money")
		return
	}
	hotel.hireEmployee(payment)
	Printer.Println("New employee hired!")
}

func view() {
	printLine()
	Printer.Println("Hotel Summary")
	printLine()
	Printer.Println()
	Printer.Printf("Total Money: $%d\n", hotel.money)
	Printer.Printf("Total Rooms: %d\n", len(hotel.emptyRooms)+len(hotel.occupiedRooms)+len(hotel.dirtyRooms))
	Printer.Printf("Hotel Occupancy: %.1f%%\n", float64(len(hotel.occupiedRooms))/float64(len(hotel.emptyRooms)+len(hotel.dirtyRooms)+len(hotel.occupiedRooms))*100)
	Printer.Printf("Hotel Vacancy: %.1f%%\n", float64(len(hotel.emptyRooms))/float64(len(hotel.emptyRooms)+len(hotel.dirtyRooms)+len(hotel.occupiedRooms))*100)
	Printer.Println()
	printLine()
	Printer.Printf("Empty Rooms: ")
	if len(hotel.emptyRooms) == 0 {
		Printer.Println(" No empty rooms")
	} else {
		for _, emptyRoom := range hotel.emptyRooms {
			Printer.Printf(" %d", emptyRoom)
		}
	}
	Printer.Println()
	printLine()
	Printer.Println()
	Printer.Println("Occupied Rooms")
	printLine()
	keys := []int{}
	for k := range hotel.occupiedRooms {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, occupiedRoom := range keys {
		person := hotel.occupiedRooms[occupiedRoom]
		Printer.Printf("Room %d: Occupied by %s\n", occupiedRoom, person)
	}
	Printer.Println()
	printLine()
	Printer.Printf("Dirty Rooms:")
	if len(hotel.dirtyRooms) == 0 {
		Printer.Println(" No dirty rooms")
	} else {
		for _, dirtyRooms := range hotel.dirtyRooms {
			Printer.Printf(" %d", dirtyRooms)
		}
	}
	Printer.Println()
	printLine()
	Printer.Println()
	Printer.Println("Employee Summary")
	printLine()
	for i, employee := range hotel.employees {
		Printer.Printf("Employee %d: Happiness level = %d\n", i, employee)
	}
	printLine()
	Printer.Println()
}

func printLine() {
	Printer.Println("--------------------------------------------------------------")
}

func printHelp() {
	Printer.Println("Commands:")
	Printer.Println()
	Printer.Println("view - view summary of the hotel")
	Printer.Println("checkin - check somebody in to the hotel")
	Printer.Println("checkout - check somebody out of the hotel")
	Printer.Println("clean - tell somebody to clean a room")
	Printer.Println("build - expand the hotel")
	Printer.Println("pay - pay a current employee")
	Printer.Println("hire - hire a new employee")
	Printer.Println("help - see all commands")
	Printer.Println()
}

func getIntFromUser(how string) (int, error) {
	str := getUserInput(how)
	num64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, errors.New("not a number")
	}
	return int(num64), nil
}
