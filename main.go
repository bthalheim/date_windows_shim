package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if err := run(os.Stdout, os.Stderr, os.Args[1:], time.Now); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(stdout, stderr io.Writer, args []string, now func() time.Time) error {
	utc := false
	format := ""

	for _, arg := range args {
		switch {
		case arg == "-u" || arg == "--utc":
			utc = true
		case arg == "-h" || arg == "--help":
			printUsage(stdout)
			return nil
		case strings.HasPrefix(arg, "+"):
			if format != "" {
				return errors.New("only one output format can be specified")
			}
			format = arg[1:]
		default:
			return fmt.Errorf("unsupported argument: %s", arg)
		}
	}

	current := now()
	if utc {
		current = current.UTC()
	}

	if format == "" {
		_, err := fmt.Fprintln(stdout, current.Format("Mon Jan _2 15:04:05 MST 2006"))
		return err
	}

	formatted, err := formatDate(current, format)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, formatted)
	return err
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: datew [-u|--utc] [+FORMAT]")
	fmt.Fprintln(w, "Print the current date and time using a small GNU date-compatible format subset.")
}

func formatDate(t time.Time, format string) (string, error) {
	var b strings.Builder

	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			b.WriteByte(format[i])
			continue
		}

		i++
		if i >= len(format) {
			return "", errors.New("dangling format specifier")
		}

		value, err := formatSpecifier(t, format[i])
		if err != nil {
			return "", err
		}
		b.WriteString(value)
	}

	return b.String(), nil
}

func formatSpecifier(t time.Time, spec byte) (string, error) {
	switch spec {
	case '%':
		return "%", nil
	case 'a':
		return t.Format("Mon"), nil
	case 'A':
		return t.Format("Monday"), nil
	case 'b', 'h':
		return t.Format("Jan"), nil
	case 'B':
		return t.Format("January"), nil
	case 'c':
		return t.Format("Mon Jan _2 15:04:05 2006"), nil
	case 'd':
		return t.Format("02"), nil
	case 'D':
		return t.Format("01/02/06"), nil
	case 'e':
		return fmt.Sprintf("%2d", t.Day()), nil
	case 'F':
		return t.Format("2006-01-02"), nil
	case 'H':
		return t.Format("15"), nil
	case 'I':
		return t.Format("03"), nil
	case 'j':
		return fmt.Sprintf("%03d", t.YearDay()), nil
	case 'm':
		return t.Format("01"), nil
	case 'M':
		return t.Format("04"), nil
	case 'p':
		return t.Format("PM"), nil
	case 'R':
		return t.Format("15:04"), nil
	case 'S':
		return t.Format("05"), nil
	case 'T':
		return t.Format("15:04:05"), nil
	case 'u':
		day := int(t.Weekday())
		if day == 0 {
			day = 7
		}
		return strconv.Itoa(day), nil
	case 'w':
		return strconv.Itoa(int(t.Weekday())), nil
	case 'x':
		return t.Format("01/02/06"), nil
	case 'X':
		return t.Format("15:04:05"), nil
	case 'y':
		return t.Format("06"), nil
	case 'Y':
		return t.Format("2006"), nil
	case 'z':
		return t.Format("-0700"), nil
	case 'Z':
		return t.Format("MST"), nil
	default:
		return "", fmt.Errorf("unsupported format specifier: %%%c", spec)
	}
}
