package main

import (
	"bufio"
	"flag"
	"fmt"
	parser "github.com/audreylim/go-markdown"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	// Specify target markdown file with -markdown=<filename>. Otherwise the default name is index.md. Same for the html file.
	fileMD := flag.String("md", "index.md", "a string")
	fileHTML := flag.String("html", "index.html", "a string")
	flag.Parse()

	// Open and read markdown file.
	f, err := os.Open(*fileMD)
	defer f.Close()
	readMDFile := bufio.NewReader(f)
	if err != nil {
		fmt.Println(err)
	}

	// Create new HTML file.
	newfileHTML, err := os.Create(*fileHTML)
	if err != nil {
		fmt.Println(err)
	}
	defer newfileHTML.Close()

	// Parse markdown file.
	p := parser.NewParser(readMDFile)
	go p.Parse()

	// Write to HTML file.
	HTMLTextSlice := []string{}
	go func() {
		for {
			formatter := <-parser.FormatterChn
			stringlit := <-parser.StringlitChn
			switch formatter {
			case "#":
				HTMLTextSlice = append(HTMLTextSlice, "<h1>"+stringlit+"</h1>")
			case "##":
				HTMLTextSlice = append(HTMLTextSlice, "<h2>"+stringlit+"</h2>")
			case "###":
				HTMLTextSlice = append(HTMLTextSlice, "<h3>"+stringlit+"</h3>")
			case "####":
				HTMLTextSlice = append(HTMLTextSlice, "<h4>"+stringlit+"</h4>")
			case "#####":
				HTMLTextSlice = append(HTMLTextSlice, "<h5>"+stringlit+"</h5>")
			case "######":
				HTMLTextSlice = append(HTMLTextSlice, "<h6>"+stringlit+"</h6>")
			//case ">":
			//	HTMLTextSlice = append(HTMLTextSlice, "<blockquote>\n<p>"+stringlit+"</p>\n</blockquote>")
			case "bullet":
				HTMLTextSlice = append(HTMLTextSlice, "<ul>\n"+stringlit+"</ul>")
			case "para":
				HTMLTextSlice = append(HTMLTextSlice, "<p>"+stringlit+"</p>")
			}
		}
	}()
	time.Sleep(time.Second)

	HTMLText := strings.Join(HTMLTextSlice, string('\n'))

	b := []byte(HTMLText)
	ioutil.WriteFile(*fileHTML, b, 0644)
}
