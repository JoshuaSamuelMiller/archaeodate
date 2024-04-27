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

// The basic date format, stores years before 1950
type date struct {
    name string
    datum int
}

type gregDate struct {
    name string
    datum int
    annoDomini bool
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

// Date handelling functions
func (Date gregDate) getGregDate() string {
    if Date.annoDomini == true {
	return "AD " + strconv.Itoa(Date.datum)
    } else {
	return strconv.Itoa(Date.datum) + " BC"
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

func printDate(d getDate) {
    fmt.Println(d.getGregDate())
}

func gregorianToDate(gregDate gregDate) date {
    if gregDate.annoDomini == true {
	var newDatum =  present - gregDate.datum
	return date{
	    name: gregDate.name,
	    datum: newDatum,
	}
    } else {
	var newDatum =  present + gregDate.datum - 1
	return date{
	    name: gregDate.name,
	    datum: newDatum,
	}
    }
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
	inputs[dateIn].Placeholder = "3800"
	inputs[dateIn].Focus()
	inputs[dateIn].CharLimit = 6
	inputs[dateIn].Width = 8
	inputs[dateIn].Prompt = ""

	inputs[systemIn] = textinput.New()
	inputs[systemIn].Placeholder = "BC"
	inputs[systemIn].Focus()
	inputs[systemIn].CharLimit = 2
	inputs[systemIn]. Width = 2
	inputs[systemIn].Prompt = ""

	var inputDate date = date{dateIn, getSystem(systemIn)}

	return model{
	    inputs: inputs,
	    result: "something",
	    focused: 0,
	    err: nil,
    }
}

// Update loop 
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs)) 

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
	headerStyle.Render("Press Ctrl-q to quit."),
    ) + "\n"

    //send to UI
    return s
}
