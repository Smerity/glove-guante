package main

import (
	"bufio"
	"compress/gzip"
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

// Add sorting by the value or a StringFloatTuple
type ByValue []*StringFloatTuple

func (r ByValue) Len() int           { return len(r) }
func (r ByValue) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByValue) Less(i, j int) bool { return r[i].Value > r[j].Value }

var words map[string][]float64

func main() {
	words := make(map[string][]float64)

	// Load the word vectors from a gzipped input file, passed in by the first command line argument
	fn := os.Args[1]
	f, _ := os.Open(fn)
	gz, _ := gzip.NewReader(f)
	defer gz.Close()
	defer f.Close()
	fscanner := bufio.NewScanner(gz)
	fmt.Println("Loading", fn, "...")
	//
	for fscanner.Scan() {
		splits := strings.SplitN(fscanner.Text(), " ", 2)
		word := splits[0]
		nums := strings.Split(splits[1], " ")
		vec := make([]float64, len(nums))
		for i, n := range nums {
			vec[i], _ = strconv.ParseFloat(n, 64)
		}
		words[word] = vec
	}
	fmt.Println("Loaded", len(words))

	// Allow the user to query in an interpreter style interface
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf(">>> ")
	for scanner.Scan() {
		testWords := strings.Split(scanner.Text(), " ")
		fmt.Println("=-=-=-=-=")
		// Ensure all the words are in the vocabulary
		notInVocab := make([]string, 0)
		for _, word := range testWords {
			if _, ok := words[word]; !ok {
				notInVocab = append(notInVocab, word)
			}
		}
		if len(notInVocab) > 0 {
			fmt.Println("Words not found in the vocabulary:", notInVocab)
			fmt.Printf(">>> ")
			continue
		}
		// Generate the test vector to compare against
		var tvec []float64
		if len(testWords) == 1 {
			tvec = words[testWords[0]]
		} else if len(testWords) == 3 {
			fmt.Println("Finding similar words for", testWords[0], "-", testWords[1], "+", testWords[2])
			tvec = VAdd(VSub(words[testWords[0]], words[testWords[1]]), words[testWords[2]])
		} else {
			fmt.Println("Testing only works for one or three words")
			fmt.Printf(">>> ")
			continue
		}
		//
		var results ByValue
		for word, wvec := range words {
			//sim := VCosine(tvec, wvec)
			sim := VDot(wvec, tvec) / VNorm(tvec) / VNorm(wvec)
			results = append(results, &StringFloatTuple{word, sim})
		}
		// Sort the results and return the top 10 to the user
		sort.Sort(results)
		for i := 0; i < 10; i++ {
			fmt.Println(results[i].Key, results[i].Value)
		}
		//
		fmt.Printf(">>> ")
	}
}
