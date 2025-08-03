package main

import (
	"image"
	"os"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Navbar struct {
	guigui.DefaultWidget

	panel        basicwidget.Panel
	panelContent navigationPanelContent
}

func (n *Navbar) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&n.panel)
}

func (n *Navbar) Build(context *guigui.Context) error {
	n.panel.SetBorder(basicwidget.PanelBorder{
		Top: true,
	})
	context.SetSize(&n.panelContent, context.ActualSize(n), n)
	n.panel.SetContent(&n.panelContent)

	context.SetBounds(&n.panel, context.Bounds(n), n)

	return nil
}

type navigationPanelContent struct {
	guigui.DefaultWidget

	buttonBack    basicwidget.Button
	buttonNext    basicwidget.Button
	buttonConfirm basicwidget.Button
	buttonCancel  basicwidget.Button
}

func (n *navigationPanelContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	model := context.Model(n, modelKeyModel).(*Model)

	appender.AppendChildWidget(&n.buttonBack)
	if model.CurrentStep().Func.Text != "" {
		appender.AppendChildWidget(&n.buttonConfirm)
	} else {
		appender.AppendChildWidget(&n.buttonNext)
	}
	appender.AppendChildWidget(&n.buttonCancel)
}

func (n *navigationPanelContent) Build(context *guigui.Context) error {
	model := context.Model(n, modelKeyModel).(*Model)

	context.SetEnabled(&n.buttonBack, model.CurrentStep().BackButton.Enabled())
	n.buttonBack.SetText("Back")
	n.buttonBack.SetOnUp(func() {
		model.BackStep()
	})

	n.buttonNext.SetText("Next")
	n.buttonNext.SetOnUp(func() {
		model.NextStep()
	})

	if s := model.CurrentStep(); s.Func.Text != "" {
		context.SetEnabled(&n.buttonConfirm, s.ConfirmButton.Enabled())
		n.buttonConfirm.SetText(s.Func.Text)
		n.buttonConfirm.SetOnUp(func() {
			if s.BackButton.Enabled() {
				s.BackButton.SetEnabled(false)
			}
			if s.Form.Enabled() {
				s.Form.SetEnabled(false)
			}
			if s.ConfirmButton.Enabled() {
				s.ConfirmButton.SetEnabled(false)
			}
			go func() {
				s.Func.Func()
				if !s.Last {
					s.Func.Text = ""
				}

			}()
		})
	}

	n.buttonCancel.SetText("Cancel")
	if s := model.CurrentStep(); s.CancelFunc != nil {
		n.buttonCancel.SetOnUp(s.CancelFunc)
	} else {
		n.buttonCancel.SetOnUp(func() {
			os.Exit(1)
		})
	}

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
	context.SetBounds(&n.buttonBack, gl.CellBounds(1, 0), n)
	if model.CurrentStep().Func.Text != "" {
		context.SetBounds(&n.buttonConfirm, gl.CellBounds(2, 0), n)
	} else {
		context.SetBounds(&n.buttonNext, gl.CellBounds(2, 0), n)
	}
	context.SetBounds(&n.buttonCancel, gl.CellBounds(3, 0), n)

	return nil
}

func (n *navigationPanelContent) DefaultSize(context *guigui.Context) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(6*int(u), context.ActualSize(&n.buttonNext).Y+u/2)
}
