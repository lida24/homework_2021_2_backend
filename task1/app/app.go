package app

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"os"
	"strconv"
	"strings"
)

type Cmd struct {
	help       bool
	count      bool
	duplicates bool
	unique     bool
	ignore     bool
	fields     int
	symbol     int
	inputFile  string
	outputFile string
}

type Configure struct {
	cmd Cmd
}

func (conf Configure) checkI(str string) string {
	if conf.cmd.ignore {
		str = strings.ToLower(str)
	}
	return str
}

func (conf Configure) checkF(str string) string {
	arr := strings.Split(str, " ")
	if len(arr) > conf.cmd.fields {
		str = strings.Join(arr[conf.cmd.fields:], " ")
	}
	return str
}

func (conf Configure) checkS(str string) string {
	if len(str) > conf.cmd.symbol {
		str = (str)[conf.cmd.symbol:]
	}
	return str
}

func (conf Configure) removeDuplicates(lines []string) []string {
	alreadySeen := make(map[string]bool)

	j := 0

	for i, line := range lines {

		line = conf.checkI(line)
		line = conf.checkF(line)
		line = conf.checkS(line)

		if !alreadySeen[line] {
			alreadySeen[line] = true
			lines[j] = lines[i]
			j++
		}
	}
	return lines[:j]
}

func (conf Configure) Counter(lines []string) []string {
	alreadySeen := make(map[string]bool)

	j := 0
	count := 0

	for i, line := range lines {

		line = conf.checkI(line)
		line = conf.checkF(line)
		line = conf.checkS(line)

		if !alreadySeen[line] {
			alreadySeen[line] = true
			count = j
			lines[j] = "1 " + (lines)[i]
			j++
		} else {
			oneLine := lines[count]
			arr := strings.Split(oneLine, " ")
			val, _ := strconv.Atoi(arr[0])
			oneLine = strings.Join(arr[1:], " ")
			lines[count] = strconv.Itoa(val+1) + " " + oneLine
		}
	}
	return lines[:j]
}

func (conf Configure) Duplicate(lines []string) []string {
	alreadySeen := make(map[string]bool)
	count := 0

	j := 0

	for i, line := range lines {

		line = conf.checkI(line)
		line = conf.checkF(line)
		line = conf.checkS(line)

		if !alreadySeen[line] {
			alreadySeen[line] = true
			count = i
		} else if count != -1 {
			lines[j] = lines[count]
			count = -1
			j++
		}
	}

	return lines[:j]
}

func (conf Configure) Unique(lines []string) []string {
	lines = conf.Counter(lines)

	j := 0

	for _, line := range lines {
		arr := strings.Split(line, " ")
		val, _ := strconv.Atoi(arr[0])
		if val == 1 {
			lines[j] = strings.Join(arr[1:], " ")
			j++
		}
	}
	return lines[:j]
}

func (conf Configure) Do(lines []string, cmd Cmd) []string {
	conf.cmd = cmd
	switch true {
	case conf.cmd.count:
		lines = conf.Counter(lines)
	case conf.cmd.duplicates:
		lines = conf.Duplicate(lines)
	case conf.cmd.unique:
		lines = conf.Unique(lines)
	default:
		lines = conf.removeDuplicates(lines)
	}
	return lines
}

type Application struct {
	cmd Cmd
}

func (a Application) scanner(r io.Reader) []string {
	var arr []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}
	return arr
}

func (a Application) inFromScan(reader io.Reader) []string {
	return a.scanner(reader)
}

func (a Application) inFromFile() ([]string, error) {
	var arr []string
	file, err := os.Open(a.cmd.inputFile)
	if err != nil {
		return arr, err
	}
	defer file.Close()
	arr = a.scanner(file)
	return arr, nil
}

func (a Application) outFromScan(arr []string) string {
	var str strings.Builder
	for _, x := range arr {
		str.WriteString(x)
	}
	return str.String()
}

func (a Application) outFromFile(arr []string) error {
	fileOut, err := os.Create(a.cmd.outputFile)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	for _, x := range arr {
		_, err = fileOut.WriteString(x + "\n")
		if err != nil {
			return err
		}
	}
	return fileOut.Close()
}

func (a *Application) Parse() {
	flag.BoolVar(&a.cmd.help, "help", false, "help")
	flag.BoolVar(&a.cmd.count, "c", false, "подсчитать количество встречаний строки во входных данных. Вывести это число перед строкой отделив пробелом.")
	flag.BoolVar(&a.cmd.duplicates, "d", false, "вывести только те строки, которые повторились во входных данных.")
	flag.BoolVar(&a.cmd.unique, "u", false, "вывести только те строки, которые не повторились во входных данных.")
	flag.BoolVar(&a.cmd.ignore, "i", false, "не учитывать регистр букв.")
	flag.IntVar(&a.cmd.fields, "f", 0, "не учитывать первые num_fields полей в строке. Полем в строке является непустой набор символов отделённый пробелом.")
	flag.IntVar(&a.cmd.symbol, "s", 0, "не учитывать первые num_chars символов в строке.")
	flag.Parse()

	var names = []string{a.cmd.inputFile, a.cmd.outputFile}
	for i := 0; i < len(flag.Args()); i++ {
		names = append(names, flag.Args()[i])
	}
}

func (a Application) filter(arrayBool []bool, f func(bool) bool) []bool {
	newArray := make([]bool, 0)
	for _, value := range arrayBool {
		if f(value) {
			newArray = append(newArray, value)
		}
	}
	return newArray
}

func (a Application) checkCDU() error {
	var arrBools = []bool{a.cmd.count, a.cmd.duplicates, a.cmd.unique}

	if len(a.filter(arrBools, func(b bool) bool { return b })) > 1 {
		return errors.New("don't use -c -d -u together")
	}
	return nil
}

func (a *Application) RunApp(reader io.Reader) (string, error) {
	a.Parse()

	if a.cmd.help {
		flag.PrintDefaults()
		return "", nil
	}

	if a.checkCDU() != nil {
		flag.PrintDefaults()
		return "", nil
	}

	var arr []string
	if a.cmd.inputFile == "" {
		arr = a.inFromScan(reader)
	} else {
		arr, _ = a.inFromFile()
	}

	var u = Configure{}
	arr = u.Do(arr, a.cmd)

	if a.cmd.outputFile == "" {
		return a.outFromScan(arr), nil
	} else {
		err := a.outFromFile(arr)
		if err != nil {
			return "", errors.New("invalid to write to file")
		}
	}
	return "", nil
}
