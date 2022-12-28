package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/valyala/fasttemplate"
)

var (
	// colours is a list of HEX values.
	// Use "colorize" extension for VSC to view colours
	// within the IDE.
	colours = []string{
		"#EF6F6C",
		"#465775",
		"#56E39F",
		"#59C9A5",
		"#D9CAB3",
		"#BC8034",
		"#90323D",
		"#F6BD60",
		"#F5CAC3",
		"#247BA0",
		"#A7ABDD",
		"#70C1B3",
		"#000000",
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := randomiseBadgeColours(colours); err != nil {
		panic(err)
	}
}

func randomiseBadgeColours(palette []string) error {
	f, err := os.OpenFile("README.md", os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		f.Close()
		return err
	}

	f.Close()

	splits := []string{
		"https://img.shields.io/badge/",
		"?style=for-the-badge",
	}

	t, err := fasttemplate.NewTemplate(
		string(contents),
		splits[0],
		splits[1],
	)
	if err != nil {
		return err
	}

	fmt.Print("Replacing tag colours: ")
	s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		fmt.Print(".")
		sheildOptions := strings.Split(tag, "-")

		// write shield title
		n, err := w.Write([]byte(splits[0] + sheildOptions[0] + "-"))
		if err != nil {
			return n, err
		}

		// write shield hex colour code
		nx, err := w.Write([]byte(palette[rand.Intn(len(palette))][1:] + splits[1]))
		return n + nx, err
	})

	fmt.Println("\nColours replaced")

	return ioutil.WriteFile("README.staging.md", []byte(s), 0666)

}
