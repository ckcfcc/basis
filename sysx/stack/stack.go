package stack

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const maxStackSize = 32

type Frame struct {
	File string
	Line int
	Name string
}

func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Name)
}

type Stack []Frame

func (s Stack) String() string {
	var b bytes.Buffer
	writeStack(&b, s)
	return b.String()
}

type Multi struct {
	stacks []Stack
}

func (m *Multi) Stacks() []Stack {
	return m.stacks
}

func (m *Multi) Add(s Stack) {
	m.stacks = append(m.stacks, s)
}

func (m *Multi) AddCallers(skip int) {
	m.Add(Callers(skip + 1))
}

func (m *Multi) String() string {
	var b bytes.Buffer
	for i, s := range m.stacks {
		if i != 0 {
			fmt.Fprintf(&b, "\n(Stack %d)\n", i+1)
		}
		writeStack(&b, s)
	}
	return b.String()
}

func (m *Multi) Copy() *Multi {
	m2 := &Multi{
		stacks: make([]Stack, len(m.stacks)),
	}
	copy(m2.stacks, m.stacks)
	return m2
}

func Caller(skip int) Frame {
	pc, file, line, _ := runtime.Caller(skip + 1)
	fun := runtime.FuncForPC(pc)
	return Frame{
		File: StripGOPATH(file),
		Line: line,
		Name: StripPackage(fun.Name()),
	}
}

func Callers(skip int) Stack {
	pcs := make([]uintptr, maxStackSize)
	num := runtime.Callers(skip+2, pcs)
	stack := make(Stack, num)
	for i, pc := range pcs[:num] {
		fun := runtime.FuncForPC(pc)
		file, line := fun.FileLine(pc - 1)
		stack[i].File = StripGOPATH(file)
		stack[i].Line = line
		stack[i].Name = StripPackage(fun.Name())
	}
	return stack
}

func CallersMulti(skip int) *Multi {
	m := new(Multi)
	m.AddCallers(skip + 1)
	return m
}

func writeStack(b *bytes.Buffer, s Stack) {
	var width int
	for _, f := range s {
		if l := len(f.File) + numDigits(f.Line) + 1; l > width {
			width = l
		}
	}
	last := len(s) - 1
	b.WriteString("\n\t\t\t {\n")
	for i, f := range s {
		b.WriteString("\t\t\t\t")
		b.WriteString(f.File)
		b.WriteRune(rune(':'))
		n, _ := fmt.Fprintf(b, "%d", f.Line)
		for i := width - len(f.File) - n; i != 0; i-- {
			b.WriteRune(rune(' '))
		}
		b.WriteString(f.Name)
		if i != last {
			b.WriteRune(rune('\n'))
		}
	}
	b.WriteString("\n\t\t\t }\n")
}

func numDigits(i int) int {
	var n int
	for {
		n++
		i = i / 10
		if i == 0 {
			return n
		}
	}
}

var (
	gopath  string
	gopaths []string
)

func init() {
	if gopath == "" {
		gopath = os.Getenv("GOPATH")
	}
	SetGOPATH(gopath)
}

func StripGOPATH(f string) string {
	for _, p := range gopaths {
		if strings.HasPrefix(f, p) {
			return f[len(p):]
		}
	}
	return f
}

func SetGOPATH(gp string) {
	gopath = gp
	gopaths = nil

	for _, p := range strings.Split(gopath, ":") {
		if p != "" {
			gopaths = append(gopaths, filepath.Join(p, "src")+"/")
		}
	}

	gopaths = append(gopaths, filepath.Join(runtime.GOROOT(), "src", "pkg")+"/")
}

func StripPackage(n string) string {
	slashI := strings.LastIndex(n, "/")
	if slashI == -1 {
		slashI = 0
	}
	dotI := strings.Index(n[slashI:], ".")
	if dotI == -1 {
		return n
	}
	return n[slashI+dotI+1:]
}
