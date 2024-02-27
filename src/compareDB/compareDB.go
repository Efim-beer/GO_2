package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/r3labs/diff/v3"
)

type Recipe struct {
	XMLName xml.Name `xml:"recipes" json:"-" diff:"-"`
	Cake    []struct {
		Name       string `xml:"name" json:"name" diff:"Name"`
		Stovetime  string `xml:"stovetime" json:"time" diff:"Time"`
		Ingredient []struct {
			Itemname  string `xml:"itemname" json:"ingredient_name" diff:"Itemname, identifier"`
			Itemcount string `xml:"itemcount" json:"ingredient_count" diff:"Itemcount"`
			Itemunit  string `xml:"itemunit" json:"ingredient_unit,omitempty" diff:"Itemunit"`
		} `xml:"ingredients>item" json:"ingredients" diff:"-"`
	} `xml:"cake" json:"cake" diff:"-"`
}

type XML Recipe

func (p *XML) Read(file []byte) Recipe {
	err := xml.Unmarshal(file, p)
	check(err)
	return Recipe(*p)
}

type Json Recipe

func (p *Json) Read(file []byte) Recipe {
	err := json.Unmarshal(file, p)
	check(err)
	return Recipe(*p)
}

type DBReader interface {
	Read(file []byte) Recipe
}

func readRecipe(fileName string) Recipe {
	recipe := new(Recipe)
	format := fileFormat(fileName)
	file, err := os.ReadFile(fileName)
	check(err)
	switch format {

	case "xml":
		fileStruct := new(XML)
		recipe1 := fileStruct.Read(file)
		recipe = &recipe1

	case "json":
		commentFree(&file)
		fileStruct := new(Json)
		recipe1 := fileStruct.Read(file)
		recipe = &recipe1
	}
	return *recipe
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

func addRemoveCake(r1 Recipe, r2 Recipe, r string) {
	for i := 0; i < len(r1.Cake); i++ {
		var identCake = false
		for x := 0; x < len(r2.Cake); x++ {
			if r1.Cake[i].Name == r2.Cake[i].Name {
				identCake = true
			}
		}
		if !identCake {
			fmt.Printf("%s cake \"%s\"\n", r, r1.Cake[i].Name)
		}
	}
}

func addRemoveIngridient(r1 Recipe, r2 Recipe, p string) {
	for i := 0; i < len(r1.Cake); i++ {
		if r1.Cake[i].Name == r2.Cake[i].Name {
			for x := 0; x < len(r1.Cake[i].Ingredient); x++ {
				var identIngridient = false
				for y := 0; y < len(r2.Cake[i].Ingredient); y++ {
					if r1.Cake[i].Ingredient[x].Itemname == r2.Cake[i].Ingredient[y].Itemname {
						identIngridient = true
					}
				}
				if !identIngridient {
					fmt.Printf("%s ingredient \"%s\" for cake \"%s\"\n", p,
						r1.Cake[i].Ingredient[x].Itemname, r1.Cake[i].Name)
				}
			}
		}
	}
}

func changeIngredient(r1 Recipe, r2 Recipe) {
	for i := 0; i < len(r1.Cake); i++ {
		if r1.Cake[i].Name == r2.Cake[i].Name {
			for x := 0; x < len(r1.Cake[i].Ingredient); x++ {
				for y := 0; y < len(r2.Cake[i].Ingredient); y++ {
					if r1.Cake[i].Ingredient[x].Itemname == r2.Cake[i].Ingredient[y].Itemname {
						ch, _ := diff.Diff(r1.Cake[i].Ingredient[x], r2.Cake[i].Ingredient[y])
						printChangeIngredient(ch, r1.Cake[i].Name, r1.Cake[i].Ingredient[x].Itemname)
					}
				}
			}
		}
	}
}

func printChangeIngredient(changelog diff.Changelog, cake, ingName string) {
	for _, change := range changelog {
		count := "count "
		if change.Path[0] != "Itemcount" {
			count = ""
		}
		if change.To != "" && change.From != "" {
			fmt.Printf("CHANGED unit %sfor ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
				count, ingName, cake, change.To, change.From)
		} else if change.To == "" {
			fmt.Printf("REMOVED unit %s\"%s\" for ingredient \"%s\" for cake \"%s\"\n",
				count, change.From, ingName, cake)
		} else {
			fmt.Printf("ADDED unit %s\"%s\" for ingredient \"%s\" for cake \"%s\"\n",
				count, change.To, ingName, cake)
		}
	}
}

func main() {
	f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	old := f1.Bool("old", false, "take old data")
	neww := f1.Bool("new", false, "take new data")
	if len(os.Args) == 5 {
		f1.Parse(os.Args[1:])
		f1.Parse(os.Args[3:])
	} else {
		panic("Not correct args.")
	}

	if *old && *neww {
		if os.Args[1] != "--old" && os.Args[1] != "-old" {
			panic("Wrong order of args. Older is first.")
		} else if os.Args[1] == "-old" {
			panic("Use '--old', not '-old'.")
		} else if os.Args[3] == "-new" {
			panic("Use '--new', not '-new'.")
		}

		fileNameOld := os.Args[2]
		recipeOld := readRecipe(fileNameOld)

		fileNameNew := os.Args[4]
		recipeNew := readRecipe(fileNameNew)

		addRemoveCake(recipeNew, recipeOld, "ADDED")
		addRemoveCake(recipeOld, recipeNew, "REMOVED")
		addRemoveIngridient(recipeNew, recipeOld, "ADDED")
		addRemoveIngridient(recipeOld, recipeNew, "REMOVED")
		changeIngredient(recipeOld, recipeNew)
	} else {
		panic("Use '--old' & '--new' flags for passing path to Args.")
	}
}
