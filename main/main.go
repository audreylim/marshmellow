package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	md "github.com/audreylim/go-markdown"
)

func main() {
	var HTMLTextSlice []string
	timebreaker := make(chan string)

	r1 := regexp.MustCompile("[a-zA-Z0-9]+.md")
	r2 := regexp.MustCompile("[a-zA-Z0-9]+")

	flag.Parse()
	mdfiles := flag.Args()

	for _, v := range mdfiles {
		if !r1.MatchString(v) {
			fmt.Printf("%s is not a markdown file", v)
		} else {
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
			p := md.NewParser(readMDFile)
			go func() {
				if v := p.Parse(); v == "eof" {
					timebreaker <- "x"
				}
			}()

			// Write to HTML file.
			go func() {
				for {
					formatter := <-md.FormatterChn
					stringlit := <-md.StringlitChn
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
					case "bullet":
						HTMLTextSlice = append(HTMLTextSlice, "<ul>\n"+stringlit+"</ul>")
					case "para":
						HTMLTextSlice = append(HTMLTextSlice, "<p>"+stringlit+"</p>")
					}

				}
			}()

			select {
			case _, ok := <-timebreaker:
				if !ok {
					time.Sleep(time.Second * 5)
				}
			}

			fmt.Println(HTMLTextSlice)
			HTMLText := strings.Join(HTMLTextSlice, string('\n'))
			b := []byte(HTMLText)
			ioutil.WriteFile(fileHTML, b, 0644)
		}
	}
}
