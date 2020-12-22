// aerth (2020)
// MIT License (Open Source)

// bookmark-extract program extracts bookmarks how you want them from browser profile files
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"
)

type Bookmark struct {
	DateAdded    string     `json:"date_added"`
	GUID         string     `json:"guid"`
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	URL          string     `json:"url,omitempty"`
	Children     []Bookmark `json:"children,omitempty"`
	DateModified string     `json:"date_modified,omitempty"`
}
type ChromBookmarkFolder struct {
	Children     []Bookmark `json:"children"`
	DateAdded    string     `json:"date_added"`
	DateModified string     `json:"date_modified"`
	GUID         string     `json:"guid"`
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
}

var helpText = `

Examples:

Print urls, one per line:
./bin/bookmark-extract -q -input ~/.config/google-chrome/Default/Bookmarks -tpl '{{.URL}}'

Print 'name,url', one per line
./bin/bookmark-extract -q -input ~/.config/google-chrome/Default/Bookmarks -tpl '{{.Name}},{{.URL}},,'

Print '<a href="URL">Name</a> ' as one line per folder
./bin/bookmark-extract -noline -q -input ~/.config/google-chrome/Default/Bookmarks -tpl '<a href="{{.URL}}">{{.Name}}</a> '

`

func main() {
	log.SetFlags(0)
	var (
		filename     = os.ExpandEnv("$HOME/.config/google-chrome/Default/Bookmarks")
		outfilename  = "-"
		templatestr  = "{{.}}"
		noNewline    = false
		noLog        = false
		templateHelp = false
	)
	if runtime.GOOS == "windows" {
		filename = "Bookmarks" // ok?!
	}
	flag.StringVar(&filename, "input", filename, "input file (chrome json)")
	flag.StringVar(&outfilename, "output", "", "output file (stdout if empty or -)")
	flag.StringVar(&templatestr, "tpl", templatestr, "template text if contains '{{', otherwise path to template file")
	flag.BoolVar(&noNewline, "noline", noNewline, "dont add new line to every bookmark output")
	flag.BoolVar(&noLog, "q", noLog, "log to /dev/null")
	flag.BoolVar(&templateHelp, "template-help", templateHelp, "show template help and exit")
	flag.Parse()

	if templateHelp {
		fmt.Println(helpText)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	var writer io.Writer = os.Stdout
	if outfilename != "" && outfilename != "-" {
		wf, err := os.Create(outfilename)
		if err != nil {
			log.Fatalln(err)
		}
		defer wf.Close()
	}

	var tpl *template.Template
	if strings.Contains(templatestr, "{{") {
		if !noNewline {
			templatestr += "\n"
		}
		tpl, err = template.New("").Parse(templatestr)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		tpl, err = template.New("").ParseFiles(templatestr)
		if err != nil {
			log.Fatalln(err)
		}
	}

	decoder := json.NewDecoder(f)
	thing := ChromBookmarkFile{}

	if err := decoder.Decode(&thing); err != nil {
		log.Fatalln(err)
	}
	if noLog {
		log.SetOutput(ioutil.Discard)
	}
	log.Println("parsed file successfully")
	template.New("tpl").Parse(templatestr)
	for folderName, folderContents := range thing.Roots {
		log.Println(folderName, folderContents.Name, folderContents.Type, len(folderContents.Children))
		for _, bookmark := range folderContents.Children {
			handleBookmark(bookmark, writer, tpl)
		}

	}
	log.Println("finished")
}

type ChromBookmarkFile struct {
	Checksum string                         `json:"checksum"`
	Roots    map[string]ChromBookmarkFolder `json:"roots"`
	Version  int                            `json:"version"`
}

func handleBookmark(bookmark Bookmark, output io.Writer, tpl *template.Template) {
	if true {
		if err := tpl.Execute(output, bookmark); err != nil {
			log.Printf("error:", err)
		}
	} else {
		fmt.Println(bookmark.Name, bookmark.Type, bookmark.URL)
	}
	for _, v := range bookmark.Children {
		handleBookmark(v, output, tpl)
	}
}
