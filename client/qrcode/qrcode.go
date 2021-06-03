package qrcode

// Reference: https://github.com/dawndiy/qrcode-terminal/blob/master/qrcode-terminal.go

import (
	"bytes"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
	"golang.org/x/term"
)

var color = map[string]string{
	"BLACK":   "\033[48;5;0m  \033[0m",
	"RED":     "\033[48;5;1m  \033[0m",
	"GREEN":   "\033[48;5;2m  \033[0m",
	"YELLOW":  "\033[48;5;3m  \033[0m",
	"BLUE":    "\033[48;5;4m  \033[0m",
	"MAGENTA": "\033[48;5;5m  \033[0m",
	"CYAN":    "\033[48;5;6m  \033[0m",
	"WHITE":   "\033[48;5;7m  \033[0m",
}

var recoveryLevel = map[string]qrcode.RecoveryLevel{
	"LOW":     qrcode.Low,
	"MEDIUM":  qrcode.Medium,
	"HIGH":    qrcode.High,
	"HIGHEST": qrcode.Highest,
}

type Config struct {
	Foreground string
	Background string
	Level      string
	Align      string
}

func ShowQRcode(text string, config *Config) {
	foreground := color["BLACK"]
	background := color["WHITE"]
	level := "LOW"
	align := "left"

	if config != nil {
		if c := config.Foreground; c != "" {
			foreground = color[strings.ToUpper(c)]
		}
		if c := config.Background; c != "" {
			background = color[strings.ToUpper(c)]
		}
		if c := config.Level; c != "" {
			level = strings.ToUpper(c)
		}
		if c := config.Align; c != "" {
			align = strings.ToLower(c)
		}
	}

	qr, _ := qrcode.New(text, recoveryLevel[level])
	bitmap := qr.Bitmap()
	output := bytes.NewBuffer([]byte{})

	// Need to know the width
	screenCols, _, _ := term.GetSize(0)

	for ir, row := range bitmap {
		lr := len(row)

		if ir == 0 || ir == 1 || ir == 2 ||
			ir == lr-1 || ir == lr-2 || ir == lr-3 {
			continue
		}

		if align == "center" {
			for spaces := 0; spaces < (screenCols/2 - lr/2 - 2*(3*2)); spaces++ {
				output.WriteByte(' ')
			}
		}

		if align == "right" {
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
				output.WriteString(foreground)
			} else {
				output.WriteString(background)
			}
		}
		output.WriteByte('\n')
	}
	output.WriteTo(os.Stdout)
}
