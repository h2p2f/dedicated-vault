package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
)

func (g *GraphicApp) credentialTab() (*widget.List, *fyne.Container) {

	// declare credentials area's widgets
	crLabel := widget.NewLabel("Credentials details:")
	crMeta := widget.NewLabel("Meta:")
	crMetaEntry := widget.NewEntry()
	crLogin := widget.NewLabel("Login:")
	crLoginEntry := widget.NewEntry()
	crPassword := widget.NewLabel("Password:")
	crPasswordEntry := widget.NewPasswordEntry()
	crUUIDLabel := widget.NewLabel("")
	crUUIDLabel.Hide()

	//get credentials list
	listData, err := g.processor.GetDataByType("cr")
	if err != nil {
		fmt.Println(err)
	}

	// construct credentials list
	credentialsList := widget.NewList(
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
	// add credentials list onSelected event
	credentialsList.OnSelected = func(id widget.ListItemID) {
		crMetaEntry.SetText(listData[id].Meta)
		crLoginEntry.SetText(listData[id].Folder.Credentials.Login)
		crPasswordEntry.SetText(listData[id].Folder.Credentials.Password)
		crUUIDLabel.SetText(listData[id].UUID)
	}

	refresh := func() {
		g.notLoggedIn()
		listData, err = g.processor.GetDataByType("cr")
		if err != nil {
			fmt.Println(err)
		}
		credentialsList.Refresh()
	}

	addButton := widget.NewButton("Add", func() {

		if crMetaEntry.Text == "" || crLoginEntry.Text == "" || crPasswordEntry.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}

		folder := models.Folder{
			Credentials: models.Credentials{
				Login:    crLoginEntry.Text,
				Password: crPasswordEntry.Text,
			},
		}
		saved := models.Data{
			Meta:     crMetaEntry.Text,
			DataType: "cr",
			Folder:   folder,
		}

		err := g.processor.SaveData(saved)
		if err != nil {
			return
		}
		r := refresh
		r()
	})

	refreshButton := widget.NewButton("Refresh", refresh)

	removeButton := widget.NewButton("Remove", func() {
		if crMetaEntry.Text == "" || crLoginEntry.Text == "" || crPasswordEntry.Text == "" || crUUIDLabel.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		folder := models.Folder{
			Credentials: models.Credentials{
				Login:    crLoginEntry.Text,
				Password: crPasswordEntry.Text,
			},
		}
		removed := models.Data{
			UUID:     crUUIDLabel.Text,
			Meta:     crMetaEntry.Text,
			DataType: "cr",
			Folder:   folder,
		}
		err := g.processor.DeleteData(removed)
		if err != nil {
			return
		}
		r := refresh
		r()
	})

	editButton := widget.NewButton("Edit", func() {
		if crMetaEntry.Text == "" || crLoginEntry.Text == "" || crPasswordEntry.Text == "" || crUUIDLabel.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}
		folder := models.Folder{
			Credentials: models.Credentials{
				Login:    crLoginEntry.Text,
				Password: crPasswordEntry.Text,
			},
		}
		edited := models.Data{
			UUID:     crUUIDLabel.Text,
			Meta:     crMetaEntry.Text,
			DataType: "cr",
			Folder:   folder,
		}
		err := g.processor.ChangeData(edited)
		if err != nil {
			return
		}
		r := refresh
		r()
	})

	credentialsDetailsBox := container.NewVBox(
		refreshButton,
		addButton, removeButton, editButton,
		crLabel,
		crMeta, crMetaEntry,
		crLogin, crLoginEntry,
		crPassword, crPasswordEntry,
		crUUIDLabel,
	)
	// End of Credentials tab
	return credentialsList, credentialsDetailsBox
}
