package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Color is given its own type for safe function signatures
type Color string

// Color codes interpretted by the terminal
// NOTE: all codes must be of the same length or they will throw off the field alignment of tabwriter
const (
	Reset                   Color = "\x1b[0000m"
	Bright                        = "\x1b[0001m"
	BlackText                     = "\x1b[0030m"
	RedText                       = "\x1b[0031m"
	GreenText                     = "\x1b[0032m"
	YellowText                    = "\x1b[0033m"
	BlueText                      = "\x1b[0034m"
	MagentaText                   = "\x1b[0035m"
	CyanText                      = "\x1b[0036m"
	WhiteText                     = "\x1b[0037m"
	DefaultText                   = "\x1b[0039m"
	BrightRedText                 = "\x1b[1;31m"
	BrightGreenText               = "\x1b[1;32m"
	BrightYellowText              = "\x1b[1;33m"
	BrightBlueText                = "\x1b[1;34m"
	BrightMagentaText             = "\x1b[1;35m"
	BrightCyanText                = "\x1b[1;36m"
	BrightWhiteText               = "\x1b[1;37m"
	BlackBackground               = "\x1b[0040m"
	RedBackground                 = "\x1b[0041m"
	GreenBackground               = "\x1b[0042m"
	YellowBackground              = "\x1b[0043m"
	BlueBackground                = "\x1b[0044m"
	MagentaBackground             = "\x1b[0045m"
	CyanBackground                = "\x1b[0046m"
	WhiteBackground               = "\x1b[0047m"
	BrightBlackBackground         = "\x1b[0100m"
	BrightRedBackground           = "\x1b[0101m"
	BrightGreenBackground         = "\x1b[0102m"
	BrightYellowBackground        = "\x1b[0103m"
	BrightBlueBackground          = "\x1b[0104m"
	BrightMagentaBackground       = "\x1b[0105m"
	BrightCyanBackground          = "\x1b[0106m"
	BrightWhiteBackground         = "\x1b[0107m"
)

// Color implements the Stringer interface for interoperability with string
func (c *Color) String() string {
	return fmt.Sprintf("%v", c)
}

func Paint(color Color, value string) string {
	return fmt.Sprintf("%v %v%v", color, value, Reset)
}

func PaintRow(colors []Color, row []string) []string {
	paintedRow := make([]string, len(row))
	for i, v := range row {
		paintedRow[i] = Paint(colors[i], v)
	}
	return paintedRow
}

func PaintRowUniformly(color Color, row []string) []string {
	colors := make([]Color, len(row))
	for i, _ := range colors {
		colors[i] = color
	}
	return PaintRow(colors, row)
}

func AnonymizeRow(row []string) []string {
	anonRow := make([]string, len(row))
	for i, v := range row {
		anonRow[i] = strings.Repeat("-", len(v))
	}
	return anonRow
}

func PrintRow(writer io.Writer, line []string) {
	fmt.Fprintln(writer, "\t"+strings.Join(line, "\t"))
}

func ColorExamples() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Header
	header := []string{"Field1", "Field2", "Field3"}
	PrintRow(writer, PaintRowUniformly(DefaultText, header))
	PrintRow(writer, PaintRowUniformly(DefaultText, AnonymizeRow(header))) // header separator

	// Row A
	colors := []Color{RedText, CyanText, MagentaText}
	PrintRow(writer, PaintRow(colors, []string{"Aaron", "Adam", "Alice"}))

	// Row B
	colors = []Color{BrightYellowText, BrightGreenText, DefaultText}
	PrintRow(writer, PaintRow(colors, []string{"Bob", "Brooke", "Betty"}))

	// Actually print to screen
	writer.Flush()
}
