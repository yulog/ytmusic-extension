package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Step3 struct {
	guigui.DefaultWidget

	buttonsForm basicwidget.Form
	buttonText  basicwidget.TextInput
	button      basicwidget.Button

	configForm basicwidget.Form
	buttonBack basicwidget.Button
	buttonNext basicwidget.Button

	model *Model
}

func (b *Step3) SetModel(model *Model) {
	b.model = model
}

func (b *Step3) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)

	b.buttonText.SetValue("Button")
	b.button.SetText("Button")
	context.SetEnabled(&b.button, b.model.Buttons().Enabled())

	b.buttonsForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &b.buttonText,
			SecondaryWidget: &b.button,
		},
	})

	b.buttonBack.SetText("Back")
	b.buttonBack.SetOnUp(func() {
		b.model.SetStep(step2)
	})
	b.buttonNext.SetText("Next")
	b.buttonNext.SetOnUp(func() {
		b.model.SetStep(step4)
	})

	b.configForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &b.buttonBack,
			SecondaryWidget: &b.buttonNext,
		},
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(b).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(b.buttonsForm.DefaultSize(context).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(b.configForm.DefaultSize(context).Y),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&b.buttonsForm, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&b.configForm, gl.CellBounds(0, 2))

	return nil
}
