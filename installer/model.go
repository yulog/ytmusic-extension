package main

const (
	step0 = iota
	step1
	step2
	step3
	step4
	step5
)

type Model struct {
	step int

	buttons ButtonsModel
	texts   TextsModel
}

func (m *Model) Step() int {
	// if m.step == "" {
	// 	return "step0"
	// }
	return m.step
}

func (m *Model) SetStep(step int) {
	m.step = step
}

func (m *Model) NextStep() int {
	if m.step+1 > step5 {
		return m.step
	}
	m.step++
	return m.step
}

func (m *Model) BackStep() int {
	if m.step-1 < step0 {
		return m.step
	}
	m.step--
	return m.step
}

func (m *Model) Buttons() *ButtonsModel {
	return &m.buttons
}

func (m *Model) Texts() *TextsModel {
	return &m.texts
}

type ButtonsModel struct {
	disabled bool
}

func (b *ButtonsModel) Enabled() bool {
	return !b.disabled
}

func (b *ButtonsModel) SetEnabled(enabled bool) {
	b.disabled = !enabled
}

type TextsModel struct {
	text    string
	textSet bool
}

func (t *TextsModel) Text() string {
	if !t.textSet {
		return `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
隴西の李徴は博学才穎、天宝の末年、若くして名を虎榜に連ね、ついで江南尉に補せられたが、性、狷介、自ら恃むところ頗る厚く、賤吏に甘んずるを潔しとしなかった。`
	}
	return t.text
}

func (t *TextsModel) SetText(text string) {
	t.text = text
	t.textSet = true
}
