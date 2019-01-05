package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func wordCount(r io.Reader) map[string]int {
	m := make(map[string]int)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	re := regexp.MustCompile(`[a-z]+`)
	var line string

	for scanner.Scan() {
		line = strings.ToLower(scanner.Text())
		match := re.FindAllStringSubmatch(line, -1)
		for _, s := range match {
			m[s[0]]++
		}
	}
	return m
}

// writeMapToFile writes map to the csv file with prior sorting of the words by the number of appearances
func writeMapToFile(m map[string]int, f io.Writer) {
	type kv struct {
		key   string
		value int
	}

	var s []kv
	for k, v := range m {
		s = append(s, kv{k, v})
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].value > s[j].value
	})

	nw := csv.NewWriter(f)
	line := make([]string, 2)
	for _, v := range s {
		line[0] = v.key
		line[1] = strconv.Itoa(m[v.key])
		err := nw.Write(line)
		if err != nil {
			panic(err)
		}
	}
	nw.Flush()
}

func main() {
	t1 := time.Now()

	if len(os.Args) < 3 {
		fmt.Println(`Arguments not corect. Example: wordCount "word" "file.txt"`)
		return
	}
	wordOfInterest := os.Args[1]
	srcFilePath := os.Args[2]
	dstFilePath := srcFilePath[:len(srcFilePath)-4] + ".csv"

	f, err := os.Open(srcFilePath)
	if err != nil {
		fmt.Println("Cannot open file:", srcFilePath)
		panic(err)
	}
	defer f.Close()

	fOut, err := os.Create(dstFilePath)
	if err != nil {
		fmt.Println("Cannot create file:", dstFilePath)
		panic(err)
	}
	defer fOut.Close()

	m := wordCount(f)

	fmt.Printf("Word '%s' appears %d times in file %s \n", wordOfInterest, m[wordOfInterest], srcFilePath)

	writeMapToFile(m, fOut)

	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
