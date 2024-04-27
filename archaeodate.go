package main

import (
    "fmt"
    "strconv"
    "os"
//    "strings"
//    "log"

//    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
//    "github.com/charmbracelet/lipgloss"
)

//program consts
const (
    present int = 1950
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

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
	os.Exit(1)
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

type model struct {
    thing string
}

func initialModel() model {
    return model{
	thing: "text",
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
	switch msg.String() {
	case "ctrl+c", "q":
	    return m, tea.Quit
	}
    } 
    return m, nil
}

func (m model) View() string {
    return "something"
}
