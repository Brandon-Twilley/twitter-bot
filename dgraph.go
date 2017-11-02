package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type arrayOfNodesBefore struct {
	before []*node
}

type arrayOfNodesAfter struct {
	after []*node
}

type node struct {
	nodesBefore *arrayOfNodesBefore
	str         string
	nodesAfter  *arrayOfNodesAfter
}

type directedGraph struct {
	nodes []*node
}

/*
	This function creates our directedGraph.
	Our directedGraph is to be composed of 'nodes'.
	These nodes are linked by words directly before and after them, so they
	would then be a doubly linked directed graph.
*/
func (currentGraph *directedGraph) buildGraph(strA stringArray) *directedGraph {
	if DEBUG {
		fmt.Print("Building Graph:")
	}

	for i := 0; i < len(strA.array); i++ {
		currentNode := &node{}
		currentNode.nodesBefore = &arrayOfNodesBefore{}
		currentNode.nodesAfter = &arrayOfNodesAfter{}
		currentGraph.nodes = append(currentGraph.nodes, currentNode)
		currentGraph.nodes[i].str = strA.array[i]

		if i == 0 {
			currentGraph.nodes[i].nodesBefore.before = nil

		} else if i == (len(strA.array) - 1) {
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
*/
func (dGraph directedGraph) find_element(in string) (bool, *node) {
	sb := &node{}
	ret := false
	for i := 0; i < len(dGraph.nodes); i++ {
		//checks through all of our tree elements to see if there exists
		//the element we're looking for.
		//      TODO: add concurrency to this compare function for quicker times.

		if Compare(in, dGraph.nodes[i].str) == 0 {
			ret = true
			sb = dGraph.nodes[i]
			return ret, sb
		}
	}
	sb = nil
	return ret, sb
}

/*
	This function searches through the list of existing elements within all
	words that were found to precede the word we're looking at.  If it exists,
	the node (structure holding the word) is returned and our function also
	returns a true.  If no word was found, we return a false in our boolean
	and a null on our node pointer.
*/
func (stbbef arrayOfNodesBefore) find_element(in string) (bool, *node) {

	exists := false
	sbb := &node{}
	for i := 0; i < len(stbbef.before); i++ {
		if stbbef.before[i].str == in {
			exists = true
			sbb = stbbef.before[i]
		}
	}
	return exists, sbb
}

/*
	This function searches through the list of existing elements within all
	words that were found after the word we're looking at.  If it exists,
	the node (structure holding the word) is returned and our function also
	returns a true.  If no word was found, we return a false in our boolean
	and a null on our node pointer.
*/

func (staft arrayOfNodesAfter) find_element(in string) (bool, *node) {
	exists := false
	saft := &node{}
	for i := 0; i < len(staft.after); i++ {
		if staft.after[i].str == in {
			exists = true
			saft = staft.after[i]
		}
	}
	return exists, saft
}

/*
   This creates 2 arrays.  The first array stores all elements that exist both in the dGraph
   structure and the in structure.  Our second array stores all the elements that exist in
   our in structre, but not our dGraph structure.  Afterwards, this adds all elements into
   another array called the updateArray.
*/
func (dGraph *directedGraph) combineGraphs(in directedGraph, weight int) *directedGraph {
	if DEBUG {
		fmt.Print("COMBINING TREES: < ")
		for i := 0; i < len(in.nodes); i++ {
			fmt.Print(in.nodes[i].str, " ")
		}
		fmt.Println(" > ")
	}

	sbArray := []*node{}
	saArray := []*node{}
	for i := 0; i < len(in.nodes); i++ {
		exists, ptr := dGraph.find_element(in.nodes[i].str)
		if exists {
			//array of existing branch nodes
			sbArray = append(sbArray, ptr)
		} else {
			//array of nonexistent branch nodes

			if DEBUG {
				fmt.Println("Creating element << ", in.nodes[i].str)
			}
			ptr := &node{str: in.nodes[i].str}
			ptr.nodesAfter = &arrayOfNodesAfter{}
			ptr.nodesBefore = &arrayOfNodesBefore{}
			saArray = append(saArray, ptr)
			dGraph.nodes = append(dGraph.nodes, ptr)
		}
	}

	for i := 0; i < len(sbArray); i++ {

		exists, ptr := dGraph.find_element(in.nodes[i].str)

		if !exists {
			fmt.Println("Not all pointers were correctly added to the dGraph array.  \nTerminating on number: %i", i)
			return nil
		}
		// We check to see if our previous element for the in branch  is already in our dGraph array before branch.
		if i != 0 {
			exists, ptr2 := ptr.nodesBefore.find_element(in.nodes[i-1].str)

			//if there is no element matching the string in our before array, we create another element for this.
			/*
				-------------------------------------------------------------------------------------
							BEFORE BRANCH CREATOR
				-------------------------------------------------------------------------------------
			*/

			if !exists {
				_, ptr2 = dGraph.find_element(in.nodes[i-1].str)
				ptr.nodesBefore.before = append(ptr.nodesBefore.before, ptr2)
				if DEBUG {
					fmt.Println("Appending \"", ptr2.str, "\" -> \"", in.nodes[i].str)
				}

			}
		}

		/*
			-------------------------------------------------------------------------------------
						AFTER BRANCH CREATOR
			-------------------------------------------------------------------------------------
		*/

		if i != len(in.nodes)-1 {

			exists, ptr2 := ptr.nodesAfter.find_element(in.nodes[i+1].str)

			if !exists {
				_, ptr2 = dGraph.find_element(in.nodes[i+1].str)
				ptr.nodesAfter.after = append(ptr.nodesAfter.after, ptr2)
			}
		}
	}
	return dGraph
}

/*
	TODO: Printed results could look better.  Consider integrating vis.js
	functionality to result.
*/
func (dGraph directedGraph) printGraph() {

	fmt.Println("PRINTING NETWORK")

	for i := 0; i < len(dGraph.nodes); i++ {
		fmt.Println()

		if len(dGraph.nodes[i].nodesBefore.before) == 0 {
			fmt.Print("nil")
		} else {
			for j := 0; j < len(dGraph.nodes[i].nodesBefore.before); j++ {
				fmt.Print("<", dGraph.nodes[i].nodesBefore.before[j], ">")
			}
		}

		fmt.Print(" <- ", dGraph.nodes[i].str, " -> ")

		if len(dGraph.nodes[i].nodesAfter.after) == 0 {
			fmt.Print("nil")
		} else {
			for j := 0; j < len(dGraph.nodes[i].nodesAfter.after); j++ {
				fmt.Print("<", dGraph.nodes[i].nodesAfter.after[j], ">")
			}
		}
	}
	fmt.Println("PRINTED NETWORK")
}

func traverseGraph(dGraph directedGraph) *string {
	var buffer bytes.Buffer

	_, ptr := dGraph.find_element("\u65e5")

	for i := 0; true; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		if len(ptr.nodesAfter.after) == 0 {
			break
		}
		i := rand.Intn(len(ptr.nodesAfter.after))
		ptr = ptr.nodesAfter.after[i]
		if Compare(ptr.str, "\x98") == 0 {
			break
		}
		buffer.WriteString(ptr.str)
		buffer.WriteString(" ")
	}

	str := buffer.String()
	return &str
}

/*
	From our input, we create an array(slice) containing all sentences.
	From here, we create a beginning character (/u65e5) and a terminating
	character (/u65e6) from
*/
func toStringArray(input string) *stringArray {

	sa := &stringArray{}
	sa.array = append(sa.array, "\u65e5")
	inputArray := strings.Split(input, " ")
	for i := 0; i < len(inputArray); i++ {
		sa.array = append(sa.array, inputArray[i])
	}
	sa.array = append(sa.array, "\u65e6")

	return sa
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
