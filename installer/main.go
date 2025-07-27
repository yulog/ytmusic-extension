package main

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"

	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
)

type modelKey int

const (
	modelKeyModel modelKey = iota
)

type Root struct {
	guigui.DefaultWidget

	background basicwidget.Background
	sidebar    Sidebar
	steps      Steps
	navbar     Navbar

	model Model

	faceSourceEntries []basicwidget.FaceSourceEntry
}

func (r *Root) updateFontFaceSources(context *guigui.Context) {
	// r.locales = slices.Delete(r.locales, 0, len(r.locales))
	// r.locales = context.AppendLocales(r.locales)

	// r.faceSourceEntries = slices.Delete(r.faceSourceEntries, 0, len(r.faceSourceEntries))
	r.faceSourceEntries = cjkfont.AppendRecommendedFaceSourceEntries(r.faceSourceEntries, []language.Tag{language.Und})
	basicwidget.SetFaceSources(r.faceSourceEntries)
}

func (r *Root) Model(key any) any {
	switch key {
	case modelKeyModel:
		return &r.model
	default:
		return nil
	}
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.background)
	appender.AppendChildWidget(&r.sidebar)
	appender.AppendChildWidget(&r.steps)
	appender.AppendChildWidget(&r.navbar)
}

func (r *Root) Build(context *guigui.Context) error {
	r.updateFontFaceSources(context)

	context.SetBounds(&r.background, context.Bounds(r), r)

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	context.SetBounds(&r.sidebar, gl.CellBounds(0, 0), r)
	{
		u := basicwidget.UnitSize(context)
		gl := layout.GridLayout{
			Bounds: gl.CellBounds(1, 0).Inset(u / 2),
			Heights: []layout.Size{
				layout.FlexibleSize(1),
				layout.FixedSize(u + u/2),
			},
		}
		context.SetBounds(&r.steps, gl.CellBounds(0, 0), r)
		context.SetBounds(&r.navbar, gl.CellBounds(0, 1), r)
	}

	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title:          "Installer",
		WindowSize:     image.Pt(800, 600),
		RunGameOptions: &ebiten.RunGameOptions{
			// ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
