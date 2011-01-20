// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

import (
	"walk/drawing"
	"walk/gui"
)

type MainWindow struct {
	*gui.MainWindow
	tabWidget    *gui.TabWidget
	prevFilePath string
}

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func (mw *MainWindow) openImage() {
	dlg := &gui.FileDialog{}

	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.Title = "Select an Image"

	ok, err := dlg.ShowOpen(mw)
	panicIfErr(err)
	if !ok {
		return
	}

	mw.prevFilePath = dlg.FilePath

	img, err := drawing.NewImageFromFile(dlg.FilePath)
	panicIfErr(err)

	page, err := gui.NewTabPage()
	panicIfErr(err)
	panicIfErr(page.SetText(path.Base(strings.Replace(dlg.FilePath, "\\", "/", -1))))
	panicIfErr(page.SetLayout(gui.NewHBoxLayout()))

	imageView, err := gui.NewImageView(page)
	panicIfErr(err)

	panicIfErr(imageView.SetImage(img))

	_, err = mw.tabWidget.Pages().Add(page)
	panicIfErr(err)

	panicIfErr(mw.tabWidget.SetCurrentPage(page))
}

func runMainWindow() (int, os.Error) {
	mainWnd, err := gui.NewMainWindow()
	panicIfErr(err)
	defer mainWnd.Dispose()

	mw := &MainWindow{MainWindow: mainWnd}
	mw.ClientArea().SetLayout(gui.NewVBoxLayout())
	panicIfErr(mw.SetText("Walk Image Viewer Example"))

	mw.tabWidget, err = gui.NewTabWidget(mw.ClientArea())
	panicIfErr(err)

	imageList, err := gui.NewImageList(drawing.Size{16, 16}, 0)
	panicIfErr(err)
	mw.ToolBar().SetImageList(imageList)

	fileMenu, err := gui.NewMenu()
	panicIfErr(err)
	_, fileMenuAction, err := mw.Menu().Actions().AddMenu(fileMenu)
	panicIfErr(err)
	fileMenuAction.SetText("File")

	openBmp, err := drawing.NewBitmapFromFile("img/open.png")
	panicIfErr(err)

	openAction := gui.NewAction()
	openAction.SetImage(openBmp)
	openAction.SetText("Open")
	openAction.Triggered().Subscribe(func(args *gui.EventArgs) { mw.openImage() })
	fileMenu.Actions().Add(openAction)
	mw.ToolBar().Actions().Add(openAction)

	exitAction := gui.NewAction()
	exitAction.SetText("Exit")
	exitAction.Triggered().Subscribe(func(args *gui.EventArgs) { gui.Exit(0) })
	fileMenu.Actions().Add(exitAction)

	helpMenu, err := gui.NewMenu()
	panicIfErr(err)
	_, helpMenuAction, err := mw.Menu().Actions().AddMenu(helpMenu)
	panicIfErr(err)
	helpMenuAction.SetText("Help")

	aboutAction := gui.NewAction()
	aboutAction.SetText("About")
	aboutAction.Triggered().Subscribe(func(args *gui.EventArgs) {
		gui.MsgBox(mw, "About", "Walk Image Viewer Example", gui.MsgBoxOK|gui.MsgBoxIconInformation)
	})
	helpMenu.Actions().Add(aboutAction)

	panicIfErr(mw.SetMinSize(drawing.Size{320, 240}))
	panicIfErr(mw.SetSize(drawing.Size{800, 600}))
	mw.Show()

	return mw.RunMessageLoop()
}

func main() {
	runtime.LockOSThread()

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Error:", x)
		}
	}()

	exitCode, err := runMainWindow()
	panicIfErr(err)
	os.Exit(exitCode)
}
