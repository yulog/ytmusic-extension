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

func (s *Sidebar) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.panel)
}

func (s *Sidebar) Build(context *guigui.Context) error {
	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetSize(&s.panelContent, context.ActualSize(s), s)
	s.panel.SetContent(&s.panelContent)

	context.SetBounds(&s.panel, context.Bounds(s), s)

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.List[int]
}

func (s *sidebarContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.list)
}

func (s *sidebarContent) Build(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*Model)

	s.list.SetStyle(basicwidget.ListStyleSidebar)

	var items []basicwidget.ListItem[int]
	for _, v := range model.Steps().Steps() {
		items = append(items, basicwidget.ListItem[int]{
			Text: v.TitleText,
			ID:   v.ID,
		})
	}

	s.list.SetItems(items)
	s.list.SelectItemByID(model.Step())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			model.SetStep(0)
			return
		}
		model.SetStep(item.ID)
	})
	context.SetEnabled(&s.list, false)

	context.SetBounds(&s.list, context.Bounds(s), s)

	return nil
}
