package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var cols = 1600
var rows = 1200

func convertLine(line string) string {
	fields := strings.Split(line, ",")
	if len(fields) != 4 {
		log.Fatalf("invalid line: %s", line)
	}

	c, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Fatal(err)
	}
	c = c % cols
	a := math.Pi * 2.0 * (float64(c) / float64(cols))

	r, err := strconv.Atoi(fields[2])
	if err != nil {
		log.Fatal(err)
	}
	b := math.Pi * 2.0 * (float64(r) / float64(rows))

	d, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		log.Fatal(err)
	}

	x := d * math.Cos(b) * math.Sin(a)
	y := d * math.Cos(b) * math.Cos(a)
	z := d * math.Sin(b)

	return fmt.Sprintf("%f %f %f 0 0 0 0", x, y, z)
}

func convertFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := convertLine(scanner.Text())
		fmt.Println(l)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func main() {
	convertFile("/home/dvoros/Documents/points-10rows.txt")
}
