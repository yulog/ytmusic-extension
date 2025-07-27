package main

import (
	"fmt"
	"os"
)

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

	steps StepsModel
}

func (m *Model) Step() int {
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

func (m *Model) Steps() *StepsModel {
	return &m.steps
}

func (m *Model) CurrentStep() *StepModel {
	return m.Steps().ByIndex(m.Step())
}

type StepsModel struct {
	steps []StepModel

	location    string
	extensionID string
}

type StepModel struct {
	ID         int
	TitleText  string
	BodyText   string
	Func       FuncModel
	CancelFunc func()
	BackButton EnabledModel
	Form       DisabledModel
	Last       bool
}

type FuncModel struct {
	Text string
	Func func()
}

type DisabledModel struct {
	disabled bool
}

func (d *DisabledModel) Enabled() bool {
	return !d.disabled
}

func (d *DisabledModel) SetEnabled(enabled bool) {
	d.disabled = !enabled
}

type EnabledModel struct {
	enabled bool
}

func (e *EnabledModel) Enabled() bool {
	return e.enabled
}

func (e *EnabledModel) SetEnabled(enabled bool) {
	e.enabled = enabled
}

func (s *StepsModel) Steps() []StepModel {
	if s.steps == nil {
		s.steps = []StepModel{
			{
				TitleText: "Welcome",
				ID:        step0,
				BodyText: `このウィザードは、[製品名] をあなたのコンピューターにインストールするお手伝いをします。インストールを開始する前に、すべての実行中のアプリケーションを閉じることをお勧めします。

続行するには、「次へ」をクリックしてください。`,
			},
			{
				TitleText: "Location",
				ID:        step1,
				BodyText:  `[製品名] をインストールするフォルダーを選択してください。デフォルトの場所を使用するか、「参照...」をクリックして別のフォルダーを選択できます。`,
				Func: FuncModel{
					Text: "Extract",
					Func: func() {
						fmt.Println("Extract!!", s.location)
					},
				},
				CancelFunc: func() {
					os.Exit(1)
				},
				BackButton: EnabledModel{
					enabled: true,
				},
			},
			{
				TitleText: "Load Extension",
				ID:        step2,
				BodyText:  `このステップでは、[製品名] で使用する拡張機能を読み込みます。`,
				CancelFunc: func() {
					os.Exit(1)
				},
			},
			{
				TitleText: "Extension ID",
				ID:        step3,
				BodyText:  `前のステップで読み込んだ拡張機能のIDを入力してください。`,
				CancelFunc: func() {
					os.Exit(1)
				},
				BackButton: EnabledModel{
					enabled: true,
				},
			},
			{
				TitleText: "Registry",
				ID:        step4,
				BodyText:  `[製品名] は、設定の一部をシステムのレジストリ（Windowsの場合）に保存する必要があります。これにより、アプリケーションが正しく動作します。`,
				Func: FuncModel{
					Text: "Register",
					Func: func() {
						fmt.Println("Register!!", s.extensionID)
					},
				},
				CancelFunc: func() {
					os.Exit(1)
				},
				BackButton: EnabledModel{
					enabled: true,
				},
			},
			{
				TitleText: "Finish",
				ID:        step5,
				BodyText: `[製品名] のインストールが正常に完了しました。
				
「完了」をクリックしてインストーラーを終了してください。`,
				Func: FuncModel{
					Text: "Finish",
					Func: func() {
						fmt.Println("Finish!!")
						os.Exit(0)
					},
				},
				CancelFunc: func() {
					os.Exit(1)
				},
				Last: true,
			},
		}
	}
	return s.steps
}

func (s *StepsModel) ByIndex(i int) *StepModel {
	return &s.Steps()[i]
}

func (s *StepsModel) Location() string {
	if s.location == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return ""
		}
		s.location = dir
	}
	return s.location
}

func (s *StepsModel) SetLocation(location string) {
	s.location = location
}

func (s *StepsModel) ExtensionID() string {
	return s.extensionID
}

func (s *StepsModel) SetExtensionID(id string) {
	s.extensionID = id
}
