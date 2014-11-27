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

func VCosine(A, B []float64) float64 {
	return VDot(A, B) / VNorm(A) / VNorm(B)
}

type StringFloatTuple struct {
	Key   string
	Value float64
}

type ByValue []*StringFloatTuple

func (r ByValue) Len() int           { return len(r) }
func (r ByValue) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByValue) Less(i, j int) bool { return r[i].Value > r[j].Value }

var words map[string][]float64

func main() {
	words := make(map[string][]float64)

	f, _ := os.Open("glove.6B.50d.txt")
	fscanner := bufio.NewScanner(f)
	fmt.Println("Loading...")
	for fscanner.Scan() {
		splits := strings.SplitN(fscanner.Text(), " ", 2)
		word := splits[0]
		nums := strings.Split(splits[1], " ")
		vec := make([]float64, 50)
		for i, n := range nums {
			vec[i], _ = strconv.ParseFloat(n, 64)
		}
		words[word] = vec
	}
	fmt.Println("Loaded", len(words))

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf(">>> ")
	for scanner.Scan() {

		testWords := strings.Split(scanner.Text(), " ")
		if len(testWords) != 3 {
			fmt.Println("Testing only works for three words")
		}
		fmt.Println("=-=-=-=-=")
		fmt.Println("Finding similar words for", testWords[0], "-", testWords[1], "+", testWords[2])
		tvec := VAdd(VSub(words[testWords[0]], words[testWords[1]]), words[testWords[2]])
		//
		var results ByValue
		for word, wvec := range words {
			//sim := VCosine(tvec, wvec)
			sim := VDot(wvec, tvec) / VNorm(tvec) / VNorm(wvec)
			results = append(results, &StringFloatTuple{word, sim})
		}

		sort.Sort(results)
		for i := 0; i < 10; i++ {
			fmt.Println(results[i].Key, results[i].Value)
		}

		fmt.Printf(">>> ")
	}
}
