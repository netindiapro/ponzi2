package view

import (
	"context"
	"image"
	"runtime"
	"time"
	"unicode"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang/glog"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/btmura/ponzi2/internal/gfx"
	"github.com/btmura/ponzi2/internal/matrix"
)

// Get esc from github.com/mjibson/esc. It's used to embed resources into the binary.
//go:generate esc -o bindata.go -pkg view -include ".*(ply|png)" -modtime 1337 -private data

// Application name for the window title.
const appName = "ponzi2"

// Constants used by Run for the "game loop".
const (
	fps        = 120.0
	updateSec  = 1.0 / fps
	minUpdates = 1
	maxUpdates = 1000
)

// acceptedChars are the chars the user can enter for a symbol.
var acceptedChars = map[rune]bool{
	'A': true, 'B': true, 'C': true,
	'D': true, 'E': true, 'F': true,
	'G': true, 'H': true, 'I': true,
	'J': true, 'K': true, 'L': true,
	'M': true, 'N': true, 'O': true,
	'P': true, 'Q': true, 'R': true,
	'S': true, 'T': true, 'U': true,
	'V': true, 'W': true, 'X': true,
	'Y': true, 'Z': true,
}

// Colors used throughout the UI.
var (
	green     = [3]float32{0.25, 1, 0}
	red       = [3]float32{1, 0.3, 0}
	yellow    = [3]float32{1, 1, 0}
	purple    = [3]float32{0.5, 0, 1}
	white     = [3]float32{1, 1, 1}
	gray      = [3]float32{0.15, 0.15, 0.15}
	lightGray = [3]float32{0.35, 0.35, 0.35}
	orange    = [3]float32{1, 0.5, 0}
)

var (
	viewInputSymbolTextRenderer = gfx.NewTextRenderer(goregular.TTF, 48)
	viewInstructionsText        = newCenteredText(gfx.NewTextRenderer(goregular.TTF, 24), "Type in symbol and press ENTER...")
)

func init() {
	// This is needed to arrange that main() runs on main thread for GLFW.
	// See documentation for functions that are only allowed to be called
	// from the main thread.
	runtime.LockOSThread()
}

// The View renders the UI to view and edit the model's stocks that it observes.
type View struct {
	// win is the handle to the GLFW window.
	win *glfw.Window

	// winSize is the current window's size used to measure and draw the UI.
	winSize image.Point

	// chart renders the currently viewed stock.
	chart *Chart

	// chartThumbs renders the stocks in the sidebar.
	chartThumbs []*ChartThumb

	// inputSymbol stores and renders the symbol being entered by the user.
	inputSymbol *centeredText

	// inputSymbolSubmittedCallback is called when a new symbol is entered.
	inputSymbolSubmittedCallback func(symbol string)

	// mousePos is the current global mouse position.
	mousePos image.Point

	// mouseLeftButtonClicked is whether the left mouse button was clicked.
	mouseLeftButtonClicked bool

	// sidebarOffset stores the Y offset change due to scroll wheel events.
	sidebarOffset image.Point
}

// viewContext is passed down the view hierarchy providing drawing hints and
// event information. Meant to be passed around like a Rectangle or Point rather
// than a pointer to avoid mistakes.
type viewContext struct {
	// Bounds is the rectangle with global coords that should be drawn within.
	Bounds image.Rectangle

	// MousePos is the current global mouse position.
	MousePos image.Point

	// MouseLeftButtonClicked is whether the left mouse button was clicked.
	MouseLeftButtonClicked bool

	// Fudge is the position from 0 to 1 between the current and next frame.
	Fudge float32

	// ScheduledCallbacks are callbacks to be called at the end of Render.
	ScheduledCallbacks *[]func()
}

// LeftClickInBounds returns true if the left mouse button was clicked within
// the context's bounds. Doesn't take into account overlapping view parts.
func (vc viewContext) LeftClickInBounds() bool {
	return vc.MouseLeftButtonClicked && vc.MousePos.In(vc.Bounds)
}

// New creates a new View.
func New() *View {
	return &View{
		inputSymbol:                  newCenteredText(viewInputSymbolTextRenderer, "", centeredTextBubble(chartRounding, chartPadding)),
		inputSymbolSubmittedCallback: func(symbol string) {},
	}
}

// Init initializes the View and returns a cleanup function.
func (v *View) Init(ctx context.Context) (cleanup func(), err error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// Set the following hints for Linux compatibility.
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(800, 600, appName, nil, nil)
	if err != nil {
		return nil, err
	}
	v.win = win

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}
	glog.V(2).Infof("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 0)

	if err := gfx.InitProgram(); err != nil {
		return nil, err
	}

	gfx.SetAlpha(1.0)

	// Call the size callback to set the initial viewport.
	w, h := win.GetSize()
	v.handleSizeEvent(w, h)
	win.SetSizeCallback(func(win *glfw.Window, width, height int) {
		v.handleSizeEvent(width, height)
	})

	win.SetCharCallback(func(win *glfw.Window, char rune) {
		v.handleCharEvent(char)
	})

	win.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		v.handleKeyEvent(key, action)
	})

	win.SetCursorPosCallback(func(win *glfw.Window, xpos, ypos float64) {
		v.handleCursorPosEvent(xpos, ypos)
	})

	win.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		v.handleMouseButtonEvent(button, action)
	})

	win.SetScrollCallback(func(win *glfw.Window, xoff, yoff float64) {
		v.handleScrollEvent(yoff)
	})

	return func() { glfw.Terminate() }, nil
}

func (v *View) handleSizeEvent(width, height int) {
	glog.V(2).Infof("width:%d height:%d", width, height)
	defer v.PostEmptyEvent()

	s := image.Pt(width, height)
	if v.winSize == s {
		return
	}

	gl.Viewport(0, 0, int32(s.X), int32(s.Y))

	// Calculate the new ortho projection view matrix.
	fw, fh := float32(s.X), float32(s.Y)
	gfx.SetProjectionViewMatrix(matrix.Ortho(fw, fh, fw /* use width as depth */))

	v.winSize = s
}

func (v *View) handleCharEvent(char rune) {
	glog.V(2).Infof("char:%c", char)
	defer v.PostEmptyEvent()

	char = unicode.ToUpper(char)
	if _, ok := acceptedChars[char]; ok {
		v.inputSymbol.Text += string(char)
	}
}

func (v *View) handleKeyEvent(key glfw.Key, action glfw.Action) {
	glog.V(2).Infof("key:%v action:%v", key, action)
	defer v.PostEmptyEvent()

	if action != glfw.Release {
		return
	}

	switch key {
	case glfw.KeyEscape:
		v.inputSymbol.Text = ""

	case glfw.KeyBackspace:
		if l := len(v.inputSymbol.Text); l > 0 {
			v.inputSymbol.Text = v.inputSymbol.Text[:l-1]
		}

	case glfw.KeyEnter:
		v.inputSymbolSubmittedCallback(v.inputSymbol.Text)
		v.inputSymbol.Text = ""
	}
}

func (v *View) handleCursorPosEvent(x, y float64) {
	glog.V(2).Infof("x:%f y:%f", x, y)
	defer v.PostEmptyEvent()

	// Flip Y-axis since the OpenGL coordinate system makes lower left the origin.
	v.mousePos = image.Pt(int(x), v.winSize.Y-int(y))
}

func (v *View) handleMouseButtonEvent(button glfw.MouseButton, action glfw.Action) {
	glog.V(2).Infof("button:%v action:%v", button, action)
	defer v.PostEmptyEvent()

	if button != glfw.MouseButtonLeft {
		return
	}

	v.mouseLeftButtonClicked = action == glfw.Release
}

func (v *View) handleScrollEvent(yoff float64) {
	glog.V(2).Infof("yoff:%f", yoff)
	defer v.PostEmptyEvent()

	if yoff != -1 && yoff != +1 {
		return
	}

	if len(v.chartThumbs) == 0 {
		return
	}

	m := v.metrics()

	if !v.mousePos.In(m.sidebarRegion) {
		return
	}

	if m.sidebarRect.Dy() < v.winSize.Y {
		return
	}

	// Scroll wheel down: yoff = -1 up: yoff = +1
	off := sidebarScrollAmount.Mul(-int(yoff))
	tmpRect := m.sidebarRect.Add(off)
	if tmpRect.Min.Y > 0 {
		off.Y -= tmpRect.Min.Y
	}
	if tmpRect.Max.Y < v.winSize.Y {
		off.Y += v.winSize.Y - tmpRect.Max.Y
	}

	v.sidebarOffset = v.sidebarOffset.Add(off)
}

// Run runs the "game loop".
func (v *View) Run(preupdate func()) {
start:
	var lag float64
	animating := false
	prevTime := glfw.GetTime()
	for !v.win.ShouldClose() {
		currTime := glfw.GetTime() /* seconds */
		elapsed := currTime - prevTime
		prevTime = currTime
		lag += elapsed

		i := 0
		for ; i < minUpdates || i < maxUpdates && lag >= updateSec; i++ {
			preupdate()
			animating = v.update()
			lag -= updateSec
		}
		if lag < 0 {
			lag = 0
		}

		fudge := float32(lag / updateSec)
		if fudge < 0 || fudge > 1 {
			fudge = 0
		}

		now := time.Now()
		v.render(fudge)
		v.win.SwapBuffers()
		glog.V(2).Infof("updates:%d lag(%f)/updateSec(%f)=fudge(%f) animating:%t render:%v", i, lag, updateSec, fudge, animating, time.Since(now).Seconds())

		glfw.PollEvents()
		if !animating {
			glog.V(2).Info("wait events")
			glfw.WaitEvents()
			goto start
		}
	}
}

func (v *View) update() (animating bool) {
	if v.chart != nil {
		if v.chart.Update() {
			animating = true
		}
	}
	for _, th := range v.chartThumbs {
		if th.Update() {
			animating = true
		}
	}
	return animating
}

func (v *View) render(fudge float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	vc := viewContext{
		Bounds:                 image.Rectangle{image.ZP, v.winSize},
		MousePos:               v.mousePos,
		MouseLeftButtonClicked: v.mouseLeftButtonClicked,
		Fudge:              fudge,
		ScheduledCallbacks: new([]func()),
	}

	ogBnds := vc.Bounds.Inset(viewOuterPadding)

	// Calculate bounds for main area.
	vc.Bounds = ogBnds
	if len(v.chartThumbs) > 0 {
		vc.Bounds.Min.X += viewOuterPadding + viewChartThumbSize.X
	}

	// Render the the main chart or instructions.
	if v.chart != nil {
		v.chart.Render(vc)
	} else {
		viewInstructionsText.Render(vc.Bounds)
	}

	// Render the input symbol over the chart.
	v.inputSymbol.Render(vc.Bounds)

	// Render the sidebar thumbnails.
	if len(v.chartThumbs) != 0 {
		vc.Bounds = image.Rect(
			viewOuterPadding, ogBnds.Max.Y-viewChartThumbSize.Y,
			viewOuterPadding+viewChartThumbSize.X, ogBnds.Max.Y,
		)
		vc.Bounds = vc.Bounds.Add(v.sidebarOffset)
		for _, th := range v.chartThumbs {
			th.Render(vc)
			vc.Bounds = vc.Bounds.Sub(image.Pt(0, viewChartThumbSize.Y+viewOuterPadding))
		}
	}

	// Call any callbacks scheduled by views.
	for _, cb := range *vc.ScheduledCallbacks {
		cb()
	}

	// Reset any flags for the next viewContext.
	v.mouseLeftButtonClicked = false
}

// PostEmptyEvent wakes up the Run loop with an event if it is asleep.
func (v *View) PostEmptyEvent() {
	glfw.PostEmptyEvent()
}

// SetInputSymbolSubmittedCallback sets the callback for when a new symbol is entered.
func (v *View) SetInputSymbolSubmittedCallback(cb func(symbol string)) {
	v.inputSymbolSubmittedCallback = cb
}

// SetChart sets the View's main chart.
func (v *View) SetChart(ch *Chart) {
	defer v.PostEmptyEvent()
	v.chart = ch
}

// AddChartThumb adds the ChartThumbnail to the side bar.
func (v *View) AddChartThumb(th *ChartThumb) {
	defer v.PostEmptyEvent()
	v.chartThumbs = append(v.chartThumbs, th)
}

// RemoveChartThumb removes the ChartThumbnail from the side bar.
func (v *View) RemoveChartThumb(th *ChartThumb) {
	defer v.PostEmptyEvent()
	for i, thumb := range v.chartThumbs {
		if thumb == th {
			v.chartThumbs = append(v.chartThumbs[:i], v.chartThumbs[i+1:]...)
			break
		}
	}
}
