package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"fmt"
	"os"
	"net/http"
	"io"
	"io/ioutil"
	"net"
	"strings"
)


func main() {

	LISTEN_PORT := "9000"
	SEND_PORT := "9001"
	
	photonet := false
	always_safe:= true
	is_safe := false;
	for {
		str := listen(LISTEN_PORT);
		str2 := strings.Split(str,"`");
		if Compare(str2[0], "PING") == 0 {
			url := ""
			
			if (photonet) {
				url = scrape_photonet();
				is_safe = true;
				photonet = false;
			} else {
				L: 
				url,is_safe = scrape_danbooru();
				
				if always_safe&&!is_safe {
					goto L;
				}
				
				fmt.Printf(str);
				if is_safe {
					fmt.Println("\nour image is safe")
				} else {
					fmt.Println("\nour image isn't safe");
				}
				photonet=true;
			}
			
			download("download.jpg",url);
			if is_safe {
				send("download.jpg`true`"+str2[1],SEND_PORT);
			} else {
				send("download.jpg`false`"+str2[1],SEND_PORT);
			}
		} else {
			always_safe = !always_safe;
		}
	}
}

func download(path string, url string) {
	img,_ := os.Create(path);
	defer img.Close();
	
	resp,_ := http.Get(url);
	defer resp.Body.Close();
	
	b,_ := io.Copy(img,resp.Body);
	fmt.Println("Filesize: ",b, " bytes.");
}

func scrape_danbooru() (str string, is_safe bool) {
	doc, err := goquery.NewDocument("http://danbooru.donmai.us/posts/random")
	if err != nil {
		log.Fatal(err)
	}
	str = "";

	doc.Find("ul").Each(func (i int, s *goquery.Selection) {
		st,_ := s.Attr("itemtype")
		if st != "" {
			s.Find("li").Each(func (j int, t *goquery.Selection) {
				//fmt.Printf("%iTHIS IS THE TEXT QUERIED: %s\n", j, t.Text())
				if j == 5 {
					if t.Text() == "Rating: Safe" {
						is_safe = true;
					} else {
						is_safe = false;
					}
				}

			})
		}

	})


	doc.Find("img").Each(func (i int, s *goquery.Selection) {
		st,_ := s.Attr("src");

		str="http://danbooru.donmai.us" + st;
	})

	return str,is_safe;
}

func scrape_photonet() (str string){
	doc, err := goquery.NewDocument("http://photo.net/photodb/random-photo?category=NoNudes")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		st,_ := s.Attr("src");

		fmt.Printf(str);
		fmt.Printf("\n");
		if i == 2 {
			str = st;
		}
	})
	return str
}
func listen(port string) (out string) {
	L:
	conn, err := net.Dial("tcp", "localhost:" + port)
	if err != nil {
		goto L
	}
	
	defer conn.Close();
	
	bs,_ := ioutil.ReadAll(conn)
	return (string(bs))
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

func Compare(a string, b string) int {
	if a == b {
		return 0;
	}
	if a < b {
		return -1;
	}
	return 1;
}
