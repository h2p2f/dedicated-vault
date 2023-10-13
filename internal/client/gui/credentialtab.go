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

// credentialTab - function for creating credential tab
func (g *GraphicApp) credentialTab(ctx context.Context) (*widget.List, *fyne.Container) {
	var err error

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
	//error not handled because it called on startup
	//user may not be logged in
	listData, _ := g.processor.GetDataByType("cr")

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
			g.dialogErr(err)
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

		err := g.processor.SaveData(ctx, saved)
		if err != nil {
			g.dialogErr(err)
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
		err := g.processor.DeleteData(ctx, removed)
		if err != nil {
			g.dialogErr(err)
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
		err := g.processor.ChangeData(ctx, edited)
		if err != nil {
			g.dialogErr(err)
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
