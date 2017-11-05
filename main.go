package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jzelinskie/geddit"
)

var api *anaconda.TwitterApi

type redditParams struct {
	redditBot         *geddit.LoginSession
	subredditOptions  *geddit.ListingOptions
	redditSubmissions *[]*geddit.Submission
}

type wordArray struct {
	array []string
}

//creates a new string for every period
type sentenceArray struct {
	array []string
	strA  []wordArray
	strB  node
}

var DEBUG = false

func main() {
	//Initializes twitter and reddit connection from the conf.go configurations
	api = initializeTwitter()
	redditParameters := (*initializeReddit())
L:
	//Acquire submissions from "totallynotrobots" subreddit.
	submissions, err := redditParameters.redditBot.SubredditSubmissions(SUBREDDIT_TO_SCRAPE, geddit.HotSubmissions, (*redditParameters.subredditOptions))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	redditParameters.redditSubmissions = &submissions

	sentenceArray := createStringArray(redditParameters)
	mainGraph := buildInitialGraph(sentenceArray)

	/*
		Regular expression to filter out any characters that aren't alphanumeric
	*/
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]")

	/*
		i is used to count the amount of posts that are made each day.  After
		posting (86400/POST_RATE) posts, we refresh, meaning we've posted a
		days worth of tweets, refreshing each day.
	*/
	i := 0

	for true {
		var post string
		for true {
			subgraph := traverseGraph(mainGraph)
			post = subgraph.iterateGraph()
			//post = *traverseGraph(*mainGraph)
			if utf8.RuneCountInString(post) > 140 {
			} else {
				break
			}
		}

		processedString := reg.ReplaceAllString(post, "")
		api.PostTweet(processedString, nil)
		fmt.Println("Tweet posted: ", post)

		time.Sleep(POST_RATE * time.Second)
		i++
		if i >= (86400 / int(POST_RATE)) {
			goto L
			i = 0
		}
	}
}

func initializeTwitter() *anaconda.TwitterApi {

	//retrieves constants from conf.go to initiate twitter connection.
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(TOKEN, TOKEN_SECRET)
	return api

}

func initializeReddit() *redditParams {
	redditConnectionAttempts := 0
L:
	session, err := geddit.NewLoginSession(
		REDDIT_USERNAME,
		REDDIT_PASSWORD,
		REDDIT_BOT_NAME,
	)

	if REDDIT_BOT_NAME == "" {
		fmt.Println("No bot name was specified.  Please check the conf.go file for configurations.")
	}
	if REDDIT_USERNAME == "" {
		fmt.Println("No username was specified.  Please check the conf.go file for configuration")
	}
	if REDDIT_PASSWORD == "" {
		fmt.Println("No password was specified.  Please check the conf.go file for configuration")
	}
	if err != nil {
		fmt.Println("There was an error in communication with reddit.")
		fmt.Println(err)
		//if 3 attempts have failed, we exit the program.
		if redditConnectionAttempts == 3 {
			fmt.Println("UNABLE TO CONNECT TO REDDIT")
			os.Exit(-1)
		}
		redditConnectionAttempts = redditConnectionAttempts + 1
		time.Sleep(2 * time.Second)
		goto L
	}
	options := &geddit.ListingOptions{Limit: THREAD_SAMPLE_COUNT}
	parameters := &redditParams{redditBot: session, subredditOptions: options}
	return parameters
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

/*
	returns an array of everything written, separated by sentences (periods)
	and the end of posts
*/

func createStringArray(parameters redditParams) []*wordArray {
	//sentenceArray
	words := make([]*wordArray, 1)

	for k := 0; k < len(*(parameters.redditSubmissions)); k++ {
		comments, err := parameters.redditBot.Comments((*parameters.redditSubmissions)[k])
		if err != nil {
			fmt.Println("COULD NOT RETRIEVE COMMENTS FROM SUBREDDIT")
			os.Exit(-1)
		}
		if DEBUG {
			for i := 0; i < len(comments); i++ {
				fmt.Println(comments[i].Body)
			}
		}

		for i := 0; i < len(comments); i++ {
			array_of_comment_sentences := strings.Split(comments[i].Body, ".")
			for j := 0; j < len(array_of_comment_sentences); j++ {
				if strings.Contains(array_of_comment_sentences[j], "/") {
					break
				}
				words = append(words, toWordArray(array_of_comment_sentences[j]))
			}
		}
	}
	if DEBUG {
		fmt.Print("Creating String array: ")
		for i := 0; i < len(words); i++ {
			fmt.Println(words[i])
		}
	}

	return words
}

func buildInitialGraph(words []*wordArray) *directedGraph {
	mainGraph := &directedGraph{}
	for i := 0; i < len(words); i++ {
		secondaryGraph := &directedGraph{}
		if words[i] != nil {
			secondaryGraph.buildGraph(*words[i])
			mainGraph.combineGraphs(*secondaryGraph)
		}
	}
	return mainGraph
}
