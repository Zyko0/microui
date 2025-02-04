// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package main

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"

	"github.com/ebitengine/microui"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	logBuf       string
	logSubmitBuf string
	logUpdated   bool
	bg           = [3]float32{90, 95, 100}
	checks       = [3]bool{true, false, true}
)

func WriteLog(text string) {
	if len(logBuf) > 0 {
		logBuf += "\n"
	}
	logBuf += text
	logUpdated = true
}

func TestWindow(ctx *microui.Context) {
	if ctx.BeginWindow("Demo Window", image.Rect(40, 40, 340, 490)) {
		defer ctx.EndWindow()
		win := ctx.GetCurrentContainer()
		win.Rect.Max.X = win.Rect.Min.X + max(win.Rect.Dx(), 240)
		win.Rect.Max.Y = win.Rect.Min.Y + max(win.Rect.Dy(), 300)

		/* window info */
		if ctx.Header("Window Info") {
			win := ctx.GetCurrentContainer()
			ctx.LayoutRow(2, []int{54, -1}, 0)
			ctx.Label("Position:")
			ctx.Label(fmt.Sprintf("%d, %d", win.Rect.Min.X, win.Rect.Min.Y))
			ctx.Label("Size:")
			ctx.Label(fmt.Sprintf("%d, %d", win.Rect.Dx(), win.Rect.Dy()))
		}

		/* labels + buttons */
		if ctx.HeaderEx("Test Buttons", microui.OptExpanded) != 0 {
			ctx.LayoutRow(3, []int{86, -110, -1}, 0)
			ctx.Label("Test buttons 1:")
			if ctx.Button("Button 1") {
				WriteLog("Pressed button 1")
			}
			if ctx.Button("Button 2") {
				WriteLog("Pressed button 2")
			}
			ctx.Label("Test buttons 2:")
			if ctx.Button("Button 3") {
				WriteLog("Pressed button 3")
			}
			if ctx.Button("Popup") {
				ctx.OpenPopup("Test Popup")
			}
			if ctx.BeginPopup("Test Popup") != 0 {
				ctx.Button("Hello")
				ctx.Button("World")
				ctx.EndPopup()
			}
		}

		/* tree */
		if ctx.HeaderEx("Tree and Text", microui.OptExpanded) != 0 {
			ctx.LayoutRow(2, []int{140, -1}, 0)
			ctx.LayoutBeginColumn()
			if ctx.BeginTreeNode("Test 1") {
				if ctx.BeginTreeNode("Test 1a") {
					ctx.Label("Hello")
					ctx.Label("World")
					ctx.EndTreeNode()
				}
				if ctx.BeginTreeNode("Test 1b") {
					if ctx.Button("Button 1") {
						WriteLog("Pressed button 1")
					}
					if ctx.Button("Button 2") {
						WriteLog("Pressed button 2")
					}
					ctx.EndTreeNode()
				}
				ctx.EndTreeNode()
			}
			if ctx.BeginTreeNode("Test 2") {
				ctx.LayoutRow(2, []int{54, 54}, 0)
				if ctx.Button("Button 3") {
					WriteLog("Pressed button 3")
				}
				if ctx.Button("Button 4") {
					WriteLog("Pressed button 4")
				}
				if ctx.Button("Button 5") {
					WriteLog("Pressed button 5")
				}
				if ctx.Button("Button 6") {
					WriteLog("Pressed button 6")
				}
				ctx.EndTreeNode()
			}
			if ctx.BeginTreeNode("Test 3") {
				ctx.Checkbox("Checkbox 1", &checks[0])
				ctx.Checkbox("Checkbox 2", &checks[1])
				ctx.Checkbox("Checkbox 3", &checks[2])
				ctx.EndTreeNode()
			}
			ctx.LayoutEndColumn()

			ctx.LayoutBeginColumn()
			ctx.LayoutRow(1, []int{-1}, 0)
			ctx.Text("Lorem ipsum dolor sit amet, consectetur adipiscing " +
				"elit. Maecenas lacinia, sem eu lacinia molestie, mi risus faucibus " +
				"ipsum, eu varius magna felis a nulla.")
			ctx.LayoutEndColumn()
		}

		/* background color sliders */
		if ctx.HeaderEx("Background Color", microui.OptExpanded) != 0 {
			ctx.LayoutRow(2, []int{-78, -1}, 74)
			/* sliders */
			ctx.LayoutBeginColumn()
			ctx.LayoutRow(2, []int{46, -1}, 0)
			ctx.Label("Red:")
			ctx.Slider(&bg[0], 0, 255)
			ctx.Label("Green:")
			ctx.Slider(&bg[1], 0, 255)
			ctx.Label("Blue:")
			ctx.Slider(&bg[2], 0, 255)
			ctx.LayoutEndColumn()
			/* color preview */
			r := ctx.LayoutNext()
			ctx.DrawRect(r, color.RGBA{byte(bg[0]), byte(bg[1]), byte(bg[2]), 255})
			clr := fmt.Sprintf("#%02X%02X%02X", int(bg[0]), int(bg[1]), int(bg[2]))
			ctx.DrawControlText(clr, r, microui.ColorText, microui.OptAlignCenter)
		}
	}
}

func LogWindow(ctx *microui.Context) {
	if ctx.BeginWindow("Log Window", image.Rect(350, 40, 650, 240)) {
		defer ctx.EndWindow()
		/* output text panel */
		ctx.LayoutRow(1, []int{-1}, -25)
		ctx.BeginPanel("Log Output")
		panel := ctx.GetCurrentContainer()
		ctx.LayoutRow(1, []int{-1}, -1)
		ctx.Text(logBuf)
		ctx.EndPanel()
		if logUpdated {
			panel.Scroll.Y = panel.ContentSize.Y
			logUpdated = false
		}

		/* input textbox + submit button */
		var submitted bool
		ctx.LayoutRow(2, []int{-70, -1}, 0)
		if ctx.TextBox(&logSubmitBuf)&microui.ResSubmit != 0 {
			ctx.SetFocus(ctx.LastID)
			submitted = true
		}
		if ctx.Button("Submit") {
			submitted = true
		}
		if submitted {
			WriteLog(logSubmitBuf)
			logSubmitBuf = ""
		}
	}
}

func uint8Slider(ctx *microui.Context, fvalue *float32, value *uint8, low, high uint8) int {
	*fvalue = float32(*value)
	ctx.PushID(microui.PtrToBytes(unsafe.Pointer(fvalue)))
	res := ctx.SliderEx(fvalue, float32(low), float32(high), 0, "%.0f", microui.OptAlignCenter)
	*value = uint8(*fvalue)
	ctx.PopID()

	return res
}

var (
	fcolors = [14]struct {
		R, G, B, A float32
	}{}
	colors = []struct {
		Label   string
		ColorID int
	}{
		{"text:", microui.ColorText},
		{"border:", microui.ColorBorder},
		{"windowbg:", microui.ColorWindowBG},
		{"titlebg:", microui.ColorTitleBG},
		{"titletext:", microui.ColorTitleText},
		{"panelbg:", microui.ColorPanelBG},
		{"button:", microui.ColorButton},
		{"buttonhover:", microui.ColorButtonHover},
		{"buttonfocus:", microui.ColorButtonFocus},
		{"base:", microui.ColorBase},
		{"basehover:", microui.ColorBaseHover},
		{"basefocus:", microui.ColorBaseFocus},
		{"scrollbase:", microui.ColorScrollBase},
		{"scrollthumb:", microui.ColorScrollThumb},
	}
)

func StyleWindow(ctx *microui.Context) {
	if ctx.BeginWindow("Style Editor", image.Rect(350, 250, 650, 490)) {
		sw := int(float64(ctx.GetCurrentContainer().Body.Dx()) * 0.14)
		ctx.LayoutRow(6, []int{80, sw, sw, sw, sw, -1}, 0)
		for _, c := range colors {
			ctx.Label(c.Label)
			uint8Slider(ctx, &fcolors[c.ColorID].R, &ctx.Style.Colors[c.ColorID].R, 0, 255)
			uint8Slider(ctx, &fcolors[c.ColorID].G, &ctx.Style.Colors[c.ColorID].G, 0, 255)
			uint8Slider(ctx, &fcolors[c.ColorID].B, &ctx.Style.Colors[c.ColorID].B, 0, 255)
			uint8Slider(ctx, &fcolors[c.ColorID].A, &ctx.Style.Colors[c.ColorID].A, 0, 255)
			ctx.DrawRect(ctx.LayoutNext(), ctx.Style.Colors[c.ColorID])
		}
		ctx.EndWindow()
	}
}

func (g *Game) ProcessFrame() {
	g.ctx.Begin()
	TestWindow(g.ctx)
	LogWindow(g.ctx)
	StyleWindow(g.ctx)
	g.ctx.End()
}
