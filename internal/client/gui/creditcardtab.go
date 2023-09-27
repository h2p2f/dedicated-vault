package gui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
)

func (g *GraphicApp) creditCardTab(ctx context.Context) (*widget.List, *fyne.Container) {
	var err error
	// declare credit card area's widgets
	ccLabel := widget.NewLabel("Credit card details:")
	ccMeta := widget.NewLabel("Meta:")
	ccMetaEntry := widget.NewEntry()
	ccNumber := widget.NewLabel("Number:")
	ccNumberEntry := widget.NewEntry()
	ccOwner := widget.NewLabel("Owner:")
	ccOwnerEntry := widget.NewEntry()
	ccExpire := widget.NewLabel("Expire:")
	ccExpireEntry := widget.NewEntry()
	ccCVV := widget.NewLabel("CVV:")
	ccCVVEntry := widget.NewPasswordEntry()
	ccUUIDLabel := widget.NewLabel("")
	ccUUIDLabel.Hide()

	//get credit card list
	//error not handled because it called on startup
	//user may not be logged in
	listData, _ := g.processor.GetDataByType("cc")

	// construct credit card list
	creditCardList := widget.NewList(
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

	creditCardList.OnSelected = func(id widget.ListItemID) {
		ccMetaEntry.SetText(listData[id].Meta)
		ccNumberEntry.SetText(listData[id].Folder.Card.Number)
		ccOwnerEntry.SetText(listData[id].Folder.Card.NameOnCard)
		ccExpireEntry.SetText(listData[id].Folder.Card.ExpireDate)
		ccCVVEntry.SetText(listData[id].Folder.Card.CVV)
		ccUUIDLabel.SetText(listData[id].UUID)
	}

	refresh := func() {
		g.notLoggedIn()
		listData, err = g.processor.GetDataByType("cc")
		if err != nil {
			g.dialogErr(err)
			return
		}
		creditCardList.Refresh()
	}

	addButton := widget.NewButton("Add", func() {
		if ccMetaEntry.Text == "" ||
			ccNumberEntry.Text == "" ||
			ccOwnerEntry.Text == "" ||
			ccExpireEntry.Text == "" ||
			ccCVVEntry.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}

		folder := models.Folder{
			Card: models.CreditCard{
				Number:     ccNumberEntry.Text,
				NameOnCard: ccOwnerEntry.Text,
				ExpireDate: ccExpireEntry.Text,
				CVV:        ccCVVEntry.Text,
			},
		}

		saved := models.Data{
			Meta:     ccMetaEntry.Text,
			DataType: "cc",
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

	editButton := widget.NewButton("Edit", func() {
		if ccMetaEntry.Text == "" ||
			ccNumberEntry.Text == "" ||
			ccOwnerEntry.Text == "" ||
			ccExpireEntry.Text == "" ||
			ccCVVEntry.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}

		folder := models.Folder{
			Card: models.CreditCard{
				Number:     ccNumberEntry.Text,
				NameOnCard: ccOwnerEntry.Text,
				ExpireDate: ccExpireEntry.Text,
				CVV:        ccCVVEntry.Text,
			},
		}

		saved := models.Data{
			Meta:     ccMetaEntry.Text,
			DataType: "cc",
			Folder:   folder,
			UUID:     ccUUIDLabel.Text,
		}
		err := g.processor.ChangeData(ctx, saved)
		if err != nil {
			g.dialogErr(err)
			return
		}
		r := refresh
		r()
	})

	removeButton := widget.NewButton("Remove", func() {
		if ccMetaEntry.Text == "" ||
			ccNumberEntry.Text == "" ||
			ccOwnerEntry.Text == "" ||
			ccExpireEntry.Text == "" ||
			ccCVVEntry.Text == "" {
			g.lostData()
		}
		if g.config.User == "" || g.config.Token == "" {
			g.notLoggedIn()
		}

		folder := models.Folder{
			Card: models.CreditCard{
				Number:     ccNumberEntry.Text,
				NameOnCard: ccOwnerEntry.Text,
				ExpireDate: ccExpireEntry.Text,
				CVV:        ccCVVEntry.Text,
			},
		}

		saved := models.Data{
			Meta:     ccMetaEntry.Text,
			DataType: "cc",
			Folder:   folder,
			UUID:     ccUUIDLabel.Text,
		}
		err := g.processor.DeleteData(ctx, saved)
		if err != nil {
			g.dialogErr(err)
			return
		}
		r := refresh
		r()
	})

	refreshButton := widget.NewButton("Refresh", refresh)

	// construct credit card area
	creditCardDetailsBox := container.NewVBox(
		refreshButton,
		addButton, editButton, removeButton,
		ccLabel,
		ccMeta, ccMetaEntry,
		ccNumber, ccNumberEntry,
		ccOwner, ccOwnerEntry,
		ccExpire, ccExpireEntry,
		ccCVV, ccCVVEntry,
		ccUUIDLabel,
	)

	return creditCardList, creditCardDetailsBox

}
