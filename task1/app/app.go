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
	c          bool
	d          bool
	u          bool
	i          bool
	f          int
	s          int
	inputFile  string
	outputFile string
}

type Configure struct {
	cmd Cmd
}

func (conf Configure) checkI(str string) string {
	if conf.cmd.i {
		str = strings.ToLower(str)
	}
	return str
}

func (conf Configure) checkF(str string) string {
	arr := strings.Split(str, " ")
	if len(arr) > conf.cmd.f {
		str = strings.Join(arr[conf.cmd.f:], " ")
	}
	return str
}

func (conf Configure) checkS(str string) string {
	if len(str) > conf.cmd.s {
		str = (str)[conf.cmd.s:]
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
	count := make(map[string]int)

	j := 0

	for i, line := range lines {

		line = conf.checkI(line)
		line = conf.checkF(line)
		line = conf.checkS(line)

		if !alreadySeen[line] {
			alreadySeen[line] = true
			count[line] = j
			lines[j] = "1 " + (lines)[i]
			j++
		} else {
			oneLine := lines[count[line]]
			arr := strings.Split(oneLine, " ")
			val, _ := strconv.Atoi(arr[0])
			oneLine = strings.Join(arr[1:], " ")
			lines[count[line]] = strconv.Itoa(val+1) + " " + oneLine
		}
	}
	return lines[:j]
}

func (conf Configure) Duplicate(lines []string) []string {
	alreadySeen := make(map[string]bool)
	count := make(map[string]int)

	j := 0

	for i, line := range lines {

		line = conf.checkI(line)
		line = conf.checkF(line)
		line = conf.checkS(line)

		if !alreadySeen[line] {
			alreadySeen[line] = true
			count[line] = i
		} else {
			if count[line] != -1 {
				lines[j] = lines[count[line]]
				count[line] = -1
				j++
			}
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
	case conf.cmd.c:
		lines = conf.Counter(lines)
	case conf.cmd.d:
		lines = conf.Duplicate(lines)
	case conf.cmd.u:
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
	arr := make([]string, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}
	return arr
}

func (a Application) inFromScan(reader io.Reader) []string {
	return a.scanner(reader)
}

func (a Application) inFromFile(arr *[]string) error {
	file, err := os.Open(a.cmd.inputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	*arr = a.scanner(file)
	return nil
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
	return nil
}

func (a *Application) Parse() {
	flag.BoolVar(&a.cmd.help, "help", false, "help")
	flag.BoolVar(&a.cmd.c, "c", false, "подсчитать количество встречаний строки во входных данных. Вывести это число перед строкой отделив пробелом.")
	flag.BoolVar(&a.cmd.d, "d", false, "вывести только те строки, которые повторились во входных данных.")
	flag.BoolVar(&a.cmd.u, "u", false, "вывести только те строки, которые не повторились во входных данных.")
	flag.BoolVar(&a.cmd.i, "i", false, "не учитывать регистр букв.")
	flag.IntVar(&a.cmd.f, "f", 0, "не учитывать первые num_fields полей в строке. Полем в строке является непустой набор символов отделённый пробелом.")
	flag.IntVar(&a.cmd.s, "s", 0, "не учитывать первые num_chars символов в строке.")
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
	var arrBools = []bool{a.cmd.c, a.cmd.d, a.cmd.u}

	if len(a.filter(arrBools, func(b bool) bool { return b })) > 1 {
		return errors.New("don't use -c -d -u together")
	}
	return nil
}

func (a *Application) RunApp(reader io.Reader) string {
	a.Parse()

	if a.cmd.help {
		flag.PrintDefaults()
		return ""
	}

	if a.checkCDU() != nil {
		flag.PrintDefaults()
		return ""
	}

	arr := make([]string, 0)
	if a.cmd.inputFile == "" {
		arr = a.inFromScan(reader)
	} else {
		if a.inFromFile(&arr) != nil {
			return ""
		}
	}

	var u = Configure{}
	arr = u.Do(arr, a.cmd)

	if a.cmd.outputFile == "" {
		return a.outFromScan(arr)
	} else {
		err := a.outFromFile(arr)
		if err != nil {
			return "invalid to write to file"
		}
	}
	return ""
}
