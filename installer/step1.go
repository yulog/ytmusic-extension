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
}

func (s *Step1) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.bodyText)
	appender.AppendChildWidget(&s.buttonsForm)
}

func (s *Step1) Build(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*Model)

	u := basicwidget.UnitSize(context)

	s.bodyText.SetMultiline(true)
	s.bodyText.SetAutoWrap(true)
	s.bodyText.SetValue(model.CurrentStep().BodyText)

	s.buttonText.SetValue(model.Steps().Location())
	s.buttonText.SetOnValueChanged(func(text string, committed bool) {
		if committed {
			model.Steps().SetLocation(text)
		}
	})
	s.button.SetText("Ref...")
	context.SetEnabled(&s.buttonsForm, model.CurrentStep().Form.Enabled())

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
			layout.FixedSize(s.bodyText.DefaultSizeInContainer(context, context.ActualSize(&s.bodyText).X).Y),
			layout.FixedSize(s.buttonsForm.DefaultSize(context).Y),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	context.SetBounds(&s.bodyText, gl.CellBounds(0, 0), s)
	context.SetBounds(&s.buttonsForm, gl.CellBounds(0, 1), s)

	return nil
}
