package main

import (
	"fmt"
	"strings"
	"bytes"
	"time"
	"math/rand"
)

type stringBranchBef struct {
	before []*stringBranch;
	score_before []int;
	total_weight int;
}

type stringBranchAft struct {
	score_after []int;
	after []*stringBranch;
	total_weight int;
}

type stringBranch struct{
	str_bef *stringBranchBef;
	str string;
	str_aft *stringBranchAft;
}

type stringArray struct {
	array[] string;
}
//creates a new string for every period
type sentenceArray struct {
	array[] string;
	strA[] stringArray;
	strB stringBranch;
}

type stringTree struct {
	s_branch []*stringBranch;
}

func (str_bef* stringBranchBef) combine_elem(in_bef stringBranchBef) (*stringBranchBef) {

	if_exists := false;
	for i:=0;i<len(in_bef.before);i++ {
		if_exists = false;
		for j:=0;j<len(str_bef.before);j++ {
			if in_bef.before[i].str == str_bef.before[i].str {
				str_bef.score_before[i] = in_bef.score_before[j] + str_bef.score_before[i];
				if_exists = true;
			}
		}
		if !if_exists {
			str_bef.before = append(str_bef.before,in_bef.before[i]);
			str_bef.score_before = append(str_bef.score_before,in_bef.score_before[i]);
		}
	}

	return str_bef;
}

func (strTree *stringTree) build_tree(strA stringArray) (*stringTree) {

	fmt.Print("BUILDING TREE: < ");
	for i:=0;i<len(strA.array);i++ {
		fmt.Print(strA.array[i]," ");
	}
	fmt.Println(" > ");

	for i:=0;i<len(strA.array);i++ {
		str_bra:=&stringBranch{}; str_bra.str_bef = &stringBranchBef{total_weight:0};
		str_bra.str_aft = &stringBranchAft{total_weight:0};
		strTree.s_branch = append(strTree.s_branch,str_bra);
		strTree.s_branch[i].str = strA.array[i];

		// creates our stringTree.
		//Initializes our stringTree to be composed of 's_branch'es.
		//These 's_branch' es are initialized to not only be stringArrays,
		//but be set up in such a way that they are multidirectional hashes.
		//      TODO: add concurrency to this function for quicker times.

		if i == 0 {
			strTree.s_branch[i].str_bef.before = nil;
			strTree.s_branch[i].str_bef.score_before = nil;

		} else if i == (len(strA.array)-1) {
			strTree.s_branch[i].str_bef.before = append(strTree.s_branch[i].str_bef.before,strTree.s_branch[i-1]);
			strTree.s_branch[i].str_bef.score_before = append(strTree.s_branch[i].str_bef.score_before,0);


			strTree.s_branch[i].str_aft.after = nil;
			strTree.s_branch[i].str_aft.score_after = nil;
		} else {
			strTree.s_branch[i].str_bef.before = append(strTree.s_branch[i].str_bef.before,strTree.s_branch[i-1]);
			strTree.s_branch[i].str_bef.score_before = append(strTree.s_branch[i].str_bef.score_before,0);


			strTree.s_branch[i-1].str_aft.after = append(strTree.s_branch[i-1].str_aft.after,strTree.s_branch[i]);
			strTree.s_branch[i-1].str_aft.score_after = append(strTree.s_branch[i-1].str_aft.score_after,0);
		}
	}

	return strTree;
	//if there already were elements of that string
}

func (st stringTree) find_element(in string) (bool,*stringBranch) {
	sb := &stringBranch{};
	ret :=false;
	for i:=0;i<len(st.s_branch);i++ {
		//checks through all of our tree elements to see if there exists
		//the element we're looking for.
		//      TODO: add concurrency to this compare function for quicker times.

		if(Compare(in,st.s_branch[i].str) == 0) {
			ret = true;
			sb = st.s_branch[i]
			return ret,sb;
		}
	}
	sb = nil;
	return ret,sb;
}

func (stbbef stringBranchBef) find_element(in string) (bool, *stringBranch,int) {

	exists:=false;
	sbb:=&stringBranch{};
	i_out:=0;
	for i:=0;i<len(stbbef.before);i++ {
		if stbbef.before[i].str == in {
			exists = true;
			sbb = stbbef.before[i];
			i_out = i;
		}
	}
	return exists,sbb,i_out;
}

func (staft stringBranchAft) find_element(in string) (bool, *stringBranch,int) {
	exists:=false;
	saft:=&stringBranch{};
	i_out :=0;
	for i:=0;i<len(staft.after);i++ {
		if staft.after[i].str == in {
			exists = true;
			saft = staft.after[i];
			i_out = i;
		}

	}
	return exists,saft,i_out;
}

func (st *stringTree) combine_network(in stringTree,weight int,dbg bool) (*stringTree){
	fmt.Print("COMBINING TREES: < ");
	for i:=0;i<len(in.s_branch);i++ {
		fmt.Print(in.s_branch[i].str," ");
	}
	fmt.Println(" > ");


	/*
	    This creates 2 arrays.  The first array stores all elements that exist both in the st
	    structure and the in structure.  Our second array stores all the elements that exist in
	    our in structre, but not our st structure.  Afterwards, this adds all elements into
	    another array called the updateArray.
	*/



	sbArray:=[]*stringBranch{};
	saArray:=[]*stringBranch{};

	for i:=0;i<len(in.s_branch);i++ {

		exists,ptr:=st.find_element(in.s_branch[i].str);
		if exists {
			//array of existing branch nodes
			sbArray = append(sbArray,ptr);

		} else {
			//array of nonexistent branch nodes

			if dbg{
				fmt.Println("Creating element << ",in.s_branch[i].str);
			}

			ptr := &stringBranch{str:in.s_branch[i].str};
			ptr.str_aft = &stringBranchAft{total_weight:0};
			ptr.str_bef = &stringBranchBef{total_weight:0};
			saArray = append(saArray,ptr);
			st.s_branch = append(st.s_branch,ptr);
		}
	}

	for i := 0; i < len(sbArray); i++ {

		exists, ptr := st.find_element(in.s_branch[i].str);

		if !exists {
			fmt.Println("Not all pointers were correctly added to the st array.  \nTerminating on number: %i", i);
			return nil;
		}
		// We check to see if our previous element for the in branch  is already in our st array before branch.
		if i != 0 {
			exists, ptr2, i_out := ptr.str_bef.find_element(in.s_branch[i - 1].str);

			//if there is no element matching the string in our before array, we create another element for this.
			/*
				-------------------------------------------------------------------------------------
							BEFORE BRANCH CREATOR
				-------------------------------------------------------------------------------------
			 */


			if !exists {

				_, ptr2 = st.find_element(in.s_branch[i - 1].str);
				ptr.str_bef.before = append(ptr.str_bef.before, ptr2);
				ptr.str_bef.score_before = append(ptr.str_bef.score_before, weight);
				//update weight
				we:=0;
				for i := 0; i < len(ptr.str_bef.score_before); i++ {
					we=we+ptr.str_bef.score_before[i];
				}

				ptr.str_bef.total_weight = we;



				if dbg{
					fmt.Println("Appending \"", ptr2.str, "\" -> \"", in.s_branch[i].str );
				}

			} else {

				//if it does exist, we simply update the weight of the element
				ptr.str_bef.score_before[i_out] = ptr.str_bef.score_before[i_out] + weight;

				we:=0;
				for i := 0; i < len(ptr.str_bef.score_before); i++ {
					we=we+ptr.str_bef.score_before[i];
				}


				ptr.str_bef.total_weight = we
				if dbg {
					fmt.Println("Updating weight \"",ptr.str,"\" -> \"",ptr.str_bef.before[i_out].str,"\"\tWeight: ",ptr.str_bef.total_weight, "\t",ptr.str_bef.score_before);
				}
			}
		}

		/*
			-------------------------------------------------------------------------------------
						AFTER BRANCH CREATOR
			-------------------------------------------------------------------------------------
		*/

		if i != len(in.s_branch)-1 {

			exists, ptr2, i_out := ptr.str_aft.find_element(in.s_branch[i + 1].str);

			if !exists {

				_, ptr2 = st.find_element(in.s_branch[i + 1].str);
				ptr.str_aft.after = append(ptr.str_aft.after, ptr2);
				ptr.str_aft.score_after = append(ptr.str_aft.score_after, weight);
				//update weight
				we:=0;
				we=0;
				for i := 0; i < len(ptr.str_aft.score_after); i++ {
					we=we+ptr.str_aft.score_after[i];
				}
				ptr.str_aft.total_weight = we;

				if dbg {
					fmt.Println("Appending \"", in.s_branch[i].str , "\" <- \"", ptr2.str);
					fmt.Println("Updating weight \"",ptr.str,"\" <- \"",ptr.str_aft.after[i_out].str,"\"\tWeight: ",ptr.str_aft.total_weight, "\t",ptr.str_aft.score_after);

				}

			} else {
				ptr.str_aft.score_after[i_out] = ptr.str_aft.score_after[i_out] + weight;
				//update weight
				we:=0;
				for i := 0; i < len(ptr.str_aft.score_after); i++ {
					we=we+ptr.str_aft.score_after[i];
				}
				ptr.str_aft.total_weight = we;

				if dbg {
					fmt.Println("Updating weight \"",ptr.str,"\" <- \"",ptr.str_aft.after[i_out].str,"\"\tWeight: ",ptr.str_aft.total_weight, "\t",ptr.str_aft.score_after);
				}
			}
		}
	}


	return st;
}

func (st stringTree) print_network() {

	fmt.Println("PRINTING NETWORK");

	for i:=0;i<len(st.s_branch);i++ {
		fmt.Println();

		if len(st.s_branch[i].str_bef.before) == 0 {
			fmt.Print("nil");
		} else {
			for i:=0;i<len(st.s_branch[i].str_bef.before);i++ {
				fmt.Print("<",st.s_branch[i].str_bef.before[i],">")
			};
		}

		fmt.Print(" <- ",st.s_branch[i].str," -> ");

		if len(st.s_branch[i].str_aft.after) == 0 {
			fmt.Print("nil")
		} else {
			for i:=0;i<len(st.s_branch[i].str_aft.after);i++ {
				fmt.Print("<",st.s_branch[i].str_aft.after[i],">")
			};
		}
	}
}

func  tree_iterator(st stringTree,a int) (*string){
	var buffer bytes.Buffer;

	_,ptr:=st.find_element("\u65e5");

	for i:=0;i<a;i++ {
		rand.Seed(time.Now().UTC().UnixNano());
		i:=rand.Intn(len(ptr.str_aft.after));
		ptr = ptr.str_aft.after[i];
		if Compare(ptr.str,"\x98") == 0 {
			break;
		}
		buffer.WriteString(ptr.str);
		buffer.WriteString(" ");
	}

	str:=buffer.String()
	return &str
}

func tree_iterator1(st stringTree) (*string) {
	var buffer bytes.Buffer;

	_,ptr:=st.find_element("\u65e5");

	for i:=0;true;i++ {
		rand.Seed(time.Now().UTC().UnixNano());
		if len(ptr.str_aft.after) == 0 {
			break;
		}
		i:=rand.Intn(len(ptr.str_aft.after));
		ptr = ptr.str_aft.after[i];
		if Compare(ptr.str,"\x98") == 0 {
			break;
		}
		buffer.WriteString(ptr.str);
		buffer.WriteString(" ");
	}

	str:= buffer.String();
	return &str;
}

func toStringArray(input string) *stringArray{

	sa:= &stringArray{};
	sa.array=append(sa.array,"\u65e5");
	inputArray:=strings.Split(input," ");
	for i:=0;i<len(inputArray);i++ {
		sa.array = append(sa.array,inputArray[i]);
	}
	sa.array=append(sa.array,"\u65e6");

	return sa;
}

func Compare(a string, b string) int {
	if a == b {
		return 0;
	}
	if a < b {
		return -1;
	}
	return 1;
}
