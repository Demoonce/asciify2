package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"strconv"

	// "strings"

	// "fmt"
	// "time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	term "github.com/gookit/color"
)

// `.-~_,:!^+<>=;*/()?{}[]7123#4$5%980@
// 712345908
// `.-,:!^+<=*/(?{[#$%@

func ToText(filename string, width int, chars []rune, color bool, use_gaussian bool) (string, error) {
	data, err := imgio.Open(filename)
	if err != nil {
		return "", err
	}
	w := data.Bounds().Dx()
	h := data.Bounds().Dy()
	scale := float64(w) / float64(h)
	if width == 0 {
		width = w
	}
	var res_height = int(float64(width) / (scale + 1))
	if use_gaussian {
		data = transform.Resize(data, width, res_height, transform.Gaussian)
	} else {
		data = transform.Resize(data, width, res_height, transform.Linear)
	}

	if data == nil {
		log.Fatalln("Can't resize image")
	}
	var result string
	var base int
	var string_chars = make([]string, len(chars))
	for a := range chars {
		string_chars[a] = string(chars[a])
	}
	if 256%len(chars) == 0 {
		base = 256 / len(chars)
	} else {
		base = 256/len(chars) + 1
	}
	for a := 0; a < res_height; a++ {
		for b := 0; b < width; b++ {
			r, g, b, _ := data.At(b, a).RGBA()
			gray := int((r + g + b) / 768)
			red, green, blue := uint8(r), uint8(g), uint8(b)
			ch := string_chars[gray/base]
			if color {
				result += term.RGB(red, green, blue).Sprint(ch)
			} else {
				result += ch
			}
		}
		result += "\n\n"
	}
	return result, nil
}

type Params struct {
	width int
	chars []rune
	color bool
	filename string
}

func main() {
	a := app.NewWithID("asciify")
	w := a.NewWindow("Asciify")
	params := Params{0, []rune("`.-,:!^+<=*/(?{[#$%@"), true, ""}
	errors := widget.NewLabel("")
	width := widget.NewEntry()
	width.SetPlaceHolder("Enter width")
	width.OnChanged = func(s string) {
		_, err := strconv.Atoi(s)
		if err != nil {
			errors.SetText("Please enter a numeric value")
		} else {
			errors.SetText("")
		}
	}
	


	var result string

	chars := widget.NewEntry()
	chars.SetPlaceHolder("Enter chars(default ones are used if empty)")
	file_selected := widget.NewLabel("No file selected")
	file_open := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			return
		}
		if err != nil {
			log.Fatalln(err)
		}
		if err != nil {
			return
		}
		if reader.URI() == nil {
			return
		}
		params.filename = reader.URI().Path()
		file_selected.SetText(params.filename)

	}, w)

	
	
	select_button := widget.NewButton("Select an image", func() {

		file_open.Show()
	})
	done_label := widget.NewLabel("")
	done_label.TextStyle.Bold = true
	done_label.Alignment = fyne.TextAlignTrailing
	generate_button := widget.NewButton("Generate", func() {
		errors.SetText("")
		width_value, _ := strconv.Atoi(width.Text)
		c := chars.Text
		params.width = width_value
		if c != "" {
			params.chars = []rune(c)
		}
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			}
			if err != nil {
				errors.SetText(err.Error())
				return
			}
			result, err = ToText(params.filename, params.width, params.chars, false, false)
			if err != nil {
				errors.SetText(err.Error())
				return 
			}
			fmt.Println(result)
			writer.Write([]byte(result))
			done_label.SetText("Done!")

		}, w)
	})

	cont := container.NewVBox(
		width, 
		chars, 
		container.NewHBox(file_selected, select_button, done_label), 
		generate_button, 
		errors)
	w.Resize(fyne.NewSize(600, 600))
	w.SetContent(container.NewStack(cont))
	w.ShowAndRun()
}