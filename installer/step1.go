package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Step1 struct {
	guigui.DefaultWidget

	bodyText basicwidget.Text

	buttonsForm basicwidget.Form
	buttonText  basicwidget.TextInput
	button      basicwidget.Button

	model *Model
}

func (s *Step1) SetModel(model *Model) {
	s.model = model
}

func (s *Step1) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)

	s.bodyText.SetMultiline(true)
	s.bodyText.SetAutoWrap(true)
	s.bodyText.SetValue(s.model.CurrentStep().BodyText)

	s.buttonText.SetValue(s.model.Steps().Location())
	s.buttonText.SetOnValueChanged(func(text string, committed bool) {
		if committed {
			s.model.Steps().SetLocation(text)
		}
	})
	s.button.SetText("Ref...")
	context.SetEnabled(&s.buttonsForm, s.model.CurrentStep().Form.Enabled())

	s.buttonsForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &s.buttonText,
			SecondaryWidget: &s.button,
		},
	})
	context.SetSize(&s.buttonText, image.Point{X: context.ActualSize(&s.buttonsForm).X - context.ActualSize(&s.button).X - u, Y: s.buttonText.DefaultSize(context).Y}, s)

	gl := layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(s.bodyText.TextSize(context, context.ActualSize(&s.bodyText).X).Y),
			layout.FixedSize(s.buttonsForm.DefaultSize(context).Y),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&s.bodyText, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&s.buttonsForm, gl.CellBounds(0, 1))

	return nil
}
