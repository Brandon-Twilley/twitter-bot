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

func createNode() *node {
	currentNode := &node{}
	currentNode.nodesBefore = &arrayOfNodesBefore{}
	currentNode.nodesAfter = &arrayOfNodesAfter{}
	currentNode.str = ""

	return currentNode
}

/*
	This function creates our directedGraph.
	Our directedGraph is to be composed of 'nodes'.
	These nodes are linked by words directly before and after them, so they
	would then be a doubly linked directed graph.
*/
func (currentGraph *directedGraph) buildGraph(wordA wordArray) *directedGraph {
	if DEBUG {
		fmt.Print("Building Graph:")
	}

	for i := 0; i < len(wordA.array); i++ {
		currentNode := createNode()
		currentGraph.nodes = append(currentGraph.nodes, currentNode)
		currentGraph.nodes[i].str = wordA.array[i]

		if i == 0 {
			currentGraph.nodes[i].nodesBefore.before = nil

		} else if i == (len(wordA.array) - 1) {
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
*/
func (dGraph directedGraph) find_element(in string) (bool, *node, int) {
	sb := &node{}
	ret := false
	for i := 0; i < len(dGraph.nodes); i++ {
		//checks through all of our tree elements to see if there exists
		//the element we're looking for.
		//      TODO: add concurrency to this compare function for quicker times.

		if Compare(in, dGraph.nodes[i].str) == 0 {
			ret = true
			sb = dGraph.nodes[i]
			return ret, sb, i
		}
	}
	sb = nil
	return ret, sb, 0
}

/*
	This function searches through the list of existing elements within all
	words that were found to precede the word we're looking at.  If it exists,
	the node (structure holding the word) is returned and our function also
	returns a true.  If no word was found, we return a false in our boolean
	and a null on our node pointer.
*/
func (stbbef arrayOfNodesBefore) find_element(in string) (bool, *node, int) {

	sb := &node{}
	ret := false
	for i := 0; i < len(stbbef.before); i++ {
		//checks through all of our tree elements to see if there exists
		//the element we're looking for.
		//      TODO: add concurrency to this compare function for quicker times.

		if Compare(in, stbbef.before[i].str) == 0 {
			ret = true
			sb = stbbef.before[i]
			return ret, sb, i
		}
	}
	sb = nil
	return ret, sb, 0
}

/*
	This function searches through the list of existing elements within all
	words that were found after the word we're looking at.  If it exists,
	the node (structure holding the word) is returned and our function also
	returns a true.  If no word was found, we return a false in our boolean
	and a null on our node pointer.
*/

func (staft arrayOfNodesAfter) find_element(in string) (bool, *node, int) {
	sb := &node{}
	ret := false
	for i := 0; i < len(staft.after); i++ {
		//checks through all of our tree elements to see if there exists
		//the element we're looking for.
		//      TODO: add concurrency to this compare function for quicker times.

		if Compare(in, staft.after[i].str) == 0 {
			ret = true
			sb = staft.after[i]
			return ret, sb, i
		}
	}
	sb = nil
	return ret, sb, 0
}

/*
   This creates 2 arrays.  The first array stores all elements that exist both in the dGraph
   structure and the in structure.  Our second array stores all the elements that exist in
   our in structre, but not our dGraph structure.  Afterwards, this adds all elements into
   another array called the updateArray.
*/
func (dGraph *directedGraph) combineGraphs(in directedGraph) *directedGraph {
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
		exists, ptr, _ := dGraph.find_element(in.nodes[i].str)
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

		exists, ptr, _ := dGraph.find_element(in.nodes[i].str)

		if !exists {
			fmt.Println("Not all pointers were correctly added to the dGraph array.  \nTerminating on number: %i", i)
			return nil
		}
		// We check to see if our previous element for the in branch  is already in our dGraph array before branch.
		if i != 0 {
			exists, ptr2, _ := ptr.nodesBefore.find_element(in.nodes[i-1].str)

			//if there is no element matching the string in our before array, we create another element for this.
			/*
				-------------------------------------------------------------------------------------
							BEFORE BRANCH CREATOR
				-------------------------------------------------------------------------------------
			*/

			if !exists {
				_, ptr2, _ = dGraph.find_element(in.nodes[i-1].str)
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

			exists, ptr2, _ := ptr.nodesAfter.find_element(in.nodes[i+1].str)

			if !exists {
				_, ptr2, _ = dGraph.find_element(in.nodes[i+1].str)
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

/*
	from our graph, we create a subgraph (a specific traversal of the graph that doesn't
	have cycles within it).  This is the phrase we return and post to twitter.

	dGraph is our main graph.  subgraph is the graph we're creating
*/
func traverseGraph(dGraph *directedGraph) *directedGraph {
	//create graph and begin our new graph with the initial element.

	_, begin, _ := dGraph.find_element("\u65e5")
	subgraph := &directedGraph{}
	currentMainGraphNode := begin
	var currentSubGraphNode *node
	var previousSubGraphNode *node

	restart := false

	for true {

		exists, _, _ := subgraph.find_element(currentMainGraphNode.str)
		if exists {
			restart = true
			break
		} else {
			//create our currentsubgraphnode
			currentSubGraphNode = createNode()
			currentSubGraphNode.str = currentMainGraphNode.str
			subgraph = subgraph.appendNode(currentSubGraphNode)

			if Compare(currentSubGraphNode.str, begin.str) == 0 {

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
			if Compare(currentMainGraphNode.str, "\x98") == 0 {
				break
			}

		}
	}
	if restart {
		subgraph = traverseGraph(dGraph)
	}
	return subgraph
}

func (subgraph *directedGraph) iterateGraph() string {
	i := 0
	var buffer bytes.Buffer
	var n *node
	n = subgraph.nodes[i]
	for true {

		buffer.WriteString(n.str)
		buffer.WriteString(" ")

		if Compare(n.str, "\u65e6") == 0 {
			break
		}
		if len(n.nodesAfter.after) == 0 {
			break
		} else {
			n = n.nodesAfter.after[0]
		}
	}
	return buffer.String()

}

/*
	adds node to directed graph
*/
func (dGraph directedGraph) appendNode(n *node) *directedGraph {
	dGraph.nodes = append(dGraph.nodes, n)
	return &dGraph
}

/*
	removes node from graph and replaces array slot with the last element
*/
func (dGraph directedGraph) removeNode(n *node) *directedGraph {
	is_found, _, number := dGraph.find_element(n.str)
	if is_found {
		dGraph.nodes[number] = dGraph.nodes[len(dGraph.nodes)-1]
		dGraph.nodes = dGraph.nodes[:len(dGraph.nodes)-1]
	}

	return &dGraph
}

/*
	From our input, we create an array(slice) containing all sentences.
	From here, we create a beginning character (/u65e5) and a terminating
	character (/u65e6) from
*/
func toWordArray(input string) *wordArray {

	sa := &wordArray{}
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
