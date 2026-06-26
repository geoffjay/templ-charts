package grid

import (
	"math"
	"slices"
	"strconv"
)

// ComputeCellDimensionsArgs are the inputs to ComputeCellDimensions.
type ComputeCellDimensionsArgs struct {
	Width, Height float64
	Rows, Columns int
	Padding       float64
	Square        bool
}

// ComputeCellDimensions mirrors @nivo/grid computeCellDimensions: returns the
// [cellWidth, cellHeight] for a grid of `columns` x `rows` cells fitting into
// `width` x `height` with `padding` between cells. When `square` is true the
// smaller dimension is used for both so cells are square.
func ComputeCellDimensions(args ComputeCellDimensionsArgs) (cellWidth, cellHeight float64) {
	cellWidth = (args.Width - float64(args.Columns-1)*args.Padding) / float64(args.Columns)
	cellHeight = (args.Height - float64(args.Rows-1)*args.Padding) / float64(args.Rows)
	if !args.Square {
		return cellWidth, cellHeight
	}
	min := math.Min(cellWidth, cellHeight)
	return min, min
}

// GenerateGridArgs are the inputs to GenerateGrid (base form, no extender).
type GenerateGridArgs struct {
	Width, Height float64
	Columns, Rows int
	Padding       float64
	FillDirection GridFillDirection
	Square        bool
}

// GenerateGrid mirrors @nivo/grid generateGrid (base form): lays out
// `columns` x `rows` cells centered within `width` x `height`, reorders them
// by `fillDirection`, and assigns each cell its fill-order `index`.
//
// Defaults (matching upstream): padding=0, fillDirection="bottom", square=true.
func GenerateGrid(args GenerateGridArgs) GeneratedGrid[GridCell] {
	return GenerateGridWith[GridCell](GenerateGridWithArgs[GridCell]{
		Width:         args.Width,
		Height:        args.Height,
		Columns:       args.Columns,
		Rows:          args.Rows,
		Padding:       args.Padding,
		FillDirection: args.FillDirection,
		Square:        args.Square,
	})
}

// GenerateGridWithArgs extends GenerateGridArgs with an optional cell
// extender. When Extend is nil the base GridCells are returned as-is.
type GenerateGridWithArgs[C any] struct {
	Width, Height float64
	Columns, Rows int
	Padding       float64
	FillDirection GridFillDirection
	Square        bool
	Extend        CellExtender[C]
}

// GenerateGridWith mirrors @nivo/grid generateGrid<C>: the generic form that
// augments each base cell via `extend` (e.g. attaching a value/color for
// heatmap/waffle). When extend is nil the base GridCells are returned
// unchanged (C must be GridCell in that case).
func GenerateGridWith[C any](args GenerateGridWithArgs[C]) GeneratedGrid[C] {
	if args.FillDirection == "" {
		args.FillDirection = GridFillBottom
	}
	cellWidth, cellHeight := ComputeCellDimensions(ComputeCellDimensionsArgs{
		Width:   args.Width,
		Height:  args.Height,
		Rows:    args.Rows,
		Columns: args.Columns,
		Padding: args.Padding,
		Square:  args.Square,
	})

	origin := Vertex{
		(args.Width - (cellWidth*float64(args.Columns) + args.Padding*float64(args.Columns-1))) / 2,
		(args.Height - (cellHeight*float64(args.Rows) + args.Padding*float64(args.Rows-1))) / 2,
	}

	cells := make([]GridCell, 0, args.Rows*args.Columns)
	for row := 0; row < args.Rows; row++ {
		for column := 0; column < args.Columns; column++ {
			cells = append(cells, GridCell{
				Key:    strconv.Itoa(row) + "." + strconv.Itoa(column),
				Index:  0, // adjusted below by fillDirection
				Column: column,
				Row:    row,
				X:      float64(column) * cellWidth,
				Y:      float64(row) * cellHeight,
				Width:  cellWidth,
				Height: cellHeight,
			})
		}
	}

	switch args.FillDirection {
	case GridFillLeft:
		slices.SortFunc(cells, func(a, b GridCell) int {
			if a.Column != b.Column {
				return b.Column - a.Column
			}
			return b.Row - a.Row
		})
	case GridFillTop:
		slices.SortFunc(cells, func(a, b GridCell) int {
			if a.Row != b.Row {
				return b.Row - a.Row
			}
			return b.Column - a.Column
		})
	case GridFillRight:
		slices.SortFunc(cells, func(a, b GridCell) int {
			if a.Column != b.Column {
				return a.Column - b.Column
			}
			return a.Row - b.Row
		})
	default: // bottom — natural row-major order, no sort needed.
	}

	for i := range cells {
		cells[i].Index = i
	}

	out := make([]C, len(cells))
	for i, cell := range cells {
		if args.Extend != nil {
			out[i] = args.Extend(cell, origin)
		} else {
			// No extender: C is expected to be GridCell.
			if c, ok := any(cell).(C); ok {
				out[i] = c
			}
		}
	}

	return GeneratedGrid[C]{
		X:          origin[0],
		Y:          origin[1],
		CellWidth:  cellWidth,
		CellHeight: cellHeight,
		Cells:      out,
	}
}
