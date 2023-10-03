package gui

import (
	"context"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
strongly need to refactor this code
it's too long, hard to read and understand
but it's not simple to do - a lot of dependencies of GUI's struct, so it's hard to split it
one more thing - the fyne library has a lot of callbacks, so it's hard to split it too
one of the solutions is to use some kind of dependency injection in the future
*/

// settingsTab - function for creating settings tab
func (g *GraphicApp) settingsTab(ctx context.Context) *fyne.Container {

	LoginLabel := widget.NewLabel("Login")
	login := widget.NewEntry()
	userLabel := widget.NewLabel("User")
	if g.config.User == "" {
		userLabel.Hide()
	}
	fullSyncButton := widget.NewButton("Full sync", func() {
		err := g.processor.FullSync(ctx)
		if err != nil {
			g.dialogErr(err)
			return
		}
	})
	if g.config.User == "" {
		fullSyncButton.Hide()
	}
	passwordLabel := widget.NewLabel("Password")
	password := widget.NewPasswordEntry()
	passphraseLabel := widget.NewLabel("Passphrase")
	passphrase := widget.NewPasswordEntry()

	hideAndShow := func(s string) {
		userLabel.SetText("User logged in: " + s)
		userLabel.Show()
		fullSyncButton.Show()
		LoginLabel.Hide()
		login.Hide()
		passwordLabel.Hide()
		password.Hide()
		passphraseLabel.Hide()
		passphrase.Hide()
	}

	loginButton := widget.NewButton("Login", func() {
		if login.Text == "" || password.Text == "" || passphrase.Text == "" {
			g.dialogErr(errors.New("empty fields"))
			return
		}
		err := g.processor.LoginUser(ctx, login.Text, password.Text, passphrase.Text)
		if err != nil {
			g.dialogErr(err)
			return
		}
		hideAndShow(login.Text)

	})

	registerButton := widget.NewButton("Register", func() {
		if login.Text == "" || password.Text == "" || passphrase.Text == "" {
			g.dialogErr(errors.New("empty fields"))
			return
		}
		err := g.processor.CreateUser(ctx, login.Text, password.Text, passphrase.Text)
		if err != nil {
			g.dialogErr(err)
			return
		}
		hideAndShow(login.Text)
	})
	exitButton := widget.NewButton("Exit", func() {
		g.mainWindow.Close()
	})

	settingsContainer := container.NewVBox(
		userLabel,
		LoginLabel, login,
		passwordLabel, password,
		passphraseLabel, passphrase,
		loginButton, registerButton,
		fullSyncButton,
		exitButton,
	)

	return settingsContainer

}
