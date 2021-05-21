package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Params struct {
	XResolution int
	YResolution int
	NumRows     int
	NumPoints   uint64
}

var params = Params{}

var headerLine = regexp.MustCompile(`^#\s*(\S+)+ (\S+)$`)

func convertLine(line string, outputFormat string) string {
	fields := strings.Split(line, ",")
	if len(fields) != 4 {
		log.Fatalf("invalid line: %s", line)
	}

	c, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Fatal(err)
	}
	c = c % params.XResolution
	a := math.Pi * 2.0 * (float64(c) / float64(params.XResolution))

	r, err := strconv.Atoi(fields[2])
	if err != nil {
		log.Fatal(err)
	}
	b := math.Pi * 2.0 * (float64(r) / float64(params.YResolution))

	d, err := strconv.ParseFloat(fields[3], 64)
	d = math.Max(d-1.7, 0) // Compensating for LIDAR distance from shaft

	if err != nil {
		log.Fatal(err)
	}

	x := d * math.Cos(b) * math.Sin(a)
	y := d * math.Cos(b) * math.Cos(a)
	z := d * math.Sin(b)

	switch outputFormat {
	case "pcd":
		return fmt.Sprintf("%f %f %f", x/100, y/100, z/100)
	case "asc":
		return fmt.Sprintf("%f %f %f 0 0 0 0", x/100, y/100, z/100)
	default:
		log.Fatalf("unknown format: %s\n", outputFormat)
		return ""
	}
}

func printPCDHeader(p Params) {
	p.NumPoints = (uint64)(p.NumRows) * (uint64)(p.XResolution)
	tmpl, err := template.New("pcd-header").Parse(`# .PCD v.7 - Point Cloud Data file format
VERSION .7
FIELDS x y z
SIZE 4 4 4
TYPE F F F
COUNT 1 1 1
WIDTH {{.XResolution}}
HEIGHT {{.NumRows}}
VIEWPOINT 0 0 0 1 0 0 0
POINTS {{.NumPoints}}
DATA ascii
`)
	if err != nil {
		log.Fatalf("unable to parse PCD header template: %v\n", err)
	}
	err = tmpl.Execute(os.Stdout, p)
	if err != nil {
		log.Fatalf("unable to execute PCD header template: %v\n", err)
	}
}

func processHeader(line string) {
	headerLine.MatchString(line)

	groups := headerLine.FindStringSubmatch(line)
	if len(groups) != 3 {
		log.Fatalf("invalid header: %s\n", line)
	}

	var err error
	if groups[1] == "version" {
		if groups[2] != "1" {
			log.Fatalf("unknown version in header: %s\n", groups[2])
		}
	}
	if groups[1] == "x-resolution" {
		params.XResolution, err = strconv.Atoi(groups[2])
		if err != nil {
			log.Fatalf("invalid header: %s\n", line)
		}
	}
	if groups[1] == "y-resolution" {
		params.YResolution, err = strconv.Atoi(groups[2])
		if err != nil {
			log.Fatalf("invalid header: %s\n", line)
		}
	}
	if groups[1] == "num-rows" {
		params.NumRows, err = strconv.Atoi(groups[2])
		if err != nil {
			log.Fatalf("invalid header: %s\n", line)
		}
	}
}

func convertFile(fileName string, outputFormat string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	firstPoint := true
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "#") {
			processHeader(t)
		} else {
			if firstPoint && outputFormat == "pcd" {
				printPCDHeader(params)
				firstPoint = false
			}
			l := convertLine(scanner.Text(), outputFormat)
			fmt.Println(l)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

var inputFile = flag.String("input", "", "Input .asdp file to convert.")
var outputFormat = flag.String("outputFormat", "pcd", "Output format. Must be either 'asc' or 'pcd'.")

func main() {
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Must provide input file name with --input.")
	}

	if *outputFormat != "asc" && *outputFormat != "pcd" {
		log.Fatal("Output format must be either 'asc' or 'pcd'.")
	}

	convertFile(*inputFile, *outputFormat)
}
