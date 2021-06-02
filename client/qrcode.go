package client

import (
	"bytes"
	"os"

	"github.com/skip2/go-qrcode"
	"golang.org/x/term"
)

// Reference: https://github.com/dawndiy/qrcode-terminal/blob/master/qrcode-terminal.go

const (
	NormalBlack   = "\033[38;5;0m  \033[0m"
	NormalRed     = "\033[38;5;1m  \033[0m"
	NormalGreen   = "\033[38;5;2m  \033[0m"
	NormalYellow  = "\033[38;5;3m  \033[0m"
	NormalBlue    = "\033[38;5;4m  \033[0m"
	NormalMagenta = "\033[38;5;5m  \033[0m"
	NormalCyan    = "\033[38;5;6m  \033[0m"
	NormalWhite   = "\033[38;5;7m  \033[0m"

	BrightBlack   = "\033[48;5;0m  \033[0m"
	BrightRed     = "\033[48;5;1m  \033[0m"
	BrightGreen   = "\033[48;5;2m  \033[0m"
	BrightYellow  = "\033[48;5;3m  \033[0m"
	BrightBlue    = "\033[48;5;4m  \033[0m"
	BrightMagenta = "\033[48;5;5m  \033[0m"
	BrightCyan    = "\033[48;5;6m  \033[0m"
	BrightWhite   = "\033[48;5;7m  \033[0m"
)

var (
	frontColor      string
	backgroundColor string
	levelString     string
	codeJustify     string
)

func ShowQRcode(text string) {
	front := BrightBlack
	back := BrightWhite
	justify := "left"

	screenCols, _, _ := term.GetSize(0)

	qr, _ := qrcode.New(text, qrcode.Low)
	bitmap := qr.Bitmap()
	output := bytes.NewBuffer([]byte{})

	for ir, row := range bitmap {
		lr := len(row)

		if ir == 0 || ir == 1 || ir == 2 ||
			ir == lr-1 || ir == lr-2 || ir == lr-3 {
			continue
		}

		if justify == "center" {
			for spaces := 0; spaces < (screenCols/2 - lr/2 - 2*(3*2)); spaces++ {
				output.WriteByte(' ')
			}
		}

		if justify == "right" {
			for spaces := 0; spaces < (screenCols - 2*(lr-3*2)); spaces++ {
				output.WriteByte(' ')
			}
		}

		for ic, col := range row {
			lc := len(bitmap)
			if ic == 0 || ic == 1 || ic == 2 ||
				ic == lc-1 || ic == lc-2 || ic == lc-3 {
				continue
			}
			if col {
				output.WriteString(front)
			} else {
				output.WriteString(back)
			}
		}
		output.WriteByte('\n')
	}
	output.WriteTo(os.Stdout)
}
