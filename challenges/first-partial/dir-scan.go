package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func printOutput(pathname string, dir, symlinks, devices, sockets, other int) {
	fmt.Println("+-------------------+-----+")
	fmt.Printf("|Path               | %s |\n", pathname)
	fmt.Println("+-------------------+-----+")
	fmt.Printf("|Directories:       |%d |\n", dir)
	fmt.Printf("|Symbolic Links:    |%d |\n", symlinks)
	fmt.Printf("|Devices:           |%d |\n", devices)
	fmt.Printf("|Sockets:           |%d |\n", sockets)
	fmt.Printf("|Other Files:       |%d |\n", other)
	fmt.Println("+-------------------+-----+")
}
   
// scanDir stands for the directory scanning implementation
func scanDir(dir string) error {

	directories, symlinks, devices, sockets, other := 0,0,0,0,0

    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error [%v] en path [%q]\n", err, path)
			os.Exit(2)
		}
		Mode := info.Mode()
		switch {
		case Mode.IsDir():
			directories++
		case Mode&os.ModeSymlink != 0:
			symlinks++
		case Mode&os.ModeDevice != 0:
			devices++
		case Mode&os.ModeSocket != 0:
			sockets++
		default:
			other++
		}
        return nil
	})
    if err != nil {
		fmt.Printf("Error [%v]\n", err)
		os.Exit(3)
	}
	printOutput(dir, directories, symlinks, devices, sockets, other)
	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./dir-scan <directory>")
		os.Exit(1)
	}
	scanDir(os.Args[1])
}
