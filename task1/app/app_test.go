package app

import (
	"testing"
)

func TestCountMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          true,
		d:          false,
		u:          false,
		f:          0,
		s:          0,
		i:          false,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.", "I love music.", "I love music.", "\n", "I love music of Kartik.", "I love music of Kartik.", "Thanks."}
	expected := []string{"3 I love music.", "1 \n", "2 I love music of Kartik.", "1 Thanks."}
	u := Configure{}
	u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestDuplicatedMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          true,
		u:          false,
		f:          0,
		s:          0,
		i:          false,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.", "I love music.", "I love music.", "\n", "I love music of Kartik.", "I love music of Kartik.", "Thanks."}
	expected := []string{"I love music.", "I love music of Kartik."}
	u := Configure{}
	u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestUniqueMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          false,
		u:          true,
		f:          0,
		s:          0,
		i:          false,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.", "I love music.", "I love music.", "\n", "I love music of Kartik.", "I love music of Kartik.", "Thanks."}
	expected := []string{"\n", "Thanks."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestAnyCaseMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          false,
		u:          false,
		f:          0,
		s:          0,
		i:          true,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I LOVE MUSIC.",
		"I love music.",
		"I LoVe MuSiC.",
		"\n",
		"I love MuSIC of Kartik.",
		"I love music of kartik.",
		"Thanks."}
	expected := []string{"I LOVE MUSIC.", "\n", "I love MuSIC of Kartik.", "Thanks."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestBaseMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          false,
		u:          false,
		f:          0,
		s:          0,
		i:          false,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.", "I love music.", "I love music.", "\n", "I love music of Kartik.", "I love music of Kartik.", "Thanks."}
	expected := []string{"I love music.", "\n", "I love music of Kartik.", "Thanks."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestNumFieldsMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          false,
		u:          false,
		f:          1,
		s:          0,
		i:          false,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"We love music.",
		"I love music.",
		"They love music.",
		"\n",
		"I love music of Kartik.",
		"We love music of Kartik.",
		"Thanks."}
	expected := []string{"We love music.", "\n", "I love music of Kartik.", "Thanks."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestNumCharsMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          false,
		u:          false,
		f:          0,
		s:          1,
		i:          false,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.",
		"A love music.",
		"C love music.",
		"\n",
		"I love music of Kartik.",
		"We love music of Kartik.",
		"Thanks."}
	expected := []string{
		"I love music.",
		"\n",
		"I love music of Kartik.",
		"We love music of Kartik.",
		"Thanks.",
	}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestDuplicatedAndAnyCaseMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          true,
		u:          false,
		f:          0,
		s:          1,
		i:          true,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.", "I LoVe musiC.", "I Love Music.", "\n", "I love music of Kartik.", "I love music of Kartik.", "Thanks."}
	expected := []string{"I love music.", "I love music of Kartik."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestCountAndAnyCaseMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          true,
		d:          false,
		u:          false,
		f:          0,
		s:          1,
		i:          true,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.", "I LoVe musiC.", "I Love Music.", "\n", "I love music of Kartik.", "I love music of Kartik.", "Thanks."}
	expected := []string{"3 I love music.", "1 \n", "2 I love music of Kartik.", "1 Thanks."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestUniqueAndAnyCaseMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          false,
		u:          true,
		f:          0,
		s:          1,
		i:          true,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{"I love music.", "I LoVe musiC.", "I Love Music.", "\n", "I love music of Kartik.", "I love music of Kartik.", "Thanks."}
	expected := []string{"\n", "Thanks."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}

func TestNumCharFieldsMode(t *testing.T) {
	configure := Cmd{
		h:          false,
		c:          false,
		d:          false,
		u:          false,
		f:          1,
		s:          2,
		i:          false,
		inputFile:  "",
		outputFile: "",
	}
	current := []string{
		"We A love music.",
		"I B love music.",
		"They love music.",
		"\n",
		"I love music of Kartik.",
		"We love music of Kartik.",
		"Thanks."}
	expected := []string{"We A love music.", "They love music.", "\n", "I love music of Kartik.", "Thanks."}
	u := Configure{}
	current = u.Do(current, configure)
	for i := range expected {
		if expected[i] != current[i] {
			t.Error("Doesn't Equal")
		}
	}
}
