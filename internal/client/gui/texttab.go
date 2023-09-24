package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
)

func (g *GraphicApp) textTab() (*widget.List, *fyne.Container) {

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
	listData, err := g.processor.GetDataByType("tx")
	if err != nil {
		fmt.Println(err)
	}

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
		err := g.processor.SaveData(data)
		if err != nil {
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
		err := g.processor.ChangeData(data)
		if err != nil {
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
		err := g.processor.DeleteData(data)
		if err != nil {
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
