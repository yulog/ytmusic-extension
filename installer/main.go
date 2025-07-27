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

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.updateFontFaceSources(context)

	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	r.sidebar.SetModel(&r.model)
	r.steps.SetModel(&r.model)
	r.navbar.SetModel(&r.model)

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&r.sidebar, gl.CellBounds(0, 0))
	{
		u := basicwidget.UnitSize(context)
		gl := layout.GridLayout{
			Bounds: gl.CellBounds(1, 0).Inset(u / 2),
			Heights: []layout.Size{
				layout.FlexibleSize(1),
				layout.FixedSize(u + u/2),
			},
		}
		appender.AppendChildWidgetWithBounds(&r.steps, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&r.navbar, gl.CellBounds(0, 1))
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
