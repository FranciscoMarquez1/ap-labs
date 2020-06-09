package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-collections/collections/stack"
)

type cell struct {
	id     int
	x      int
	y      int
	dir    int
	hasCar bool
	green  bool
}

type car struct {
	no    int
	x     int
	y     int
	speed int
	idle  int
	path  []cell
}

type semaphore struct {
	cells []cell
	index int
	speed int
}

type qElement struct {
	prev *qElement
	c    cell
}

const STREET = 0
const BUILDING = 1

const NO_DIR = -1
const U = 0
const R = 1
const D = 2
const L = 3
const RU = 4
const DR = 5
const UL = 6
const LD = 7

var width int
var nCars int
var nSemaphores int
var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var board [][]cell
var paths [][]cell
var cars []car

func getIntersections() [][]cell {
	intersections := make([][]cell, 0)

	for i := 0; i < width; i += 9 {
		for j := 0; j < width; j += 9 {
			intersection := make([]cell, 0)

			if i > 0 {
				intersection = append(intersection, board[i-1][j])
			}

			if j > 0 {
				intersection = append(intersection, board[i+1][j-1])
			}

			if i < width-2 {
				intersection = append(intersection, board[i+2][j+1])
			}

			if j < width-2 {
				intersection = append(intersection, board[i][j+2])
			}

			intersections = append(intersections, intersection)
		}
	}

	return intersections
}

func main() {
	wFlag := flag.Int("w", 29, "Positive integer with a value of 9n+2")
	cFlag := flag.Int("c", 16, "Positive integer with a value lower or equal to width")
	sFlag := flag.Int("s", 16, "Positive integer")
	flag.Parse()

	width = *wFlag
	nCars = *cFlag
	nSemaphores = *sFlag

	if width < 0 || (width-2)%9 != 0 {
		log.Fatalf("Width must have a value of 9n+2 (11, 20, 29, 38, 47).")
	}

	var streetCells []cell
	board, streetCells = createBoard()

	if nCars > width || nCars < 0 {
		log.Fatalf("Number of cars must be a positive integer, but it can't exceed the width of the map due to visualization limitations.")
	}

	intersections := getIntersections()

	if nSemaphores < 0 || nSemaphores > len(intersections) {
		log.Fatalf("In a board of width %d, you have %d intersections, meaning that the number of semaphores must be equal or less to that.", width, len(intersections))
	}

	showSplashScreen()

	initSemaphores(intersections)

	ch := make(chan int, nCars)
	initCars(streetCells, &ch)
	displaySim(&ch)

	fmt.Print("Do you want to visualize each car's individual report? (y/n): ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	for strings.ToLower(input.Text()) != "y" && strings.ToLower(input.Text()) != "n" {
		fmt.Print("Do you want to visualize each car's individual report? (y/n): ")
		input = bufio.NewScanner(os.Stdin)
		input.Scan()
	}

	if strings.ToLower(input.Text()) == "y" {
		printReports()
	}
}

func initSemaphores(intersections [][]cell) {
	semaphores := make([]semaphore, 0)

	for i := 0; i < nSemaphores; i++ {
		n := r.Intn(len(intersections))

		intersection := intersections[n]

		intersections[len(intersections)-1], intersections[n] = intersections[n], intersections[len(intersections)-1]
		intersections = intersections[:len(intersections)-1]

		speed := r.Intn(1200-800) + 800
		s := semaphore{intersection, 0, speed}
		setInitState(&s)
		semaphores = append(semaphores, s)
	}

	for i := 0; i < len(semaphores); i++ {
		index := i

		go func() {
			for {
				changeState(&semaphores[index])
				time.Sleep(time.Duration(semaphores[index].speed) * time.Millisecond)
			}
		}()
	}

}

func setInitState(s *semaphore) {
	for i := 0; i < len(s.cells); i++ {
		x := s.cells[i].x
		y := s.cells[i].y
		board[x][y].green = false
	}
}

func changeState(sem *semaphore) {
	length := len(sem.cells)

	cX := sem.cells[sem.index].x
	cY := sem.cells[sem.index].y
	board[cX][cY].green = false

	sem.index = (sem.index + 1) % length

	nX := sem.cells[sem.index].x
	nY := sem.cells[sem.index].y
	board[nX][nY].green = true
}

func initCars(streetCells []cell, ch *chan int) {
	cars = make([]car, 0)
	for i := 0; i < nCars; i++ {

		n := r.Intn(len(streetCells))
		cell1 := streetCells[n]

		streetCells[len(streetCells)-1], streetCells[n] = streetCells[n], streetCells[len(streetCells)-1]
		streetCells = streetCells[:len(streetCells)-1]

		n2 := r.Intn(len(streetCells))
		cell2 := streetCells[n2]

		speed := r.Intn(200) + 50

		path := getPath(cell1, cell2)
		paths = append(paths, path)

		c := car{i, cell1.x, cell1.y, speed, 0, path}

		cars = append(cars, c)
		addCar(c)
	}

	for i := 0; i < len(cars); i++ {
		index := i

		go func() {
			for len(cars[index].path) > 0 {
				time.Sleep(time.Duration(cars[index].speed) * time.Millisecond)
				moveCar(&cars[index])
			}

			*ch <- cars[index].no
			cars[index].speed = 0
			removeCar(&cars[index])
		}()
	}
}

func addCar(c car) {
	if !board[c.x][c.y].hasCar {
		board[c.x][c.y].hasCar = true
	}
}

func removeCar(c *car) {
	if board[c.x][c.y].hasCar {
		board[c.x][c.y].hasCar = false
	}
}

func moveCar(c *car) {
	cX := c.x
	cY := c.y

	nX := c.path[0].x
	nY := c.path[0].y

	if !board[nX][nY].hasCar && board[cX][cY].green {
		board[cX][cY].hasCar = false

		c.x = nX
		c.y = nY
		board[c.x][c.y].hasCar = true

		c.path = c.path[1:]

		if c.speed > 80 {
			c.speed -= 10
		}

		c.idle = 0
	} else {
		if c.idle <= 2 {
			c.idle++

			if c.speed < 450 {
				c.speed += 10
			}
		} else {
			c.speed = 300
		}
	}
}

func getPath(origin cell, destination cell) []cell {
	q := qElement{nil, origin}
	queue := make([]qElement, 0)
	queue = append(queue, q)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr.c == destination {
			return build(&curr)
		}
		neighbors := getNeighbors(curr.c)
		for i := 0; i < len(neighbors); i++ {
			q2 := qElement{&curr, neighbors[i]}
			queue = append(queue, q2)
		}
	}

	return nil
}

func getNeighbors(source cell) []cell {
	var neighbors = make([]cell, 0)
	x := source.x
	y := source.y

	d := source.dir

	switch d {
	case L:
		if y > 0 {
			c := board[x][y-1]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	case R:
		if y < width-1 {
			c := board[x][y+1]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	case U:
		if x > 0 {
			c := board[x-1][y]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	case D:
		if x < width-1 {
			c := board[x+1][y]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	case DR:
		if x < width-1 {
			c := board[x+1][y]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
		if y < width-1 {
			c := board[x][y+1]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	case RU:
		if y < width-1 {
			c := board[x][y+1]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
		if x > 0 {
			c := board[x-1][y]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	case LD:
		if y > 0 {
			c := board[x][y-1]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
		if x < width-1 {
			c := board[x+1][y]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	case UL:
		if x > 0 {
			c := board[x-1][y]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
		if y > 0 {
			c := board[x][y-1]
			if c.id != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	}

	return neighbors
}

func build(tail *qElement) []cell {
	path := make([]cell, 0)
	stk := stack.New()

	for tail != nil {
		stk.Push(tail.c)
		tail = tail.prev
	}

	stk.Pop()

	for stk.Len() > 0 {
		path = append(path, stk.Pop().(cell))
	}

	return path
}

func createBoard() ([][]cell, []cell) {

	var board = make([][]cell, 0)
	var streetCells = make([]cell, 0)

	for i := 0; i < width; i++ {
		var line = make([]cell, 0)
		iMod := i % 9

		for j := 0; j < width; j++ {
			jMod := j % 9

			if iMod < 2 || jMod < 2 {
				dir := NO_DIR

				if iMod == 0 && jMod == 0 {
					dir = LD
				} else if iMod == 1 && jMod == 0 {
					dir = DR
				} else if iMod == 0 && jMod == 1 {
					dir = UL
				} else if iMod == 1 && jMod == 1 {
					dir = RU
				} else if iMod == 0 {
					dir = L
				} else if iMod == 1 {
					dir = R
				} else if jMod == 0 {
					dir = D
				} else if jMod == 1 {
					dir = U
				}

				c := cell{STREET, i, j, dir, false, true}

				if dir == L || dir == R || dir == U || dir == D {
					streetCells = append(streetCells, c)
				}

				line = append(line, c)
			} else {
				c := cell{BUILDING, i, j, NO_DIR, false, true}
				line = append(line, c)
			}
		}

		board = append(board, line)
	}

	return board, streetCells
}

func showSplashScreen() {
	fmt.Println("\033[H\033[2J")
	fmt.Println("	░█████╗░██╗████████╗██╗░░░██╗  ████████╗██████╗░░█████╗░███████╗███████╗██╗░█████╗░          ")
	fmt.Println("	██╔══██╗██║╚══██╔══╝╚██╗░██╔╝  ╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██╔════╝██║██╔══██╗          ")
	fmt.Println("	██║░░╚═╝██║░░░██║░░░░╚████╔╝░  ░░░██║░░░██████╔╝███████║█████╗░░█████╗░░██║██║░░╚═╝          ")
	fmt.Println("	██║░░██╗██║░░░██║░░░░░╚██╔╝░░  ░░░██║░░░██╔══██╗██╔══██║██╔══╝░░██╔══╝░░██║██║░░██╗          ")
	fmt.Println("	╚█████╔╝██║░░░██║░░░░░░██║░░░  ░░░██║░░░██║░░██║██║░░██║██║░░░░░██║░░░░░██║╚█████╔╝          ")
	fmt.Println("	░╚════╝░╚═╝░░░╚═╝░░░░░░╚═╝░░░  ░░░╚═╝░░░╚═╝░░╚═╝╚═╝░░╚═╝╚═╝░░░░░╚═╝░░░░░╚═╝░╚════╝░          ")
	fmt.Println()
	fmt.Println("	░██████╗██╗███╗░░░███╗██╗░░░██╗██╗░░░░░░█████╗░████████╗░█████╗░██████╗░  ░░███╗░░░░░░█████╗░")
	fmt.Println("	██╔════╝██║████╗░████║██║░░░██║██║░░░░░██╔══██╗╚══██╔══╝██╔══██╗██╔══██╗  ░████║░░░░░██╔══██╗")
	fmt.Println("	╚█████╗░██║██╔████╔██║██║░░░██║██║░░░░░███████║░░░██║░░░██║░░██║██████╔╝  ██╔██║░░░░░██║░░██║")
	fmt.Println("	░╚═══██╗██║██║╚██╔╝██║██║░░░██║██║░░░░░██╔══██║░░░██║░░░██║░░██║██╔══██╗  ╚═╝██║░░░░░██║░░██║")
	fmt.Println("	██████╔╝██║██║░╚═╝░██║╚██████╔╝███████╗██║░░██║░░░██║░░░╚█████╔╝██║░░██║  ███████╗██╗╚█████╔╝")
	fmt.Println("	╚═════╝░╚═╝╚═╝░░░░░╚═╝░╚═════╝░╚══════╝╚═╝░░╚═╝░░░╚═╝░░░░╚════╝░╚═╝░░╚═╝  ╚══════╝╚═╝░╚════╝░")
	fmt.Println("                                                                                                ")
	fmt.Println("          _[]_          	                                   _._")
	fmt.Println("         [____]         	                              _.-=\"_-         _")
	fmt.Println("     .----'  '----.     	                         _.-=\"   _-          | ||\"\"\"\"\"\"\"---._______     __..")
	fmt.Println(" .===|    .==.    |===. 	             ___.===\"\"\"\"-.______-,,,,,,,,,,,,`-''----\" \"\"\"\"\"       \"\"\"\"\"  __'")
	fmt.Println(" \\   |   /####\\   |   / 	      __.--\"\"     __        ,'                   o \\           __        [__|")
	fmt.Println(" /   |   \\####/   |   \\      __-\"\"=======.--\"\"  \"\"--.=================================.--\"\"  \"\"--.=======:")
	fmt.Println(" '===|    `\"\"`    |===' 	]       [w] : /        \\ : |========================|    : /        \\ :  [w] :")
	fmt.Println(" .===|    .==.    |===. 	V___________:|          |: |========================|    :|          |:   _-\"")
	fmt.Println(" \\   |   /::::\\   |   / 	 V__________: \\        / :_|=======================/_____: \\        / :__-\"")
	fmt.Println(" /   |   \\::::/   |   \\ 	 -----------'  \"-____-\"  `-------------------------------'  \"-____-\"")
	fmt.Println(" '===|    `\"\"`    |==='")
	fmt.Println(" .===|    .==.    |===.")
	fmt.Println(" \\   |   /&&&&\\   |   /")
	fmt.Println(" /   |   \\&&&&/   |   \\")
	fmt.Println(" '===|    `\"\"`    |==='")
	fmt.Println("     '--.______.--'")
}

func displaySim(ch *chan int) {
	for {
		printBoard()
		fmt.Printf("Cars remaining: %v\n", nCars-len(*ch))
		if len(*ch) >= nCars {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func printBoard() {
	fmt.Println("\033[H\033[2J")

	for i := 0; i < width; i++ {
		line := ""

		for j := 0; j < width; j++ {
			if board[i][j].hasCar {
				switch board[i][j].dir {
				case L:
					line += "—◀—"
				case R:
					line += "—▶—"
				case D:
					line += "▏▼▕"
				case U:
					line += "▏▲▕"
				case LD:
					line += "⌜■┌"
				case DR:
					line += "⌞■└"
				case UL:
					line += "┐■⌝"
				case RU:
					line += "┘■⌟"
				}
			} else {
				switch board[i][j].id {
				case STREET:
					switch board[i][j].dir {
					case L:
						line += "———"
					case R:
						line += "———"
					case D:
						line += "▏╎▕"
					case U:
						line += "▏╎▕"
					case LD:
						line += "⌜ ┌"
					case DR:
						line += "⌞ └"
					case UL:
						line += "┐ ⌝"
					case RU:
						line += "┘ ⌟"
					default:
						line += "   "
					}
					break
				case BUILDING:
					line += "███"
				}
			}
		}
		if i < nCars {
			cIndex := strconv.Itoa(i)
			if i < 10 {
				cIndex = "0" + cIndex
			}

			if cars[i].speed == 0 {
				line += "	Car " + cIndex + ": Arrived at destination."
			} else {
				if cars[i].speed > 240 {
					line += "	Car " + cIndex + "'s speed: 0 mph"
				} else {
					line += "	Car " + cIndex + "'s speed: " + strconv.Itoa(1250/cars[i].speed) + "mph"
				}
			}
		}
		fmt.Println(line)
	}
}

func printReports() {
	for i := 0; i < len(paths); i++ {
		fmt.Printf("Car %d report: \n", i+1)
		fmt.Printf("Origin> [%d, %d] --- Destination> [%d, %d] \n", paths[i][0].x, paths[i][0].y, paths[i][len(paths[i])-1].x, paths[i][len(paths[i])-1].y)
		fmt.Println("Route followed:")
		fmt.Printf("[X: %d,Y: %d]", paths[i][0].x, paths[i][0].y)

		for j := 1; j < len(paths[i]); j++ {
			fmt.Printf(" -> [X: %d,Y: %d]", paths[i][j].x, paths[i][j].y)
		}

		fmt.Println("\n\n")
	}
}
