package genetic

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Child struct {
	genes   [][]byte
	fitness int
}

func printChild(child Child) {
	for i := 0; i < len(child.genes); i++ {
		fmt.Println(child.genes[i])
	}
	fmt.Printf("ERRORS: %d\n", child.fitness)
}

func printTopOfGen(children []Child, numOfGen int) {
	fmt.Printf("**** GENERATION: %d ****\n", numOfGen)
	for i := 0; i < 3; i++ {
		printChild(children[i])
	}
}

func getHorizontalErrors(line []byte) int {
	var errors int
	for i, check := range line {
		for j := i + 1; j < len(line); j++ {
			if check == line[j] {
				errors++
			}
		}
	}
	return errors
}

func getVerticalErrors(field [][]byte, col int) int {
	var errors int
	for j := 0; j < len(field); j++ {
		check := field[j][col]
		for k := j + 1; k < len(field); k++ {
			if check == field[k][col] {
				errors++
			}
		}
	}
	return errors
}

func getFitness(candidate [][]byte) int {
	var errors int
	for i := 0; i < len(candidate); i++ {
		errors += getHorizontalErrors(candidate[i])
		errors += getVerticalErrors(candidate, i)
	}
	return errors
}

func generateParent(original [][]byte) [][]byte {
	parent := make([][]byte, len(original))

	for i := 0; i < len(original); i++ {
		parent[i] = make([]byte, len(original))
		for j := 0; j < len(original); j++ {
			if original[i][j] == 0 {
				parent[i][j] = byte(rand.Intn(len(original)) + 1)
			} else {
				parent[i][j] = original[i][j]
			}
		}
	}
	return parent
}

func mutateParent(original [][]byte, parentA [][]byte, parentB [][]byte) Child {
	child := make([][]byte, len(original))
	for i := 0; i < len(original); i++ {
		child[i] = make([]byte, len(original))
		for j := 0; j < len(original); j++ {
			if original[i][j] == 0 {
				parentchoice := rand.Intn(12)
				if parentchoice >= 10 {
					child[i][j] = byte(rand.Intn(len(original)) + 1)
				} else if parentchoice < 5 {
					child[i][j] = parentA[i][j]
				} else {
					child[i][j] = parentB[i][j]
				}
			} else {
				child[i][j] = original[i][j]
			}
		}
	}
	return Child{genes: child, fitness: GEN}
}

func createFirstGeneration(original [][]byte) []Child {
	generation := make([]Child, GEN)
	for i := 0; i < GEN; i++ {
		generation[i].genes = generateParent(original)
	}
	return generation
}

func getChild(genes [][]byte, childch chan Child) {
	errorCount := getFitness(genes)
	childch <- Child{genes: genes, fitness: errorCount}
}

func getAllFitness(gen []Child) []Child {
	var children []Child
	childch := make(chan Child)
	for i := 0; i < GEN; i++ {
		go getChild(gen[i].genes, childch)
		child := <-childch
		children = append(children, child)
	}
	return children
}

type ByFit []Child

func (a ByFit) Len() int           { return len(a) }
func (a ByFit) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFit) Less(i, j int) bool { return a[i].fitness < a[j].fitness }

func sortByFitness(gen []Child) []Child {
	children := getAllFitness(gen)
	sort.Sort(ByFit(children))
	return children
}

func createNextGeneration(original [][]byte, children []Child) []Child {
	for i := GEN / 4; i < GEN-(GEN/4); i++ {
		pa := i - GEN/4
		pb := GEN - GEN/4 - i
		children[i] = mutateParent(original, children[pa].genes, children[pb].genes)
	}
	for i := GEN - (GEN / 4); i < GEN; i++ {
		children[i].genes = generateParent(original)
		children[i].fitness = GEN
	}

	return children
}

var GEN int = 10000

func Solve(original [][]byte) [][]byte {
	//	bestFitness := 0
	rand.Seed(time.Now().UnixNano())
	generation := createFirstGeneration(original)

	children := sortByFitness(generation)
	printTopOfGen(children, 0)

	var numOfGen int
	for children[0].fitness != 0 && numOfGen < 10000 {
		numOfGen++
		children = sortByFitness(children)
		children = createNextGeneration(original, children)
		printTopOfGen(children, numOfGen)
	}

	return original
}
