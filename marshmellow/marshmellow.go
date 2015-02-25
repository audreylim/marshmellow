package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/audreylim/go-marshmellow"
)

func main() {
	var HTMLTextSlice []string

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
	}
}
