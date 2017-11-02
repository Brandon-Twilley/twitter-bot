package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type arrayOfNodesBefore struct {
	before       []*node
	score_before []int
	total_weight int
}

type arrayOfNodesAfter struct {
	score_after  []int
	after        []*node
	total_weight int
}

type node struct {
	nodesBefore *arrayOfNodesBefore
	str         string
	nodesAfter  *arrayOfNodesAfter
}

type stringArray struct {
	array []string
}

//creates a new string for every period
type sentenceArray struct {
	array []string
	strA  []stringArray
	strB  node
}

type directedGraph struct {
	nodes []*node
}

func (nodesBefore *arrayOfNodesBefore) combine_elem(in_bef arrayOfNodesBefore) *arrayOfNodesBefore {

	if_exists := false
	for i := 0; i < len(in_bef.before); i++ {
		if_exists = false
		for j := 0; j < len(nodesBefore.before); j++ {
			if in_bef.before[i].str == nodesBefore.before[i].str {
				nodesBefore.score_before[i] = in_bef.score_before[j] + nodesBefore.score_before[i]
				if_exists = true
			}
		}
		if !if_exists {
			nodesBefore.before = append(nodesBefore.before, in_bef.before[i])
			nodesBefore.score_before = append(nodesBefore.score_before, in_bef.score_before[i])
		}
	}

	return nodesBefore
}

func (currentGraph *directedGraph) buildGraph(strA stringArray) *directedGraph {
	if DEBUG {
		fmt.Print("BUILDING TREE: < ")
	}
	for i := 0; i < len(strA.array); i++ {
		if DEBUG {
			fmt.Print(strA.array[i], " ")
		}
	}
	if DEBUG {
		fmt.Println(" > ")
	}

	for i := 0; i < len(strA.array); i++ {
		str_bra := &node{}
		str_bra.nodesBefore = &arrayOfNodesBefore{total_weight: 0}
		str_bra.nodesAfter = &arrayOfNodesAfter{total_weight: 0}
		currentGraph.nodes = append(currentGraph.nodes, str_bra)
		currentGraph.nodes[i].str = strA.array[i]

		// creates our directedGraph.
		//Initializes our directedGraph to be composed of 'nodes'es.
		//These 'nodes' es are initialized to not only be stringArrays,
		//but be set up in such a way that they are multidirectional hashes.
		//      TODO: add concurrency to this function for quicker times.

		if i == 0 {
			currentGraph.nodes[i].nodesBefore.before = nil
			currentGraph.nodes[i].nodesBefore.score_before = nil

		} else if i == (len(strA.array) - 1) {
			currentGraph.nodes[i].nodesBefore.before = append(currentGraph.nodes[i].nodesBefore.before, currentGraph.nodes[i-1])
			currentGraph.nodes[i].nodesBefore.score_before = append(currentGraph.nodes[i].nodesBefore.score_before, 0)

			currentGraph.nodes[i].nodesAfter.after = nil
			currentGraph.nodes[i].nodesAfter.score_after = nil
		} else {
			currentGraph.nodes[i].nodesBefore.before = append(currentGraph.nodes[i].nodesBefore.before, currentGraph.nodes[i-1])
			currentGraph.nodes[i].nodesBefore.score_before = append(currentGraph.nodes[i].nodesBefore.score_before, 0)

			currentGraph.nodes[i-1].nodesAfter.after = append(currentGraph.nodes[i-1].nodesAfter.after, currentGraph.nodes[i])
			currentGraph.nodes[i-1].nodesAfter.score_after = append(currentGraph.nodes[i-1].nodesAfter.score_after, 0)
		}
	}

	return currentGraph
	//if there already were elements of that string
}

func (st directedGraph) find_element(in string) (bool, *node) {
	sb := &node{}
	ret := false
	for i := 0; i < len(st.nodes); i++ {
		//checks through all of our tree elements to see if there exists
		//the element we're looking for.
		//      TODO: add concurrency to this compare function for quicker times.

		if Compare(in, st.nodes[i].str) == 0 {
			ret = true
			sb = st.nodes[i]
			return ret, sb
		}
	}
	sb = nil
	return ret, sb
}

func (stbbef arrayOfNodesBefore) find_element(in string) (bool, *node, int) {

	exists := false
	sbb := &node{}
	i_out := 0
	for i := 0; i < len(stbbef.before); i++ {
		if stbbef.before[i].str == in {
			exists = true
			sbb = stbbef.before[i]
			i_out = i
		}
	}
	return exists, sbb, i_out
}

func (staft arrayOfNodesAfter) find_element(in string) (bool, *node, int) {
	exists := false
	saft := &node{}
	i_out := 0
	for i := 0; i < len(staft.after); i++ {
		if staft.after[i].str == in {
			exists = true
			saft = staft.after[i]
			i_out = i
		}

	}
	return exists, saft, i_out
}

func (st *directedGraph) combineGraphs(in directedGraph, weight int) *directedGraph {
	if DEBUG {
		fmt.Print("COMBINING TREES: < ")
		for i := 0; i < len(in.nodes); i++ {
			fmt.Print(in.nodes[i].str, " ")
		}
		fmt.Println(" > ")
	}

	/*
	   This creates 2 arrays.  The first array stores all elements that exist both in the st
	   structure and the in structure.  Our second array stores all the elements that exist in
	   our in structre, but not our st structure.  Afterwards, this adds all elements into
	   another array called the updateArray.
	*/

	sbArray := []*node{}
	saArray := []*node{}

	for i := 0; i < len(in.nodes); i++ {

		exists, ptr := st.find_element(in.nodes[i].str)
		if exists {
			//array of existing branch nodes
			sbArray = append(sbArray, ptr)

		} else {
			//array of nonexistent branch nodes

			if DEBUG {
				fmt.Println("Creating element << ", in.nodes[i].str)
			}

			ptr := &node{str: in.nodes[i].str}
			ptr.nodesAfter = &arrayOfNodesAfter{total_weight: 0}
			ptr.nodesBefore = &arrayOfNodesBefore{total_weight: 0}
			saArray = append(saArray, ptr)
			st.nodes = append(st.nodes, ptr)
		}
	}

	for i := 0; i < len(sbArray); i++ {

		exists, ptr := st.find_element(in.nodes[i].str)

		if !exists {
			fmt.Println("Not all pointers were correctly added to the st array.  \nTerminating on number: %i", i)
			return nil
		}
		// We check to see if our previous element for the in branch  is already in our st array before branch.
		if i != 0 {
			exists, ptr2, i_out := ptr.nodesBefore.find_element(in.nodes[i-1].str)

			//if there is no element matching the string in our before array, we create another element for this.
			/*
				-------------------------------------------------------------------------------------
							BEFORE BRANCH CREATOR
				-------------------------------------------------------------------------------------
			*/

			if !exists {

				_, ptr2 = st.find_element(in.nodes[i-1].str)
				ptr.nodesBefore.before = append(ptr.nodesBefore.before, ptr2)
				ptr.nodesBefore.score_before = append(ptr.nodesBefore.score_before, weight)
				//update weight
				we := 0
				for i := 0; i < len(ptr.nodesBefore.score_before); i++ {
					we = we + ptr.nodesBefore.score_before[i]
				}

				ptr.nodesBefore.total_weight = we

				if DEBUG {
					fmt.Println("Appending \"", ptr2.str, "\" -> \"", in.nodes[i].str)
				}

			} else {

				//if it does exist, we simply update the weight of the element
				ptr.nodesBefore.score_before[i_out] = ptr.nodesBefore.score_before[i_out] + weight

				we := 0
				for i := 0; i < len(ptr.nodesBefore.score_before); i++ {
					we = we + ptr.nodesBefore.score_before[i]
				}

				ptr.nodesBefore.total_weight = we
				if DEBUG {
					fmt.Println("Updating weight \"", ptr.str, "\" -> \"", ptr.nodesBefore.before[i_out].str, "\"\tWeight: ", ptr.nodesBefore.total_weight, "\t", ptr.nodesBefore.score_before)
				}
			}
		}

		/*
			-------------------------------------------------------------------------------------
						AFTER BRANCH CREATOR
			-------------------------------------------------------------------------------------
		*/

		if i != len(in.nodes)-1 {

			exists, ptr2, i_out := ptr.nodesAfter.find_element(in.nodes[i+1].str)

			if !exists {

				_, ptr2 = st.find_element(in.nodes[i+1].str)
				ptr.nodesAfter.after = append(ptr.nodesAfter.after, ptr2)
				ptr.nodesAfter.score_after = append(ptr.nodesAfter.score_after, weight)
				//update weight
				we := 0
				we = 0
				for i := 0; i < len(ptr.nodesAfter.score_after); i++ {
					we = we + ptr.nodesAfter.score_after[i]
				}
				ptr.nodesAfter.total_weight = we

				if DEBUG {
					fmt.Println("Appending \"", in.nodes[i].str, "\" <- \"", ptr2.str)
					fmt.Println("Updating weight \"", ptr.str, "\" <- \"", ptr.nodesAfter.after[i_out].str, "\"\tWeight: ", ptr.nodesAfter.total_weight, "\t", ptr.nodesAfter.score_after)

				}

			} else {
				ptr.nodesAfter.score_after[i_out] = ptr.nodesAfter.score_after[i_out] + weight
				//update weight
				we := 0
				for i := 0; i < len(ptr.nodesAfter.score_after); i++ {
					we = we + ptr.nodesAfter.score_after[i]
				}
				ptr.nodesAfter.total_weight = we

				if DEBUG {
					fmt.Println("Updating weight \"", ptr.str, "\" <- \"", ptr.nodesAfter.after[i_out].str, "\"\tWeight: ", ptr.nodesAfter.total_weight, "\t", ptr.nodesAfter.score_after)
				}
			}
		}
	}

	return st
}

func (st directedGraph) printGraph() {

	fmt.Println("PRINTING NETWORK")

	for i := 0; i < len(st.nodes); i++ {
		fmt.Println()

		if len(st.nodes[i].nodesBefore.before) == 0 {
			fmt.Print("nil")
		} else {
			for j := 0; j < len(st.nodes[i].nodesBefore.before); j++ {
				fmt.Print("<", st.nodes[i].nodesBefore.before[j], ">")
			}
		}

		fmt.Print(" <- ", st.nodes[i].str, " -> ")

		if len(st.nodes[i].nodesAfter.after) == 0 {
			fmt.Print("nil")
		} else {
			for j := 0; j < len(st.nodes[i].nodesAfter.after); j++ {
				fmt.Print("<", st.nodes[i].nodesAfter.after[j], ">")
			}
		}
	}
	fmt.Println("PRINTED NETWORK")
}

func traverseGraph(st directedGraph) *string {
	var buffer bytes.Buffer

	_, ptr := st.find_element("\u65e5")

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
	}

	str := buffer.String()
	return &str
}

func toStringArray(input string) *stringArray {
	/*
		From our input, we create an array(slice) containing all sentences.
		From here, we create a beginning character (/u65e5) and a terminating
		character (/u65e6) from
	*/
	sa := &stringArray{}
	sa.array = append(sa.array, "\u65e5")
	inputArray := strings.Split(input, " ")
	for i := 0; i < len(inputArray); i++ {
		sa.array = append(sa.array, inputArray[i])
	}
	sa.array = append(sa.array, "\u65e6")

	return sa
}

func Compare(a string, b string) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}
