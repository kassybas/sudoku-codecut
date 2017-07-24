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

func printTopOfGen(children []Child, numOfGen int, genID rune) {
	fmt.Printf("**** %c :: GENERATION: %d ****\n", genID, numOfGen)
	for i := 0; i < 3; i++ {
		printChild(children[i])
	}
}

func createValues(a [][]byte) []byte {
	DIM := len(a)
	values := make([]byte, DIM*DIM)
	for i := 0; i < DIM*DIM; i++ {
		values[i] = byte(i%DIM) + 1
	}
	for i := 0; i < DIM; i++ {
		for j := 0; j < DIM; j++ {
			if a[i][j] == 0 {
				continue
			}
			for k := 0; k < DIM*DIM; k++ {
				if a[i][j] == values[k] {
					values = append(values[:k], values[k+1:]...)
					break
				}
			}
		}
	}
	return values
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

// func getBoxErrors(field [][]byte) int {
// 	size := math.sqrt(len(field))
// 	for i := 0; i < len(field); i += size {
// 		for j := 0; j < len(field); j += size {
// 			checkBox(field, i,j)
// 		}
// 	}
// }

func getFitness(candidate [][]byte) int {
	var errors int
	for i := 0; i < len(candidate); i++ {
		errors += getHorizontalErrors(candidate[i])
		errors += getVerticalErrors(candidate, i)

	}
	return errors
}

func generateParent(original [][]byte, values []byte) [][]byte {
	parent := make([][]byte, len(original))

	for i := 0; i < len(original); i++ {
		parent[i] = make([]byte, len(original))
		for j := 0; j < len(original); j++ {
			if original[i][j] == 0 {
				valChoice := rand.Intn(len(values))
				parent[i][j] = values[valChoice]
				values = append(values[:valChoice], values[valChoice+1:]...)
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
				parentchoice := rand.Intn(2)
				if parentchoice == 0 {
					child[i][j] = parentA[i][j]
				} else {
					child[i][j] = parentB[i][j]
				}
			} else {
				child[i][j] = original[i][j]
			}
		}
	}
	return Child{genes: child, fitness: POP}
}

func createFirstGeneration(original [][]byte) []Child {
	generation := make([]Child, POP)
	values := createValues(original)

	for i := 0; i < POP; i++ {
		generation[i].genes = generateParent(original, values)
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
	for i := 0; i < POP; i++ {
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
	values := createValues(original)

	for i := POP / 4; i < POP; i++ {
		pa := i % (POP / 2)
		pb := rand.Intn(POP)
		children[i] = mutateParent(original, children[pa].genes, children[pb].genes)
	}
	for i := (POP / 3) * 2; i < POP; i++ {
		children[i].genes = generateParent(original, values)
		children[i].fitness = POP
	}
	return children
}

var POP int = 10000
var GEN int = 5000

func generateSubPopulation(original [][]byte, populationSize int, subPupID rune, genReq int, popID rune) chan []Child{
	chChild := make(chan []Child, populationSize)
	children:=sortByFitness(createFirstGeneration(original))
	numOfSubGen:=0
	for children[0].fitness != 0 && numOfSubGen < genReq {
		numOfSubGen++;
		children = sortByFitness(children)
		children = createNextGeneration(original, children)
		printTopOfGen(children, numOfSubGen, popID)
	}
	for _,val := range children{
		chChild <- val
	}
	children = sortByFitness(children)
	return chChild
}

func Solve(original [][]byte) [][]byte {
	rand.Seed(time.Now().UnixNano())
	generation := createFirstGeneration(original)

	children := sortByFitness(generation)
	printTopOfGen(children, 0, '0')

	var numOfGen int
	subPopID:='A'
	for children[0].fitness != 0 && numOfGen < 10 {
		numOfGen++
		fmt.Printf("\nSub-generation: %c\n", subPopID)
		childChan:= generateSubPopulation(original, GEN, subPopID,100, subPopID)[:GEN/10]
		for val := range childChan {
			children = append(children, val)
		}
		children= sortByFitness(children)
		subPopID++
	}
	numOfGen=0
	for children[0].fitness!=0 && numOfGen<GEN/2{
		numOfGen++;
		children = sortByFitness(children)
		children = createNextGeneration(original, children)
		printTopOfGen(children, numOfGen, 'X')
	}

	return original
}
