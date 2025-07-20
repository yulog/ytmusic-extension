package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Step4 struct {
	guigui.DefaultWidget

	form       basicwidget.Form
	buttonBack basicwidget.Button
	buttonNext basicwidget.Button
	sampleText basicwidget.Text

	model *Model
}

func (t *Step4) SetModel(model *Model) {
	t.model = model
}

func (t *Step4) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {

	// t.buttonPrev.SetText("Prev")
	t.buttonNext.SetText("Finish")

	t.form.SetItems([]basicwidget.FormItem{
		{
			// PrimaryWidget:   &t.buttonPrev,
			SecondaryWidget: &t.buttonNext,
		},
	})

	t.sampleText.SetMultiline(true)
	// t.sampleText.SetHorizontalAlign(basicwidget.HorizontalAlignLeft)
	// t.sampleText.SetVerticalAlign(basicwidget.VerticalAlignTop)
	t.sampleText.SetAutoWrap(true)
	// t.sampleText.SetBold(false)
	// t.sampleText.SetSelectable(false)
	// t.sampleText.SetEditable(false)
	// t.sampleText.SetOnValueChanged(func(text string, committed bool) {
	// 	if committed {
	// 		t.model.Texts().SetText(text)
	// 	}
	// })
	// t.sampleText.SetOnKeyJustPressed(func(key ebiten.Key) bool {
	// 	if !t.sampleText.IsEditable() {
	// 		return false
	// 	}
	// 	if key == ebiten.KeyTab {
	// 		t.sampleText.ReplaceValueAtSelection("\t")
	// 		return true
	// 	}
	// 	return false
	// })
	t.sampleText.SetValue(t.model.Texts().Text())

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(t.form.DefaultSize(context).Y),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&t.sampleText, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&t.form, gl.CellBounds(0, 1))

	return nil
}
