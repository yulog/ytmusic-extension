package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Step4 struct {
	guigui.DefaultWidget

	bodyText basicwidget.Text
}

func (s *Step4) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.bodyText)
}

func (s *Step4) Build(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*Model)

	s.bodyText.SetMultiline(true)
	s.bodyText.SetAutoWrap(true)
	s.bodyText.SetValue(model.CurrentStep().BodyText)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	context.SetBounds(&s.bodyText, gl.CellBounds(0, 0), s)

	return nil
}
