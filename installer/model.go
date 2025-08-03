package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/yulog/ytmusic-extension/dist"
	"github.com/yulog/ytmusic-extension/tools"
)

const (
	step0 = iota
	step1
	step2
	step3
	step4
	step5
)

const (
	name = "YouTube Music Notifier"
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

	output string
}

type StepModel struct {
	ID            int
	TitleText     string
	BodyText      string
	Func          FuncModel
	CancelFunc    func()
	BackButton    EnabledModel
	ConfirmButton DisabledModel
	Form          DisabledModel
	Last          bool
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
		log.SetOutput(io.MultiWriter(s, os.Stdout))
		s.steps = []StepModel{
			{
				TitleText: "Welcome",
				ID:        step0,
				BodyText: fmt.Sprintf(`このウィザードは、%s をあなたのコンピューターにインストールするお手伝いをします。

続行するには、「次へ」をクリックしてください。`, name),
			},
			{
				TitleText: "Location",
				ID:        step1,
				BodyText:  fmt.Sprintf(`%s をインストールするフォルダーを選択してください。デフォルトの場所を使用するか、「参照...」をクリックして別のフォルダーを選択できます。`, name),
				Func: FuncModel{
					Text: "Extract",
					Func: func() {
						log.Println("=== Extract ===")
						sub, _ := fs.Sub(dist.Contents, "contents")
						os.CopyFS(s.location, sub)
						log.Printf("Created files at: %s", s.location)
						log.Println("=== Extract ===")
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
				BodyText: fmt.Sprintf(`このステップでは、%s で使用する拡張機能を読み込みます。

1. chrome://extensions からデベロッパーモードを有効化します。
2. [パッケージ化されていない拡張機能を読み込む]から前のステップで展開された chrome-extension フォルダを指定します。`, name),
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
				BodyText:  fmt.Sprintf(`%s は、設定の一部をシステムのレジストリ（Windowsの場合）に保存する必要があります。これにより、アプリケーションが正しく動作します。`, name),
				Func: FuncModel{
					Text: "Register",
					Func: func() {
						log.Println("=== Register ===")
						// log.Println("Register!!", s.extensionID)
						manifestPath := tools.CreateManifest(s.location, "native-app.exe", s.extensionID)
						tools.Register(manifestPath)
						time.Sleep(3 * time.Second)
						log.Println("3 seconds")
						log.Println("=== Register ===")
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
				BodyText: fmt.Sprintf(`%s のインストールが正常に完了しました。
				
「完了」をクリックしてインストーラーを終了してください。`, name),
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
		s.location = filepath.Join(dir, "ytmusic-extension")
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

func (s *StepsModel) Write(p []byte) (n int, err error) {
	s.output += string(p)
	return len(p), nil
}
