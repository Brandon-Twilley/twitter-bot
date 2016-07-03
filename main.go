package main

import(
	"github.com/jzelinskie/geddit"
	"github.com/ChimeraCoder/anaconda"
	"time"
	"strings"
	"fmt"
	"unicode/utf8"
	"io/ioutil"
	"net/url"
	"strconv"
	"encoding/base64"
	"net"
	"io"
)

var api *anaconda.TwitterApi;

func main() {
	
	api = initializeTwitter();
	api.PostTweet("this is an api test",nil);

	L:
	twit_bot := initializeTwitter();
	redd_bot,subOpts,err := initializeReddit();

	if err != nil {
		fmt.Errorf(err.Error());
		time.Sleep(2*time.Second)
		goto L
	}

	submissions,err := redd_bot.SubredditSubmissions(SUBREDDIT_TO_SCRAPE, geddit.HotSubmissions,*subOpts)
	if err != nil {
		fmt.Errorf(err.Error());
	}
	sArray := createStringArray(redd_bot, subOpts,submissions)
	submissions, err = redd_bot.SubredditSubmissions("totallynotrobots", geddit.HotSubmissions, *subOpts)

	sArray2 := createStringArray(redd_bot, subOpts, submissions)

	sArray = append(sArray, sArray2...)

	mainTree := buildStringTree(sArray);



	//PostImageToTwitter("output.jpg","Testing API", twit_bot);

	//post image end
	i := 0

	for true {
		var str string;
		for true {
			str=*tree_iterator1(*mainTree)
			if utf8.RuneCountInString(str) >140 {

			} else {
				break;
			}
		}
		str = str[0:len(str)-1]


		time.Sleep(POST_RATE * time.Second);
		
		path,is_sfw := build_image(str);
				
		PostImageToTwitter(path,str,is_sfw,twit_bot);
		fmt.Println("Tweet posted: ",str);
		i++
		if i >= (86400/int(POST_RATE)) {
			goto L
			i=0;
		}
	}
}

func build_image(post string) (path string, is_sfw bool) {
	
	is_sfw = false;
	fmt.Println("Pinging scraper: ")
	send("PING`"+post, "9000");
	
	str := listen("9002");
	str2 := strings.Split(str,"`");
	path = str2[0];
	if Compare(str2[1],"true")==0 {
		is_sfw=true;
	} else {
		is_sfw=false;
	}
	
	return path, is_sfw;
}

func send(in string, port string) {
	ln, err := net.Listen("tcp", "localhost:"+port);
	if err != nil {
		panic(err)
	}
	
	defer ln.Close();
	
	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	
	io.WriteString(conn, fmt.Sprint(in));
	
	conn.Close();
}

func listen(port string) (out string) {
	L:
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		goto L
	}
	
	defer conn.Close();
	
	bs,_ := ioutil.ReadAll(conn)
	return (string(bs))
}

func PostImageToTwitter(path string, post string,is_sfw bool, twit_bot *anaconda.TwitterApi) {
	data, err := ioutil.ReadFile(path)
	if err != nil{
		fmt.Println(err)
	}

	mediaResponse, err := twit_bot.UploadMedia(base64.StdEncoding.EncodeToString(data))
	if err != nil {
		fmt.Println(err)
	}

	v := url.Values{}
	v.Set("media_ids", strconv.FormatInt(mediaResponse.MediaID, 10))
	v.Add("possibly_sensitive",strconv.FormatBool(is_sfw));
	result, err := twit_bot.PostTweet(post, v)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

//HOW TO TWEET: api.PostTweet("<TEXT>", nil);

func initializeTwitter() (*anaconda.TwitterApi) {


	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(TOKEN, TOKEN_SECRET);

	return api

}

func initializeReddit() (session *geddit.LoginSession, subOpts *geddit.ListingOptions, err error) {
	session, err = geddit.NewLoginSession(
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
	}

	subOpts = &geddit.ListingOptions{Limit:THREAD_SAMPLE_COUNT}

	return session, subOpts, err;
}


func createStringArray(session *geddit.LoginSession, subOpts *geddit.ListingOptions, submissions []*geddit.Submission) (sArray2 []*stringArray) {

	sArray2 =make([]*stringArray,1);

	for k:=0;k<len(submissions);k++ {
		comments,_ :=session.Comments(submissions[k]);


		for i:=0;i<len(comments);i++ {
			fmt.Println(comments[i].Body)
		}


		for i:=0;i<len(comments);i++ {

			cmt:=strings.Split(comments[i].Body,".");
			for j:=0;j<len(cmt);j++ {

				if strings.Contains(cmt[j],"/") {
					break;
				}
				
				//cmt[j] = strings.toLower(??);
				
				sArray2 = append(sArray2,toStringArray(cmt[j]));
			}
		}
	}
	return sArray2;
}

func buildStringTree(sArray []*stringArray) (*stringTree) {
	mainTree := &stringTree{};
	for i:=0;i<len(sArray);i++ {
		sTree2 := &stringTree{};
		if sArray[i] != nil {
			sTree2.build_tree(*sArray[i]);
			mainTree.combine_network(*sTree2,1,false);
		}
	}
	return mainTree
}

