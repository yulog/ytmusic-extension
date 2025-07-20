package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	panel        basicwidget.Panel
	panelContent sidebarContent
}

func (s *Sidebar) SetModel(model *Model) {
	s.panelContent.SetModel(model)
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetSize(&s.panelContent, context.ActualSize(s), s)
	s.panel.SetContent(&s.panelContent)

	appender.AppendChildWidgetWithBounds(&s.panel, context.Bounds(s))

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.List[int]

	model *Model
}

func (s *sidebarContent) SetModel(model *Model) {
	s.model = model
}

func (s *sidebarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.ListItem[int]{
		{
			Text: "Welcome",
			ID:   step0,
		},
		{
			Text: "Location",
			ID:   step1,
		},
		{
			Text: "Load Extension",
			ID:   step2,
		},
		{
			Text: "Extension ID",
			ID:   step3,
		},
		{
			Text: "Registry",
			ID:   step4,
		},
		{
			Text: "Finish",
			ID:   step5,
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByID(s.model.Step())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			s.model.SetStep(0)
			return
		}
		s.model.SetStep(item.ID)
	})
	context.SetEnabled(&s.list, false)

	appender.AppendChildWidgetWithBounds(&s.list, context.Bounds(s))

	return nil
}
