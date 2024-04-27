package main

import (
    "fmt"
    "strconv"
    "os"
//    "strings"
//    "log"

    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type (
    errMsg error
)

//program consts
const (
    present int = 1950
    
    lightGrey 	= lipgloss.Color("#949494")
    dullRed	= lipgloss.Color("#bc3838")
)

const (
    dateIn = iota
    systemIn
)

var (
    headerStyle	= lipgloss.NewStyle().Foreground(lightGrey)
    inputStyle	= lipgloss.NewStyle().Foreground(dullRed)

)

// The basic date struct containing a datum and datesystem
// 	System 	Type
//	0	BP
//	1	BC
//	2	AD
type date struct {
    datum int
    system int
}

type getDate interface {
    getGregDate() string
//    getBPDate() string
}


// Main function, starts the tea loop
func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("error, please check date format: %v", err)
	os.Exit(1)
    }
}

func (Date date) getGregDate() string {
    var gregorianDate = Date.datum - present
    if  gregorianDate > 0 {
	gregorianDate = gregorianDate + 1
	return strconv.Itoa(gregorianDate) + " BC"
    } else if gregorianDate == 0 { 
	gregorianDate = gregorianDate + 1
	return "AD " + strconv.Itoa(gregorianDate)
    } else {
	gregorianDate = -gregorianDate
	return "AD " + strconv.Itoa(gregorianDate)
    }
}

func getSystem(s string) int {
    switch {
    case s == "BP", s == "bp", s == "Bp":
	return 0
    case s == "BC", s == "bc", s == "Bc", s == "BCE", s == "bce":
	return 1
    case s == "AD", s == "ad", s == "Ad", s == "CE", s == "ce":
	return 2
    default:
	return 2
    }
}

func printSystem(s int) string {
    switch {
    case s == 0:
	return "BP"
    case s == 1:
	return "BC"
    case s == 2:
	return "AD"
    default:
	return "AD"
    }
}

func convertDate(d date) date {
    var newDatum int
    var newSystem int
    switch {
    case d.system == 0:
	newDatum = d.datum - present
	switch {
	case newDatum > 0:
	    newDatum = newDatum + 1
	    newSystem = 1
	case newDatum == 0:
	    newDatum = newDatum + 1
	    newSystem = 2
	case newDatum < 0:
	    newDatum = - newDatum
	    newSystem = 2
	}
    case d.system == 1:
	newDatum = d.datum + present - 1
	newSystem = 0
    case d.system == 2:
	newDatum = present - d.datum
	newSystem = 0
    }
    return date{newDatum, newSystem}
}

func getDatum(s string) int {
    var newDatum int
    var err error
    newDatum, err = strconv.Atoi(s)
    if err != nil{
	return 1950 
    }
    return newDatum
}

// display functions
func (m *model) nextInput() {
    m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *model) prevInput() {
    m.focused--
    if m.focused < 0 {
	m.focused = len(m.inputs) -1
    }
}

func (m model) Init() tea.Cmd {
    return textinput.Blink 
}

// State model
type model struct {
    result string
    focused int
    inputs []textinput.Model
    err error
}

// intial state of model
func initialModel() model {
    var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[dateIn] = textinput.New()
	inputs[dateIn].Placeholder = "1950"
	inputs[dateIn].Focus()
	inputs[dateIn].CharLimit = 6
	inputs[dateIn].Width = 8
	inputs[dateIn].Prompt = ""

	inputs[systemIn] = textinput.New()
	inputs[systemIn].Placeholder = "AD"
	inputs[systemIn].Focus()
	inputs[systemIn].CharLimit = 2
	inputs[systemIn]. Width = 2
	inputs[systemIn].Prompt = ""

	

	return model{
	    inputs: inputs,
	    result: "0 BP",
	    focused: 0,
	    err: nil,
    }
}

// Update loop 
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs)) 
    var inputDate date
    var outputDate date

    inputDate = date{getDatum(m.inputs[dateIn].Value()), getSystem(m.inputs[systemIn].Value())}
    outputDate = convertDate(inputDate)
    m.result = strconv.Itoa(outputDate.datum) + " " + printSystem(outputDate.system)

    switch msg := msg.(type) {
    case tea.KeyMsg:
	switch msg.Type {
	case tea.KeyEnter:
	    m.nextInput()
	case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlQ:
	    return m, tea.Quit
	case tea.KeyShiftTab, tea.KeyCtrlP:
	    m.prevInput()
	case tea.KeyTab, tea.KeyCtrlN:
	    m.nextInput()
	} 
	for i := range m.inputs {
	    m.inputs[i].Blur()
	}
	m.inputs[m.focused].Focus()

    case errMsg:
	m.err = msg
	return m, nil
    }
    for i := range m.inputs {
	m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
    }
    return m, tea.Batch(cmds...)
}

// This 'renders' renders the UI
func (m model) View() string {
    s := fmt.Sprintf(
`%s

    %s
    %s %s

    %s %s

%s
`,
	headerStyle.Render("archaeodate"),
	inputStyle.Render("Enter Date"),
	m.inputs[dateIn].View(),
	m.inputs[systemIn].View(),
	inputStyle.Render("Result:"),
	m.result,
	headerStyle.Render("Press Esc to quit."),
    ) + "\n"

    //send to UI
    return s
}
