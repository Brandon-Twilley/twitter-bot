package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var TERMINATING_CHARACTER = "\u65e6"
var INITIATING_CHARACTER = "\u65e5"

type arrayOfNodesBefore struct {
	before []*node
}

type arrayOfNodesAfter struct {
	after []*node
}

type node struct {
	nodesBefore *arrayOfNodesBefore
	word        string
	nodesAfter  *arrayOfNodesAfter
}

type directedGraph struct {
	nodes []*node
}

func createNode() *node {
	currentNode := &node{}
	currentNode.nodesBefore = &arrayOfNodesBefore{}
	currentNode.nodesAfter = &arrayOfNodesAfter{}
	currentNode.word = ""

	return currentNode
}

/*
	This function creates our directedGraph.
	Our directedGraph is to be composed of 'nodes'.
	These nodes are linked by words directly before and after them, so they
	would then be a doubly linked directed graph.
*/
func (currentGraph *directedGraph) buildGraph(array_of_words wordArray) *directedGraph {
	if DEBUG {
		fmt.Print("Building Graph:")
	}

	for i := 0; i < len(array_of_words.array); i++ {
		currentNode := createNode()
		currentGraph.nodes = append(currentGraph.nodes, currentNode)
		currentGraph.nodes[i].word = array_of_words.array[i]

		if i == 0 {
			currentGraph.nodes[i].nodesBefore.before = nil

		} else if i == (len(array_of_words.array) - 1) {
			currentGraph.nodes[i].nodesBefore.before = append(currentGraph.nodes[i].nodesBefore.before, currentGraph.nodes[i-1])
			currentGraph.nodes[i].nodesAfter.after = nil
		} else {
			currentGraph.nodes[i].nodesBefore.before = append(currentGraph.nodes[i].nodesBefore.before, currentGraph.nodes[i-1])
			currentGraph.nodes[i-1].nodesAfter.after = append(currentGraph.nodes[i-1].nodesAfter.after, currentGraph.nodes[i])
		}
	}

	return currentGraph
	//if there already were elements of that string
}

/*
	This function searches through the list of existing elements within the
	entire directed graph structure to find if a matching word exists.  If so,
	the node (structure holding the word) is returned and our function also
	returns a true.  If no word was found, we return a false in our boolean
	return and a null on our node pointer.

	TODO: create hash table and/or concurrent processes for fast traversal
	through directed graph.
	TODO: Combine find_element function for []*node datatype
		DONE
*/

func find_element(in string, nodeArray []*node) (bool, *node, int) {
	currentNode := &node{}
	nodeExists := false
	for i := 0; i < len(nodeArray); i++ {
		//checks through all of our tree elements to see if there exists
		//the element we're looking for.
		//      TODO: add concurrency to this compare function for quicker times.

		if Compare(in, nodeArray[i].word) == 0 {
			nodeExists = true
			currentNode = nodeArray[i]
			return nodeExists, currentNode, i
		}
	}
	currentNode = nil
	return nodeExists, currentNode, 0
}

func (mainGraph *directedGraph) combineGraphs(secondaryGraph directedGraph) *directedGraph {
	if DEBUG {
		fmt.Print("COMBINING TREES: < ")
		for i := 0; i < len(secondaryGraph.nodes); i++ {
			fmt.Print(secondaryGraph.nodes[i].word, " ")
		}
		fmt.Println(" > ")
	}

	sbArray := []*node{}
	saArray := []*node{}
	for i := 0; i < len(secondaryGraph.nodes); i++ {
		exists, ptr, _ := find_element(secondaryGraph.nodes[i].word, mainGraph.nodes)
		if exists {
			//array of existing nodes
			sbArray = append(sbArray, ptr)
		} else {
			//array of nonexistent nodes

			if DEBUG {
				fmt.Println("Creating element << ", secondaryGraph.nodes[i].word)
			}
			ptr := &node{word: secondaryGraph.nodes[i].word}
			ptr.nodesAfter = &arrayOfNodesAfter{}
			ptr.nodesBefore = &arrayOfNodesBefore{}
			saArray = append(saArray, ptr)
			mainGraph.nodes = append(mainGraph.nodes, ptr)
		}
	}

	for i := 0; i < len(sbArray); i++ {

		exists, ptr, _ := find_element(secondaryGraph.nodes[i].word, mainGraph.nodes)

		if !exists {
			fmt.Println("Not all pointers were correctly added to the mainGraph array.  \nTerminating on number: %i", i)
			return nil
		}
		// We check to see if our previous element for the in branch  is already in our mainGraph array before branch.
		if i != 0 {
			exists, ptr2, _ := find_element(secondaryGraph.nodes[i-1].word, ptr.nodesBefore.before)
			//if there is no element matching the string in our before array, we create another element for this.
			/*
				-------------------------------------------------------------------------------------
							BEFORE BRANCH CREATOR
				-------------------------------------------------------------------------------------
			*/

			if !exists {
				_, ptr2, _ = find_element(secondaryGraph.nodes[i-1].word, mainGraph.nodes)
				ptr.nodesBefore.before = append(ptr.nodesBefore.before, ptr2)
				if DEBUG {
					fmt.Println("Appending \"", ptr2.word, "\" -> \"", secondaryGraph.nodes[i].word)
				}

			}
		}

		/*
			-------------------------------------------------------------------------------------
						AFTER BRANCH CREATOR
			-------------------------------------------------------------------------------------
		*/

		if i != len(secondaryGraph.nodes)-1 {
			exists, ptr2, _ := find_element(secondaryGraph.nodes[i+1].word, ptr.nodesAfter.after)

			if !exists {
				_, ptr2, _ = find_element(secondaryGraph.nodes[i+1].word, mainGraph.nodes)
				ptr.nodesAfter.after = append(ptr.nodesAfter.after, ptr2)
			}
		}
	}
	return mainGraph
}

/*
	TODO: Printed results could look better.  Consider integrating vis.js
	functionality to result.
*/
func (mainGraph directedGraph) printGraph() {

	fmt.Println("PRINTING NETWORK")

	for i := 0; i < len(mainGraph.nodes); i++ {
		fmt.Println()

		if len(mainGraph.nodes[i].nodesBefore.before) == 0 {
			fmt.Print("nil")
		} else {
			for j := 0; j < len(mainGraph.nodes[i].nodesBefore.before); j++ {
				fmt.Print("<", mainGraph.nodes[i].nodesBefore.before[j], ">")
			}
		}

		fmt.Print(" <- ", mainGraph.nodes[i].word, " -> ")

		if len(mainGraph.nodes[i].nodesAfter.after) == 0 {
			fmt.Print("nil")
		} else {
			for j := 0; j < len(mainGraph.nodes[i].nodesAfter.after); j++ {
				fmt.Print("<", mainGraph.nodes[i].nodesAfter.after[j], ">")
			}
		}
	}
	fmt.Println("PRINTED NETWORK")
}

/*
	from our graph, we create a subgraph (a specific traversal of the graph that doesn't
	have cycles within it).  This is the phrase we return and post to twitter.

	mainGraph is our main graph.  subgraph is the graph we're creating
*/
func traverseGraph(mainGraph *directedGraph) *directedGraph {
	//create graph and begin our new graph with the initial element.
	_, begin, _ := find_element(INITIATING_CHARACTER, mainGraph.nodes)
	subgraph := &directedGraph{}
	currentMainGraphNode := begin
	var currentSubGraphNode *node
	var previousSubGraphNode *node

	restart := false

	for true {
		exists, _, _ := find_element(currentMainGraphNode.word, subgraph.nodes)
		if exists {
			restart = true
			break
		} else {
			//create our currentsubgraphnode
			currentSubGraphNode = createNode()
			currentSubGraphNode.word = currentMainGraphNode.word
			subgraph = subgraph.appendNode(currentSubGraphNode)

			if Compare(currentSubGraphNode.word, begin.word) == 0 {

			} else { //link previous node with current node if our current node isn't the initial node.
				fmt.Println("PREVIOUS NODE: ", previousSubGraphNode)
				fmt.Println("CURRENT NODE: ", currentSubGraphNode)
				previousSubGraphNode.nodesAfter.after = append(previousSubGraphNode.nodesAfter.after, currentSubGraphNode)
				currentSubGraphNode.nodesBefore.before = append(currentSubGraphNode.nodesBefore.before, previousSubGraphNode)
			}

			//update our previousSubgraphNode
			previousSubGraphNode = currentSubGraphNode
			//traverse our main graph.
			rand.Seed(time.Now().UTC().UnixNano())
			if len(currentMainGraphNode.nodesAfter.after) == 0 {
				//this case only exists for the terminating character.  We break the
				//loop and return the string we constructed.
				break
			}
			i := rand.Intn(len(currentMainGraphNode.nodesAfter.after))
			currentMainGraphNode = currentMainGraphNode.nodesAfter.after[i]
			if Compare(currentMainGraphNode.word, "\x98") == 0 {
				break
			}

		}
	}
	if restart {
		subgraph = traverseGraph(mainGraph)
	}
	return subgraph
}

func (subgraph *directedGraph) iterateGraph() string {
	i := 0
	var buffer bytes.Buffer
	var currentNode *node
	currentNode = subgraph.nodes[i]
	for true {

		buffer.WriteString(currentNode.word)
		buffer.WriteString(" ")

		if Compare(currentNode.word, TERMINATING_CHARACTER) == 0 {
			break
		}
		if len(currentNode.nodesAfter.after) == 0 {
			break
		} else {
			currentNode = currentNode.nodesAfter.after[0]
		}
	}
	return buffer.String()

}

/*
	adds node to directed graph
*/
func (mainGraph directedGraph) appendNode(currentNode *node) *directedGraph {
	mainGraph.nodes = append(mainGraph.nodes, currentNode)
	return &mainGraph
}

/*
	removes node from graph and replaces array slot with the last element
*/
func (mainGraph directedGraph) removeNode(currentNode *node) *directedGraph {
	is_found, _, number := find_element(currentNode.word, mainGraph.nodes)
	if is_found {
		mainGraph.nodes[number] = mainGraph.nodes[len(mainGraph.nodes)-1]
		mainGraph.nodes = mainGraph.nodes[:len(mainGraph.nodes)-1]
	}

	return &mainGraph
}

/*
	From our input, we create an array(slice) containing all sentences.
	From here, we create a beginning character (/u65e5) and a terminating
	character (/u65e6) from
*/
func toWordArray(input string) *wordArray {

	array_of_words := &wordArray{}
	array_of_words.array = append(array_of_words.array, INITIATING_CHARACTER)
	inputArray := strings.Split(input, " ")
	for i := 0; i < len(inputArray); i++ {
		array_of_words.array = append(array_of_words.array, inputArray[i])
	}
	array_of_words.array = append(array_of_words.array, TERMINATING_CHARACTER)

	return array_of_words
}

//simple string comparison operator.
func Compare(a string, b string) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}
