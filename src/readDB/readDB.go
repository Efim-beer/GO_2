package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"unicode/utf8"
)

type Recipes struct {
	XMLName xml.Name `xml:"recipes" json:"-"`
	Cake    []struct {
		Name       string `xml:"name" json:"name"`
		Stovetime  string `xml:"stovetime" json:"time"`
		Ingredient []struct {
			Itemname  string `xml:"itemname" json:"ingredient_name"`
			Itemcount string `xml:"itemcount" json:"ingredient_count"`
			Itemunit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
		} `xml:"ingredients>item" json:"ingredients"`
	} `xml:"cake" json:"cake"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileFormat(fileName string) string {
	var c = utf8.RuneCountInString(fileName)
	if fileName[c-5:] == ".json" {
		return "json"
	} else if fileName[c-4:] == ".xml" {
		return "xml"
	} else {
		panic("Wrong file format")
	}
}

func commentFree(file *[]byte) {
	s := string(*file)
	var x, y int
	x = -1
	r := []rune(s)
	// fmt.Println(r)
	var buf []rune
	for i, _ := range r {
		if r[i] == '/' && r[i+1] == '/' {
			x = i
		}
		if x >= 0 && r[i] == '\n' {
			y = i
			buf = append(r[:x], r[y:]...)
			x = -1
		}
	}

	if len(buf) > 0 {
		newS := string(buf)
		*file = []byte(newS)
	}
}

func main() {
	var f bool

	flag.BoolVar(&f, "f", false, "display path")
	flag.Parse()

	if f {
		fileName := os.Args[2]
		format := fileFormat(fileName)
		file, err := os.ReadFile(fileName)
		check(err)

		recipe := new(Recipes)
		var res []byte
		switch format {
		case "json":
			commentFree(&file)
			err := json.Unmarshal([]byte(file), recipe)
			check(err)
			res, err = xml.MarshalIndent(recipe, "", "    ")
			check(err)
		case "xml":
			err := xml.Unmarshal([]byte(file), recipe)
			check(err)
			res, err = json.MarshalIndent(recipe, "", "    ")
			check(err)
		}
		fmt.Printf("%s\n", res)
	} else {
		panic("not flag -f")
	}
}
