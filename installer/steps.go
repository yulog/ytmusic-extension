package main

import (
	"github.com/hajimehoshi/guigui"
)

type Steps struct {
	guigui.DefaultWidget

	step0 Step0
	step1 Step1
	step2 Step2
	step3 Step3
	step4 Step4
	step5 Step5

	model *Model
}

func (s *Steps) SetModel(model *Model) {
	s.model = model
	s.step0.SetModel(model)
	s.step1.SetModel(model)
	s.step2.SetModel(model)
	s.step3.SetModel(model)
	s.step4.SetModel(model)
	s.step5.SetModel(model)
}

func (s *Steps) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	switch s.model.Step() {
	case step0:
		appender.AppendChildWidgetWithBounds(&s.step0, context.Bounds(s))
	case step1:
		appender.AppendChildWidgetWithBounds(&s.step1, context.Bounds(s))
	case step2:
		appender.AppendChildWidgetWithBounds(&s.step2, context.Bounds(s))
	case step3:
		appender.AppendChildWidgetWithBounds(&s.step3, context.Bounds(s))
	case step4:
		appender.AppendChildWidgetWithBounds(&s.step4, context.Bounds(s))
	case step5:
		appender.AppendChildWidgetWithBounds(&s.step5, context.Bounds(s))
	}

	return nil
}
