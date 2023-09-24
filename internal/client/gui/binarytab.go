package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
	"io"
)

func (g *GraphicApp) binaryTab() (*widget.List, *fyne.Container) {
	var binaryData []byte

	// declare binary area's widgets
	binaryLabel := widget.NewLabel("Binary details:")
	binaryMeta := widget.NewLabel("Meta:")
	binaryMetaEntry := widget.NewEntry()
	binaryName := widget.NewLabel("Name:")
	binaryNameEntry := widget.NewEntry()
	loadButton := widget.NewButton("Load from disk", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			binaryNameEntry.SetText(reader.URI().Name())
			binaryData, err = io.ReadAll(reader)
		}, g.mainWindow)
	})
	saveButton := widget.NewButton("Save to disk", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			_, err = writer.Write(binaryData)
			binaryNameEntry.SetText(writer.URI().Name())
		}, g.mainWindow)

	})
	binaryUUIDLabel := widget.NewLabel("")
	binaryUUIDLabel.Hide()

	//get binary list
	listData, err := g.processor.GetDataByType("bi")
	if err != nil {
		fmt.Println(err)
	}

	// construct binary list
	binaryList := widget.NewList(
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
	// add binary list onSelected event
	binaryList.OnSelected = func(id widget.ListItemID) {
		binaryMetaEntry.SetText(listData[id].Meta)
		binaryNameEntry.SetText(listData[id].Folder.Binary.Name)
		binaryUUIDLabel.SetText(listData[id].UUID)
		binaryData = listData[id].Folder.Binary.Data
	}

	refresh := func() {
		g.notLoggedIn()
		listData, err = g.processor.GetDataByType("bi")
		if err != nil {
			fmt.Println(err)
		}
		binaryList.Refresh()
	}

	addButton := widget.NewButton("Add", func() {
		if binaryMetaEntry.Text == "" || binaryNameEntry.Text == "" || binaryData == nil {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		folder := models.Folder{
			Binary: models.BinaryData{
				Name: binaryNameEntry.Text,
				Data: binaryData,
			},
		}
		data := models.Data{
			Meta:     binaryMetaEntry.Text,
			DataType: "bi",
			Folder:   folder,
		}
		err := g.processor.SaveData(data)
		if err != nil {
			return
		}
		r := refresh
		r()
	})

	editButton := widget.NewButton("Edit", func() {
		if binaryMetaEntry.Text == "" || binaryNameEntry.Text == "" || binaryData == nil {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		data := models.Data{
			UUID:   binaryUUIDLabel.Text,
			Meta:   binaryMetaEntry.Text,
			Folder: models.Folder{Binary: models.BinaryData{Name: binaryNameEntry.Text, Data: binaryData}},
		}
		err := g.processor.SaveData(data)
		if err != nil {
			return
		}
		r := refresh
		r()
	})

	removeButton := widget.NewButton("Remove", func() {
		if binaryUUIDLabel.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		data := models.Data{
			UUID:   binaryUUIDLabel.Text,
			Meta:   binaryMetaEntry.Text,
			Folder: models.Folder{Binary: models.BinaryData{Name: binaryNameEntry.Text, Data: binaryData}},
		}
		err := g.processor.DeleteData(data)
		if err != nil {
			return
		}
		r := refresh
		r()
	})

	refreshButton := widget.NewButton("Refresh", refresh)
	// construct binary area

	binaryDetailBox := container.NewVBox(
		refreshButton,
		addButton, editButton, removeButton,
		binaryLabel,
		binaryMeta, binaryMetaEntry,
		binaryName, binaryNameEntry,
		loadButton, saveButton,
	)

	return binaryList, binaryDetailBox

}
