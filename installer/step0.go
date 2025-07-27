package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Step0 struct {
	guigui.DefaultWidget

	titleText basicwidget.Text
	bodyText  basicwidget.Text

	model *Model
}

func (s *Step0) SetModel(model *Model) {
	s.model = model
}

func (s *Step0) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.titleText)
	appender.AppendChildWidget(&s.bodyText)
}

func (s *Step0) Build(context *guigui.Context) error {

	s.titleText.SetMultiline(true)
	s.titleText.SetAutoWrap(true)
	s.titleText.SetBold(true)
	s.titleText.SetScale(2)
	s.titleText.SetValue("Welcome to Installer")

	s.bodyText.SetMultiline(true)
	s.bodyText.SetAutoWrap(true)
	s.bodyText.SetValue(s.model.CurrentStep().BodyText)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(s.titleText.DefaultSize(context).Y),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	context.SetBounds(&s.titleText, gl.CellBounds(0, 0), s)
	context.SetBounds(&s.bodyText, gl.CellBounds(0, 1), s)

	return nil
}
