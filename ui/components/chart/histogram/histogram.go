package histogram

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/context"
	"github.com/danielcmessias/sawsy/utils"
)

type Model struct {
	data []float64

	width  int
	height int

	// If true, the histogram will not sample and instead draw empty columns or truncate the data to exactly fit the width.
	TrueSize bool

	YMin          *float64
	YMax          *float64
	yAxisMaxWidth int

	viewContents string

	pane.Pane
}

type HistogramSpec struct {
	pane.BaseSpec
}

type NewDataMsg struct {
	Page          string
	PaneId        int
	GalleryPaneId int // Refers to the pane within the gallery
	Data          []float64
}

func (s HistogramSpec) NewFromSpec(ctx *context.ProgramContext, spec pane.PaneSpec) pane.Pane {
	histogramSpec, ok := spec.(HistogramSpec)
	if !ok {
		log.Fatal("invalid spec type, expected HistogramSpec")
	}
	return New(ctx, histogramSpec)
}

func New(ctx *context.ProgramContext, spec HistogramSpec) *Model {
	var data []float64
	for i := 0; i < 40; i++ {
		// data = append(data, rand.Float64()*100)
	}

	// var data []float64
	return &Model{
		Pane: pane.New(spec.BaseSpec),

		data:     data,
		TrueSize: false,

		YMin: utils.Float64Ptr(0),
	}
}

func (m *Model) View() string {
	return m.viewContents
}

func (m *Model) updateView() {
	if m.data == nil || len(m.data) == 0 {
		m.viewContents = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Width(m.width).
			Height(m.height).
			Render("No data")
		return
	}
	axisWidth := m.yAxisMaxWidth + 1
	plotWidth := m.width - axisWidth
	var nCols int
	var unitsPerCol float64

	if m.TrueSize {
		nCols = utils.Min(plotWidth, len(m.data))
		unitsPerCol = 1
	} else {
		nCols = plotWidth
		unitsPerCol = float64(len(m.data)) / float64(nCols)
	}

	var yMin, yMax float64
	if m.YMin == nil {
		yMin = Min(m.data)
	} else {
		yMin = *m.YMin
	}
	if m.YMax == nil {
		yMax = Max(m.data)
	} else {
		yMax = *m.YMax
	}
	yRange := yMax - yMin
	unitsPerRow := yRange / float64(m.height-1)

	cells := make([][]string, m.height)
	for i := range cells {
		cells[i] = make([]string, nCols)
		for j := range cells[i] {
			cells[i][j] = " "
		}
	}

	for j := 0; j < nCols; j++ {
		sample := m.data[int(math.Floor(float64(j)*unitsPerCol))]

		nRows := (sample - yMin) / unitsPerRow
		if math.IsNaN(nRows) {
			nRows = 0
		}
		fullRows := int(nRows)
		remainder := nRows - float64(fullRows)

		for i := 0; i < fullRows; i++ {
			cells[i][j] = barsStyle.Render(BLOCKS[8])
		}
		tmp := int(math.Round(remainder * 8))
		cells[fullRows][j] = barsStyle.Render(BLOCKS[tmp])
	}

	rows := make([]string, m.height)
	for i, row := range cells {
		// Add one for the right padding
		axisLabel := strings.Repeat(" ", m.yAxisMaxWidth+1)
		// Only label every other line
		if i%2 == 0 {
			val := yMin + float64(i)*unitsPerRow

			// Kill decimals for now, max life easy!
			formatter := "%" + fmt.Sprint(m.yAxisMaxWidth) + ".0f "
			axisLabel = fmt.Sprintf(formatter, val)
		}
		rows[len(rows)-i-1] = axisStyle.Render(axisLabel) + lipgloss.JoinHorizontal(lipgloss.Left, row...)
	}
	m.viewContents = lipgloss.JoinVertical(lipgloss.Top, rows...)
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.updateView()
}

func (m *Model) SetData(datapoints []float64) {
	m.data = datapoints
	for _, d := range datapoints {
		m.yAxisMaxWidth = utils.Max(m.yAxisMaxWidth, lipgloss.Width(axisStyle.Render(fmt.Sprintf("%.0f", d))))
	}
	m.updateView()
}
