**Directed Graph Twitter Bot**

**Authors: Brandon Twilley**

Description: this twitter bot uses directed graphs to reorganize sentences to mostly meaningless spaghetti phrases, but occasionally says something funny.  Try it out with different input data or combine directed graphs to create strange sentences.

System requirements:

golang > Written using go version 1.8.1

golang twitter API Plugin (github.com/ChimeraCoder/anaconda)

golang reddit API plugin (github.com/jzelinskie/geddit)

**To aquire these:**

`go get github.com/jzelinskie/geddit`

`go get github.com/ChimeraCoder/anaconda`

*Instructions: *
1.  install both libraries required.
2.  acquire twitter API credentials (https://apps.twitter.com/)
3.  enter your twitter credentials on conf_sample.go and input your reddit username, password, and bot name.  Select your subreddit you want to use as sample data, and run file.
4.  enjoy.

*Todo: *
1.  Build hash structure for directed graph to search through nodes more efficiently.
2.  Build database functionality.
3.  Integrate with vis.js in order to visualize directed graph.

**EXAMPLE**

The bot can be seen running at the URL: `www.twitter.com/justbotthings1`