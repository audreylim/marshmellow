package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/audreylim/marshmellow"
)

func main() {
	flag.Parse()
	c := genMdFiles(flag.Args())
	var wg sync.WaitGroup
	for mdfile := range c {
		wg.Add(1)
		go performParseFile(mdfile, &wg)
	}
	wg.Wait()
}


func genMdFiles(mdfiles []string) <-chan string {
	out := make(chan string)
	r1 := regexp.MustCompile("[a-zA-Z0-9]+.md")
	go func() {
		for _, v := range mdfiles {
			if !r1.MatchString(v) {
				fmt.Printf("%s is not a markdown file", v)
			} else {
				// emit filename out of channel
				out <- v
			}	
		}
		close(out)
	}()
	return out
}

func performParseFile(v string, wg *sync.WaitGroup) {
	defer wg.Done()
	var HTMLTextSlice []string
	r2 := regexp.MustCompile("[a-zA-Z0-9]+")

	filename := r2.FindString(v)
	// Open and read markdown file.
	f, err := os.Open(v)
	defer f.Close()
	readMDFile := bufio.NewReader(f)
	if err != nil {
		fmt.Println(err)
	}

	// Create new HTML file.
	fileHTML := filename + ".html"
	newfileHTML, err := os.Create(fileHTML)
	if err != nil {
		fmt.Println(err)
	}
	defer newfileHTML.Close()

	// Reset slice for next file.
	HTMLTextSlice = []string{}

	// Parse markdown file.
	p := mm.NewParser(readMDFile)
	p.Parse()
	
	// Write to HTML file.
	for i := 0; i < len(p.Formatter); i++ {
		switch p.Formatter[i] {
		case "#":
			HTMLTextSlice = append(HTMLTextSlice, "<h1>"+p.Stringlit[i]+"</h1>")
		case "##":
			HTMLTextSlice = append(HTMLTextSlice, "<h2>"+p.Stringlit[i]+"</h2>")
		case "###":
			HTMLTextSlice = append(HTMLTextSlice, "<h3>"+p.Stringlit[i]+"</h3>")
		case "####":
			HTMLTextSlice = append(HTMLTextSlice, "<h4>"+p.Stringlit[i]+"</h4>")
		case "#####":
			HTMLTextSlice = append(HTMLTextSlice, "<h5>"+p.Stringlit[i]+"</h5>")
		case "######":
			HTMLTextSlice = append(HTMLTextSlice, "<h6>"+p.Stringlit[i]+"</h6>")
		case "bullet":
			HTMLTextSlice = append(HTMLTextSlice, "<ul>\n"+p.Stringlit[i]+"</ul>")
		case "para":
			HTMLTextSlice = append(HTMLTextSlice, "<p>"+p.Stringlit[i]+"</p>")
		}
	}

	HTMLText := strings.Join(HTMLTextSlice, string('\n'))
	b := []byte(HTMLText)
	ioutil.WriteFile(fileHTML, b, 0644)	
}