package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Step4 struct {
	guigui.DefaultWidget

	bodyText basicwidget.Text

	model *Model
}

func (s *Step4) SetModel(model *Model) {
	s.model = model
}

func (s *Step4) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {

	s.bodyText.SetMultiline(true)
	s.bodyText.SetAutoWrap(true)
	s.bodyText.SetValue(s.model.CurrentStep().BodyText)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&s.bodyText, gl.CellBounds(0, 0))

	return nil
}
