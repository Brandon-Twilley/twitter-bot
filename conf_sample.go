package main

import "time"

//Twitter's consumer key
const CONSUMER_KEY string = ""

//twitters consumer secret hash string
const CONSUMER_SECRET string = ""

//twitter's token
const TOKEN string = ""

//twitter's token secret
const TOKEN_SECRET string = ""

//basic reddit credentials
const REDDIT_USERNAME string = ""
const REDDIT_PASSWORD string = ""
const REDDIT_BOT_NAME string = "MARKOV POLO"

const SUBREDDIT_TO_SCRAPE string = ""

//number of threads you retrieve.
const THREAD_SAMPLE_COUNT int = 10
const POST_RATE time.Duration = 60 //The rate this bot will post to twitter (in seconds)
//this bot will refresh every day

//NOTE, THIS ONLY SCRAPES THE COMMENTS OF EACH SUBREDDIT.
