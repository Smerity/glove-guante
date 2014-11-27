package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func VAddSub(A, B []float64, polarity float64) []float64 {
	if len(A) != len(B) {
		panic("Requires equal length")
	}
	vec := make([]float64, len(A))
	for i := range A {
		vec[i] = A[i] + polarity*B[i]
	}
	return vec
}

func VAdd(A, B []float64) []float64 {
	return VAddSub(A, B, +1)
}

func VSub(A, B []float64) []float64 {
	return VAddSub(A, B, -1)
}

func VDot(A, B []float64) float64 {
	dot := 0.0
	for i := range A {
		dot += A[i] * B[i]
	}
	return dot
}

func VNorm(A []float64) float64 {
	dot := VDot(A, A)
	return math.Sqrt(dot)
}

type StringFloatTuple struct {
	Key   string
	Value float64
}

type ByValue []*StringFloatTuple

func (r ByValue) Len() int           { return len(r) }
func (r ByValue) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByValue) Less(i, j int) bool { return r[i].Value > r[j].Value }

func main() {
	words := make(map[string][]float64)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Loading...")
	for scanner.Scan() {
		splits := strings.SplitN(scanner.Text(), " ", 2)
		word := splits[0]
		nums := strings.Split(splits[1], " ")
		vec := make([]float64, 50)
		for i, n := range nums {
			vec[i], _ = strconv.ParseFloat(n, 64)
		}
		words[word] = vec
	}
	fmt.Println("Loaded", len(words))

	for t := 0; t < 100; t++ {
		var results ByValue
		test := "king-man+woman"
		fmt.Println("Find similar words for", test)
		tvec := VAdd(VSub(words["king"], words["man"]), words["woman"])
		for word, wvec := range words {
			sim := VDot(wvec, tvec) / VNorm(tvec) / VNorm(wvec)
			results = append(results, &StringFloatTuple{word, sim})
		}

		sort.Sort(results)
		for i := 0; i < 10; i++ {
			fmt.Println(results[i].Key, results[i].Value)
		}
	}
}
