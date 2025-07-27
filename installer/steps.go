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

func (s *Steps) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	switch s.model.Step() {
	case step0:
		appender.AppendChildWidget(&s.step0)
	case step1:
		appender.AppendChildWidget(&s.step1)
	case step2:
		appender.AppendChildWidget(&s.step2)
	case step3:
		appender.AppendChildWidget(&s.step3)
	case step4:
		appender.AppendChildWidget(&s.step4)
	case step5:
		appender.AppendChildWidget(&s.step5)
	}
}

func (s *Steps) Build(context *guigui.Context) error {
	switch s.model.Step() {
	case step0:
		context.SetBounds(&s.step0, context.Bounds(s), s)
	case step1:
		context.SetBounds(&s.step1, context.Bounds(s), s)
	case step2:
		context.SetBounds(&s.step2, context.Bounds(s), s)
	case step3:
		context.SetBounds(&s.step3, context.Bounds(s), s)
	case step4:
		context.SetBounds(&s.step4, context.Bounds(s), s)
	case step5:
		context.SetBounds(&s.step5, context.Bounds(s), s)
	}

	return nil
}
