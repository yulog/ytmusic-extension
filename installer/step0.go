package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Step0 struct {
	guigui.DefaultWidget

	form       basicwidget.Form
	buttonBack basicwidget.Button
	buttonNext basicwidget.Button
	// pageSegmentedControl basicwidget.SegmentedControl[string]
	buttonCancel           basicwidget.Button
	titleText              basicwidget.Text
	bodyText               basicwidget.Text
	navigationPanel        basicwidget.Panel
	navigationPanelContent navigationPanelContent

	model *Model
}

func (t *Step0) SetModel(model *Model) {
	t.model = model
}

func (t *Step0) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {

	t.titleText.SetMultiline(true)
	t.titleText.SetAutoWrap(true)
	t.titleText.SetBold(true)
	t.titleText.SetScale(2)
	t.titleText.SetValue("Welcome to Installer")

	t.bodyText.SetMultiline(true)
	// t.sampleText.SetHorizontalAlign(basicwidget.HorizontalAlignLeft)
	// t.sampleText.SetVerticalAlign(basicwidget.VerticalAlignTop)
	t.bodyText.SetAutoWrap(true)
	t.bodyText.SetValue(t.model.Texts().Text())

	t.buttonBack.SetText("Back")
	t.buttonNext.SetText("Next")
	t.buttonNext.SetOnUp(func() {
		t.model.SetStep(step1)
	})

	// t.pageSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[string]{
	// 	{
	// 		Text: "Back",
	// 		ID:   "back",
	// 	},
	// 	{
	// 		Text: "Next",
	// 		ID:   "next",
	// 	},
	// })
	// t.pageSegmentedControl.SetOnItemSelected(func(index int) {
	// 	item, ok := t.pageSegmentedControl.ItemByIndex(index)
	// 	if !ok {
	// 		return
	// 	}
	// 	switch item.ID {
	// 	case "back":
	// 		// t.model.SetMode("step1")
	// 	case "next":
	// 		t.model.SetMode("step1")
	// 	default:
	// 	}
	// })
	// t.pageSegmentedControl.SelectItemByID("")

	t.buttonCancel.SetText("Cancel")
	t.buttonCancel.SetOnUp(func() {
		// TODO: Cancel sequence
	})

	t.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &t.buttonBack,
			SecondaryWidget: &t.buttonNext,
		},
		{
			SecondaryWidget: &t.buttonCancel,
		},
	})
	// TODO: Formではなく、PanelとGridにする
	t.navigationPanel.SetContent(&t.navigationPanelContent)
	t.navigationPanel.SetBorder(basicwidget.PanelBorder{Top: true})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(t.titleText.DefaultSize(context).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(t.form.DefaultSize(context).Y),
			layout.FixedSize(u + u/2),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&t.titleText, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&t.bodyText, gl.CellBounds(0, 1))
	appender.AppendChildWidgetWithBounds(&t.form, gl.CellBounds(0, 2))
	context.SetSize(&t.navigationPanelContent, image.Pt(gl.CellBounds(0, 3).Dx(), guigui.AutoSize), t)
	appender.AppendChildWidgetWithBounds(&t.navigationPanel, gl.CellBounds(0, 3))

	return nil
}
