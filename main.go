package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/yulog/ytmusic-extension/resources/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// --- Data Structures ---
type SongInfo struct {
	Source     string `json:"source"`
	Title      string `json:"title"`
	Byline     string `json:"byline,omitempty"`
	Artist     string `json:"artist,omitempty"`
	Album      string `json:"album,omitempty"`
	ArtworkURL string `json:"artworkUrl"`
}

type Message struct {
	Text string `json:"text"`
}

// --- Ebitengine Game ---
const (
	screenWidth  = 400
	screenHeight = 120
	artworkSize  = 100
	baseFontSize = 18
)

var (
	songInfoChan = make(chan SongInfo, 1)
)

type Game struct {
	currentSong       SongInfo
	artwork           *ebiten.Image
	currentArtworkURL string
	mu                sync.RWMutex
	dragger           *WindowDragger
	fontFace          *text.GoTextFace
	fontInitialized   bool
	titleScrollX      float64
	line2ScrollX      float64
	titleTextWidth    float64
	line2TextWidth    float64
	deviceScaleFactor float64 // To cache the device scale factor
}

// WindowDragger handles the logic for dragging the game window.
type WindowDragger struct {
	dragging         bool
	dragStartWindowX int
	dragStartWindowY int
	dragStartCursorX int
	dragStartCursorY int
	cursorToWindowX  float64
	cursorToWindowY  float64
}

func NewWindowDragger() *WindowDragger {
	return &WindowDragger{}
}

func (d *WindowDragger) Update(g *Game) {
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		d.dragging = false
	}
	if !d.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		d.dragging = true
		d.dragStartWindowX, d.dragStartWindowY = ebiten.WindowPosition()
		d.dragStartCursorX, d.dragStartCursorY = ebiten.CursorPosition()
	}

	if d.dragging {
		cx, cy := ebiten.CursorPosition()
		dx := int(float64(cx-d.dragStartCursorX) * d.cursorToWindowX)
		dy := int(float64(cy-d.dragStartCursorY) * d.cursorToWindowY)
		wx, wy := ebiten.WindowPosition()
		ebiten.SetWindowPosition(wx+dx, wy+dy)
	}
}

func (g *Game) initFont() {
	// file, err := os.Open("Mplus1p-Regular.ttf")
	// if err != nil {
	// 	log.Fatalf("Failed to open font file: %v", err)
	// }
	// defer file.Close() // Ensure the file is closed after use

	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular))
	if err != nil {
		log.Fatal(err)
	}
	scale := g.deviceScaleFactor
	log.Printf("HiDPI: Device scale factor is %.2f", scale)
	g.fontFace = &text.GoTextFace{
		Source: source,
		Size:   baseFontSize * scale,
	}
	g.fontInitialized = true
}

func (g *Game) Update() error {
	// Handle window dragging
	g.dragger.Update(g)

	// Handle song info updates
	select {
	case newSong := <-songInfoChan:
		g.mu.Lock()
		g.currentSong = newSong
		if g.currentSong.ArtworkURL != "" && g.currentSong.ArtworkURL != g.currentArtworkURL {
			g.currentArtworkURL = g.currentSong.ArtworkURL
			go g.downloadArtwork(g.currentArtworkURL)
		}
		g.resetTextScroll()
		g.mu.Unlock()
	default:
	}

	// Update text scrolling animation
	g.updateTextScroll()

	return nil
}

// updateTextScroll handles the logic for scrolling long text.
func (g *Game) updateTextScroll() {
	g.mu.Lock()
	defer g.mu.Unlock()

	s := g.deviceScaleFactor
	textX := (10 + artworkSize + 10) * s
	textAreaWidth := float64(screenWidth*s) - textX - (10 * s)

	// Title scrolling
	if g.titleTextWidth > textAreaWidth {
		g.titleScrollX -= 0.5 // Adjust scroll speed as needed
		if g.titleScrollX < -g.titleTextWidth {
			g.titleScrollX = textAreaWidth
		}
	} else {
		g.titleScrollX = 0
	}

	// Line 2 scrolling
	if g.line2TextWidth > textAreaWidth {
		g.line2ScrollX -= 0.5 // Adjust scroll speed as needed
		if g.line2ScrollX < -g.line2TextWidth {
			g.line2ScrollX = textAreaWidth
		}
	} else {
		g.line2ScrollX = 0
	}
}

// resetTextScroll resets the text scroll position and recalculates text widths.
func (g *Game) resetTextScroll() {
	g.titleScrollX = 0
	g.line2ScrollX = 0
	if g.fontInitialized {
		title := g.currentSong.Title
		line2 := ""
		if g.currentSong.Source == "MediaSessionAPI" {
			line2 = fmt.Sprintf("%s - %s", g.currentSong.Artist, g.currentSong.Album)
		} else {
			line2 = g.currentSong.Byline
		}
		g.titleTextWidth, _ = text.Measure(title, g.fontFace, 0)
		g.line2TextWidth, _ = text.Measure(line2, g.fontFace, 0)
	}
}

func (g *Game) downloadArtwork(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download artwork: %v", err)
		return
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Failed to decode artwork: %v", err)
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.currentArtworkURL == url {
		g.artwork = ebiten.NewImageFromImage(img)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	s := g.deviceScaleFactor

	// Draw artwork
	if g.artwork != nil {
		op := &ebiten.DrawImageOptions{}
		imgW, imgH := g.artwork.Bounds().Dx(), g.artwork.Bounds().Dy()
		scale := float64(artworkSize*s) / float64(imgW)
		if imgH > imgW {
			scale = float64(artworkSize*s) / float64(imgH)
		}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(10*s, float64(screenHeight*s-artworkSize*s)/2)
		screen.DrawImage(g.artwork, op)
	}

	// Define the text area
	textX := (10 + artworkSize + 10) * s
	textY1 := 40 * s
	textY2 := 80 * s
	textAreaWidth := float64(screenWidth*s) - textX - (10 * s)

	// Create a sub-image for the entire text area to clip the text
	textClipArea := ebiten.NewImage(int(textAreaWidth), int(screenHeight*s))
	textClipArea.Fill(color.Transparent)

	// --- Draw Title ---
	title := g.currentSong.Title
	if title == "" {
		title = "Waiting for song..."
	}
	titleOp := &text.DrawOptions{}
	titleOp.ColorScale.ScaleWithColor(color.White)
	// Draw text at a position relative to the textClipArea, considering the scroll
	titleOp.GeoM.Translate(g.titleScrollX, textY1)
	text.Draw(textClipArea, title, g.fontFace, titleOp)

	// --- Draw Line 2 ---
	line2 := ""
	if g.currentSong.Source == "MediaSessionAPI" {
		line2 = fmt.Sprintf("%s - %s", g.currentSong.Artist, g.currentSong.Album)
	} else {
		line2 = g.currentSong.Byline
	}

	if line2 != "" {
		line2Op := &text.DrawOptions{}
		line2Op.ColorScale.ScaleWithColor(color.White)
		line2Op.GeoM.Translate(g.line2ScrollX, textY2)
		text.Draw(textClipArea, line2, g.fontFace, line2Op)
	}

	// Draw the clipped text area onto the main screen
	clipOp := &ebiten.DrawImageOptions{}
	clipOp.GeoM.Translate(textX, 0)
	screen.DrawImage(textClipArea, clipOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Update the device scale factor whenever the layout changes.
	g.deviceScaleFactor = ebiten.Monitor().DeviceScaleFactor()

	// Calculate the factors to convert a cursor position to a window position.
	g.dragger.cursorToWindowX = float64(outsideWidth) / (float64(screenWidth) * g.deviceScaleFactor)
	g.dragger.cursorToWindowY = float64(outsideHeight) / (float64(screenHeight) * g.deviceScaleFactor)

	return int(float64(screenWidth) * g.deviceScaleFactor), int(float64(screenHeight) * g.deviceScaleFactor)
}

// --- Main and Message Reader ---
func main() {
	setupLogging()
	go readMessages() // Run the message reader in a separate goroutine

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("YouTube Music Notifier")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)

	sw, sh := ebiten.Monitor().Size()
	ebiten.SetWindowPosition(sw-screenWidth-10, sh-screenHeight-60)

	game := &Game{
		dragger:           NewWindowDragger(),
		deviceScaleFactor: ebiten.Monitor().DeviceScaleFactor(), // Initial setting
	}

	game.initFont() // Initialize font before starting the game

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func readMessages() {
	log.Println("Message reader goroutine started.")
	for {
		var length uint32
		if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				log.Println("Stdin closed, exiting message reader.")
				break
			}
			log.Fatalf("Failed to read message length: %v", err)
		}
		if length == 0 {
			continue
		}
		msgBytes := make([]byte, length)
		if _, err := io.ReadFull(os.Stdin, msgBytes); err != nil {
			log.Fatalf("Failed to read message content: %v", err)
		}
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Failed to unmarshal outer message: %v", err)
			continue
		}
		var songInfo SongInfo
		if err := json.Unmarshal([]byte(msg.Text), &songInfo); err != nil {
			log.Printf("Failed to unmarshal song info: %v", err)
			continue
		}
		logSongInfo(songInfo)
		songInfoChan <- songInfo // Send received info to the game loop
	}
}

func setupLogging() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	logDir := filepath.Dir(exePath)
	logFilePath := filepath.Join(logDir, "native_app.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	log.Println("---------------------")
	log.Println("Native app started.")
}
