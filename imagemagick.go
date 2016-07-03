package main

import (
	"gopkg.in/gographics/imagick.v2/imagick"
	"fmt"
	"net"
	"strings"
	"io/ioutil"
	"io"
)

func main() {
	for {
		str := listen("9001");
		str2 := strings.Split(str,"`");
		fmt.Println("Recieved request: \n",str2[0],"\n",str2[2])
		
		// file, is_sfw, string to write
		
		write_to_img(str2[2], str2[0], "output.jpg", true, "NewCenturySchlbk-Bold");
		//brand("input.jpg","output.jpg");
		 
		str = "output.jpg`"+str2[1]+"`"+str2[2];
		send(str,"9002");
	}
}

func write_to_img(text string, file1 string, file2 string, bold bool, font string) {
	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	pw.SetColor("none")
	// Create a new transparent image
	mw.NewImage(1000, 1000, pw)
	// Set up a 72 point white font
	pw.SetColor("white")
	dw.SetFillColor(pw)
	
	dw.SetFont(font);
	
	
	dw.SetFontSize(40)
	// Add a black outline to the text
	pw.SetColor("black")
	dw.SetStrokeColor(pw)
	// Turn antialias on - not sure this makes a difference
	dw.SetTextAntialias(true)
	// Now draw the text
		//dw.Annotation(25,65, text)
			
			//write mystery??Does this work===================================================
			
		str_array := separate_string(text, 15);
		
		for i := 0;i<len(str_array);i++ {
			dw.Annotation(25, float64(i*65+65), str_array[i]);
		}
	
	//dw.Annotation(25, 65, text)
	// Draw the image on to the mw
	mw.DrawImage(dw)
	// Make a copy of the text image
	cw := mw.Clone()
	// Opacity is a real number indicating (apparently) percentage
	mw.ShadowImage(70, 4, 5, 5)
	// Composite the text on top of the shadow
	mw.CompositeImage(cw, imagick.COMPOSITE_OP_OVER, 5, 5)
	cw.Destroy()
	cw = imagick.NewMagickWand()

	// Create a new image the same size as the text image and put a solid colour
	// as its background
	cw.ReadImage(file1);
	// Now composite the shadowed text over the plain background
	cw.CompositeImage(mw, imagick.COMPOSITE_OP_OVER, int(cw.GetImageWidth())/6, int(cw.GetImageHeight())/5)
	// and write the result
	cw.WriteImage(file2)
	
	fmt.Printf("Finished creating image.  Image posted as file: ");
	fmt.Printf(file2);
	fmt.Printf("\n\tOriginal file: ");
	fmt.Println(file1);

}

func separate_string(text string, break_after int) (holder []string) {
	
	holderStr := "";
	
	flag := true;
	
	j := 0;
	
	for i := 0; i<len(text);i++ {
		if flag {
			j = 0;
			flag = false
		}
		j++
		
		substr := text[i:i+1]
		holderStr += substr;

		if j >=break_after {
			if substr == " " {
				holder = append(holder, holderStr);
				holderStr = "";
				flag = true;
			}
		}
		
	}
	holder = append(holder, holderStr);
	return holder;
}

func brand(file_in string, file_out string) {
	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	pw.SetColor("none")
	// Create a new transparent image
	mw.NewImage(1000, 1000, pw)
	// Set up a 72 point white font
	pw.SetColor("white")
	dw.SetFillColor(pw)
	
	dw.SetFont("Inconsolada");
	
	
	dw.SetFontSize(40)
	// Add a black outline to the text
	pw.SetColor("black")
	dw.SetStrokeColor(pw)
	// Turn antialias on - not sure this makes a difference
	dw.SetTextAntialias(true)
	// Now draw the text
		
		dw.Annotation(25, 65, "#justrobotthings");
	
	// Draw the image on to the mw
	mw.DrawImage(dw)
	// Make a copy of the text image
	cw := mw.Clone()
	// Opacity is a real number indicating (apparently) percentage
	mw.ShadowImage(70, 4, 5, 5)
	// Composite the text on top of the shadow
	mw.CompositeImage(cw, imagick.COMPOSITE_OP_OVER, 5, 5)
	cw.Destroy()
	cw = imagick.NewMagickWand()

	// Create a new image the same size as the text image and put a solid colour
	// as its background
	cw.ReadImage(file_in);
	// Now composite the shadowed text over the plain background
	cw.CompositeImage(mw, imagick.COMPOSITE_OP_OVER, int(80), int(cw.GetImageHeight())-80)
	// and write the result
	cw.WriteImage(file_out)
	
	fmt.Printf("Finished creating image.  Image posted as file: ");
	fmt.Printf(file_out);
	fmt.Printf("\n\tOriginal file: ");
	fmt.Println(file_in);
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
