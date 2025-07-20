package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Navbar struct {
	guigui.DefaultWidget

	panel        basicwidget.Panel
	panelContent navigationPanelContent
}

func (n *Navbar) SetModel(model *Model) {
	n.panelContent.SetModel(model)
}

func (n *Navbar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	n.panel.SetBorder(basicwidget.PanelBorder{
		Top: true,
	})
	context.SetSize(&n.panelContent, context.ActualSize(n), n)
	n.panel.SetContent(&n.panelContent)

	appender.AppendChildWidgetWithBounds(&n.panel, context.Bounds(n))

	return nil
}

// type navigationWidget struct {
// 	guigui.DefaultWidget

// 	buttonBack   basicwidget.Button
// 	buttonNext   basicwidget.Button
// 	buttonCancel basicwidget.Button
// }

// func (n *navigationWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
// 	n.buttonBack.SetText("Back")
// 	n.buttonBack.SetOnUp(func() {
// 		// n.model.SetStep(step1)
// 	})
// 	// p.text.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
// 	n.buttonNext.SetText("Next")
// 	n.buttonCancel.SetText("Cancel")

// 	u := basicwidget.UnitSize(context)
// 	gl := layout.GridLayout{
// 		Bounds: context.Bounds(n),
// 		Widths: []layout.Size{
// 			layout.FlexibleSize(1),
// 			layout.FixedSize(3 * u),
// 			layout.FixedSize(3 * u),
// 			layout.FixedSize(3 * u),
// 		},
// 		ColumnGap: u / 2,
// 	}
// 	appender.AppendChildWidgetWithBounds(&n.buttonBack, gl.CellBounds(1, 0))
// 	appender.AppendChildWidgetWithBounds(&n.buttonNext, gl.CellBounds(2, 0))
// 	appender.AppendChildWidgetWithBounds(&n.buttonCancel, gl.CellBounds(3, 0))

// 	return nil
// }

// func (n *navigationWidget) DefaultSize(context *guigui.Context) image.Point {
// 	return image.Pt(6*int(basicwidget.UnitSize(context)), context.ActualSize(&n.buttonNext).Y)
// }

type navigationPanelContent struct {
	guigui.DefaultWidget

	// navigationWidget navigationWidget
	// onClearTriggered func()
	buttonBack   basicwidget.Button
	buttonNext   basicwidget.Button
	buttonCancel basicwidget.Button

	model *Model
}

func (s *navigationPanelContent) SetModel(model *Model) {
	s.model = model
}

func (n *navigationPanelContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {

	// u := basicwidget.UnitSize(context)
	// gl := layout.GridLayout{
	// 	Bounds: context.Bounds(n).Inset(u / 2),
	// 	Heights: []layout.Size{
	// 		layout.FixedSize(n.navigationWidget.DefaultSize(context).Y),
	// 	},
	// 	ColumnGap: u / 2,
	// }
	// appender.AppendChildWidgetWithBounds(&n.navigationWidget, gl.CellBounds(0, 0))

	n.buttonBack.SetText("Back")
	n.buttonBack.SetOnUp(func() {
		n.model.BackStep()
	})
	// p.text.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	n.buttonNext.SetText("Next")
	n.buttonNext.SetOnUp(func() {
		n.model.NextStep()
	})
	n.buttonCancel.SetText("Cancel")

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(n).Inset(u / 2),
		Widths: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(3 * u),
			layout.FixedSize(3 * u),
			layout.FixedSize(3 * u),
		},
		Heights: []layout.Size{
			layout.FixedSize(u),
		},
		ColumnGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&n.buttonBack, gl.CellBounds(1, 0))
	appender.AppendChildWidgetWithBounds(&n.buttonNext, gl.CellBounds(2, 0))
	appender.AppendChildWidgetWithBounds(&n.buttonCancel, gl.CellBounds(3, 0))

	return nil
}

func (n *navigationPanelContent) DefaultSize(context *guigui.Context) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(6*int(u), context.ActualSize(&n.buttonNext).Y+u/2)
}
