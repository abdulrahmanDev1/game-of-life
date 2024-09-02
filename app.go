package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

func NewGridApp(gridWidth, gridHeight int) (*GridApp, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}

	grid := NewGrid(gridWidth, gridHeight)

	return &GridApp{
		screen:     screen,
		grid:       grid,
		gridWidth:  gridWidth,
		gridHeight: gridHeight,
		paused:     true,
		changed:    true,
	}, nil
}

func cellNeighborsCount(grid [][]bool, x, y int) int {
	count := 0
	height, width := len(grid), len(grid[0])
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := (x+dx+width)%width, (y+dy+height)%height
			if grid[ny][nx] {
				count++
			}
		}
	}
	return count
}

func (app *GridApp) Run() error {
	app.screen.EnableMouse()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	lastUpdate := time.Now()

	for {
		select {
		case <-ticker.C:
			if !app.paused {
				now := time.Now()
				if now.Sub(lastUpdate) >= 200*time.Millisecond {
					app.updateGrid()
					lastUpdate = now
				}
			}
			if app.changed {
				app.draw()
				app.changed = false
			}
		default:
			// Check if there are pending events before polling
			if app.screen.HasPendingEvent() {
				ev := app.screen.PollEvent()
				switch ev := ev.(type) {
				case *tcell.EventResize:
					app.screen.Sync()
					app.changed = true
				case *tcell.EventKey:
					// Get the current screen size
					screenWidth, screenHeight := app.screen.Size()

					// Calculate maximum grid width and height
					maxGridWidth := (screenWidth - 2) / 3
					maxGridHeight := (screenHeight - 2) / 2

					switch ev.Key() {
					case tcell.KeyEscape:
						return nil
					case tcell.KeyRune:
						switch ev.Rune() {
						case 'q':
							return nil
						case 'c':
							app.ClearGrid()
							app.changed = true
						case ' ':
							app.paused = !app.paused
							app.changed = true
						}
					case tcell.KeyUp:
						if app.gridHeight > 1 {
							app.gridHeight = max(1, app.gridHeight-1)
							app.grid.Resize(app.gridWidth, app.gridHeight)
							app.changed = true
						}
					case tcell.KeyDown:
						if app.gridHeight < maxGridHeight {
							app.gridHeight++
							app.grid.Resize(app.gridWidth, app.gridHeight)
							app.changed = true
						}
					case tcell.KeyLeft:
						if app.gridWidth > 1 {
							app.gridWidth = max(1, app.gridWidth-1)
							app.grid.Resize(app.gridWidth, app.gridHeight)
							app.changed = true
						}
					case tcell.KeyRight:
						if app.gridWidth < maxGridWidth {
							app.gridWidth++
							app.grid.Resize(app.gridWidth, app.gridHeight)
							app.changed = true
						}
					}
				case *tcell.EventMouse:
					x, y := ev.Position()
					button := ev.Buttons()
					if button == tcell.Button1 {
						screenWidth, screenHeight := app.screen.Size()
						startX := (screenWidth - app.gridWidth*3) / 2
						startY := (screenHeight - app.gridHeight*2) / 2
						cellX, cellY := (x-startX)/3, (y-startY)/2
						if cellX >= 0 && cellX < app.gridWidth && cellY >= 0 && cellY < app.gridHeight {
							app.grid.ToggleCell(cellX, cellY)
							app.changed = true
						}
					}
				}
			}
		}
	}
}

func (app *GridApp) updateGrid() {
	newCells := make([][]bool, app.gridHeight)
	for i := range newCells {
		newCells[i] = make([]bool, app.gridWidth)
	}

	changed := false
	for y := 0; y < app.gridHeight; y++ {
		for x := 0; x < app.gridWidth; x++ {
			count := cellNeighborsCount(app.grid.Cells, x, y)
			if app.grid.Cells[y][x] {
				newCells[y][x] = count == 2 || count == 3
			} else {
				newCells[y][x] = count == 3
			}
			if newCells[y][x] != app.grid.Cells[y][x] {
				changed = true
			}
		}
	}

	if changed {
		app.grid.Cells = newCells
		app.changed = true
	}
}

func (app *GridApp) drawTopBorder(x, screenX, screenY int) {
	gridStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	if x == 0 {
		// Top-left corner
		app.screen.SetContent(screenX, screenY, '╔', nil, gridStyle)
	} else if x == app.gridWidth {
		// Top-right corner
		app.screen.SetContent(screenX, screenY, '╗', nil, gridStyle)
	} else {
		// Top border
		app.screen.SetContent(screenX, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX+1, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX-1, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX-2, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX+2, screenY, '═', nil, gridStyle)
	}
}

func (app *GridApp) drawBottomBorder(x, screenX, screenY int) {
	gridStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	if x == 0 {
		// Bottom-left corner
		app.screen.SetContent(screenX, screenY, '╚', nil, gridStyle)
	} else if x == app.gridWidth {
		// Bottom-right corner
		app.screen.SetContent(screenX, screenY, '╝', nil, gridStyle)
	} else {
		// Bottom border
		app.screen.SetContent(screenX, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX+1, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX-1, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX+2, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX-2, screenY, '═', nil, gridStyle)
	}
}

func (app *GridApp) drawBorders(x, y, screenX, screenY int) {
	if y == 0 {
		// Draw the top border
		app.drawTopBorder(x, screenX, screenY)
	} else if y == app.gridHeight {
		// Draw the bottom border
		app.drawBottomBorder(x, screenX, screenY)
	} else if x == 0 || x == app.gridWidth {
		// Left and right borders
		gridStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		app.screen.SetContent(screenX, screenY+1, '║', nil, gridStyle)
		app.screen.SetContent(screenX, screenY, '║', nil, gridStyle)
		app.screen.SetContent(screenX, screenY-1, '║', nil, gridStyle)
	} else {
		// Draw the grid intersections and grid lines
		gridStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		app.screen.SetContent(screenX, screenY, '╬', nil, gridStyle)
		app.screen.SetContent(screenX+1, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX-1, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX-2, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX+2, screenY, '═', nil, gridStyle)
		app.screen.SetContent(screenX, screenY+1, '║', nil, gridStyle)
		app.screen.SetContent(screenX, screenY-1, '║', nil, gridStyle)
	}
}

func (app *GridApp) draw() {
	app.screen.Clear()

	screenWidth, screenHeight := app.screen.Size()
	startX := (screenWidth - app.gridWidth*3) / 2
	startY := (screenHeight - app.gridHeight*2) / 2

	activeStyle := tcell.StyleDefault.Background(tcell.ColorWhite)
	inactiveStyle := tcell.StyleDefault.Background(tcell.ColorBlack)

	for y := 0; y <= app.gridHeight; y++ {
		for x := 0; x <= app.gridWidth; x++ {
			screenX, screenY := startX+x*3, startY+y*2

			app.drawBorders(x, y, screenX, screenY)

			if y < app.gridHeight && x < app.gridWidth {
				cellStyle := inactiveStyle
				if app.grid.Cells[y][x] {
					cellStyle = activeStyle
				}
				app.screen.SetContent(screenX+1, screenY+1, ' ', nil, cellStyle)
				app.screen.SetContent(screenX+2, screenY+1, ' ', nil, cellStyle)
			}
		}
	}

	// Define some styles
	var (
		defaultStyle = tcell.StyleDefault
		greenStyle   = defaultStyle.Foreground(tcell.ColorGreen)
		redStyle     = defaultStyle.Foreground(tcell.ColorRed)
		blueStyle    = defaultStyle.Foreground(tcell.ColorBlue)
		grayStyle    = defaultStyle.Foreground(tcell.ColorGray)
	)

	controls := []struct {
		text  string
		style tcell.Style
	}{
		{"Controls:", defaultStyle},
		{"Arrow keys: Resize grid", blueStyle},
		{"Left click: Toggle cell", blueStyle},
		{"Space: " + func() string {
			if app.paused {
				return "Resume ▷"
			} else {
				return "Pause ||"
			}
		}(), func() tcell.Style {
			if app.paused {
				return greenStyle
			} else {
				return redStyle
			}
		}()},
		{"Clear grid: C", redStyle},
		{"Quit: Q or Esc", grayStyle},
	}

	// Then, when drawing the controls, use the style associated with each control
	for i, control := range controls {
		for j, ch := range control.text {
			app.screen.SetContent(1+j, screenHeight-len(controls)+i, ch, nil, control.style)
		}
	}

	app.screen.Show()
}

func (app *GridApp) Cleanup() {
	app.screen.Fini()
}

func (app *GridApp) ClearGrid() {
	for y := 0; y < app.gridHeight; y++ {
		for x := 0; x < app.gridWidth; x++ {
			app.grid.Cells[y][x] = false
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
