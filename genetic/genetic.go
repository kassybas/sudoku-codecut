package genetic

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Child struct {
	genes   []byte
	fitness int
}

func convChildToByte(children []Child) [][]byte {
	generation := make([][]byte, GEN)
	for i := 0; i < GEN; i++ {
		generation[i] = children[i].genes
	}
	return generation
}
func printTopOfGen(children []Child, numOfGen int) {
	fmt.Printf("**** GENERATION: %d ****\n", numOfGen)
	for i := 0; i < 10; i++ {
		fmt.Println(children[i])
	}
}

func getHorizontalErrors(line []byte) int {
	var errors int = 0
	for i, check := range line {
		for j := i + 1; j < len(line); j++ {
			if check == line[j] {
				errors++
			}
		}
	}
	return errors
}

func getFitness(candidate []byte) int {
	var errors int
	errors += getHorizontalErrors(candidate)
	return errors
}

func generateParent(original []byte) []byte {
	parent := make([]byte, len(original))
	for i, val := range original {
		if val == 0 {
			parent[i] = byte(rand.Intn(len(original)) + 1)
		} else {
			parent[i] = val
		}
	}
	return parent
}

func mutateParent(original []byte, parentA []byte, parentB []byte) Child {
	child := make([]byte, len(original))
	for i, val := range original {
		if val == 0 {
			parentchoice := rand.Intn(11)
			if parentchoice == 0 {
				child[i] = byte(rand.Intn(len(original)) + 1)
			} else if parentchoice <= 5 {
				child[i] = parentA[i]
			} else {
				child[i] = parentB[i]
			}
		} else {
			child[i] = val
		}
	}
	return Child{genes: child, fitness: GEN}
}

func createFirstGeneration(original []byte) [][]byte {
	generation := make([][]byte, GEN)
	for i := 0; i < GEN; i++ {
		generation[i] = generateParent(original)
	}
	return generation
}

func getAllFitness(gen [][]byte) []Child {
	var children []Child
	for i := 0; i < GEN; i++ {
		errorCount := getFitness(gen[i])
		children = append(children, Child{genes: gen[i], fitness: errorCount})
	}
	return children
}

type ByFit []Child

func (a ByFit) Len() int           { return len(a) }
func (a ByFit) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFit) Less(i, j int) bool { return a[i].fitness < a[j].fitness }

func sortByFitness(gen [][]byte) []Child {
	children := getAllFitness(gen)
	sort.Sort(ByFit(children))

	return children
}

func createNextGeneration(original []byte, children []Child) []Child {
	for i := GEN / 2; i < GEN; i++ {
		children[i] = mutateParent(original, children[i-(GEN/2)].genes, children[i-(GEN/2-1)].genes)
	}
	return children
}

var GEN int = 100

func Solve(original []byte) []byte {
	//	bestFitness := 0
	rand.Seed(time.Now().UnixNano())
	generation := createFirstGeneration(original)
	children := sortByFitness(generation)
	printTopOfGen(children, 0)

	var numOfGen int
	for children[0].fitness != 0 && numOfGen < 100 {
		numOfGen++
		children = createNextGeneration(original, children)
		children = sortByFitness(convChildToByte(children))
		printTopOfGen(children, numOfGen)
	}

	return original
}
