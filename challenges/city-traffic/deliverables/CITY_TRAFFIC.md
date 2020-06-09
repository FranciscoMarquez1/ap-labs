# City Traffic

## Build Requirements
- Go version 1.14 or above

## Package dependency installation
```sh
$ go get "github.com/golang-collections/collections/stack"
```

## How to compile

```sh 
$ go build city-traffic.go
```

## How to run 

```sh 
$ ./city-traffic -w <size of the map> -c <number of cars> -s <number of semaphores>
```
### Notes

- width must be 9n + 2. For example: 11, 20, 29, 38, ...

- Cars must be less than or equal to the width

- Semaphores must be less than or equal to the number of intersections

- Video presentation [link](https://drive.google.com/file/d/1g5AeBB4NP5WKgyFLmLVoFhfQX7o_NF-e/view?usp=sharing)
