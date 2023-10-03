package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/h2p2f/dedicated-vault/internal/client/models"
)

/*
strongly need to refactor this code
it's too long, hard to read and understand
but it's not simple to do - a lot of dependencies of GUI's struct, so it's hard to split it
one more thing - the fyne library has a lot of callbacks, so it's hard to split it too
one of the solutions is to use some kind of dependency injection in the future
*/

// textTab - function for creating text tab
func (g *GraphicApp) textTab(ctx context.Context) (*widget.List, *fyne.Container) {
	var err error
	// declare text area's widgets
	textLabel := widget.NewLabel("Text details:")
	textMeta := widget.NewLabel("Meta:")
	textMetaEntry := widget.NewEntry()
	textContent := widget.NewLabel("Content:")
	textContentEntry := widget.NewMultiLineEntry()
	textContentEntry.SetMinRowsVisible(13)
	textUUIDLabel := widget.NewLabel("")
	textUUIDLabel.Hide()

	//get text list
	//error not handled because it called on startup
	//user may not be logged in
	listData, _ := g.processor.GetDataByType("tx")

	// construct text list
	textList := widget.NewList(
		func() int {
			return len(listData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(listData[i].Meta)
		},
	)
	// add text list onSelected event
	textList.OnSelected = func(id widget.ListItemID) {
		textMetaEntry.SetText(listData[id].Meta)
		textContentEntry.SetText(listData[id].Folder.Text.Text)
		textUUIDLabel.SetText(listData[id].UUID)
	}

	refresh := func() {
		g.notLoggedIn()
		listData, err = g.processor.GetDataByType("tx")
		if err != nil {
			g.dialogErr(err)
			return
		}
		textList.Refresh()
	}

	addButton := widget.NewButton("Add", func() {
		if textMetaEntry.Text == "" || textContentEntry.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		data := models.Data{
			Meta:     textMetaEntry.Text,
			DataType: "tx",
			Folder:   models.Folder{Text: models.TextData{Text: textContentEntry.Text}},
		}
		err := g.processor.SaveData(ctx, data)
		if err != nil {
			g.dialogErr(err)
			return
		}
		r := refresh
		r()
	})

	editButton := widget.NewButton("Edit", func() {
		if textMetaEntry.Text == "" || textContentEntry.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		data := models.Data{
			UUID:   textUUIDLabel.Text,
			Meta:   textMetaEntry.Text,
			Folder: models.Folder{Text: models.TextData{Text: textContentEntry.Text}},
		}
		err := g.processor.ChangeData(ctx, data)
		if err != nil {
			g.dialogErr(err)
			return
		}
		r := refresh
		r()
	})

	deleteButton := widget.NewButton("Delete", func() {
		if textUUIDLabel.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		data := models.Data{
			UUID:   textUUIDLabel.Text,
			Meta:   textMetaEntry.Text,
			Folder: models.Folder{Text: models.TextData{Text: textContentEntry.Text}},
		}
		err := g.processor.DeleteData(ctx, data)
		if err != nil {
			g.dialogErr(err)
			return
		}
		r := refresh
		r()
	})

	refreshButton := widget.NewButton("Refresh", refresh)

	// construct text area
	textDetailsBox := container.NewVBox(
		refreshButton,
		addButton, editButton, deleteButton,
		textLabel,
		textMeta, textMetaEntry,
		textContent, textContentEntry,
		textUUIDLabel,
	)

	return textList, textDetailsBox
}
