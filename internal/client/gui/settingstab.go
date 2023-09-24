package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GraphicApp) settingsTab() *fyne.Container {

	LoginLabel := widget.NewLabel("Login")
	login := widget.NewEntry()
	userLabel := widget.NewLabel("User")
	if g.config.User == "" {
		userLabel.Hide()
	}
	fullSyncButton := widget.NewButton("Full sync", func() {
		err := g.processor.FullSync()
		if err != nil {
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
	loginButton := widget.NewButton("Login", func() {
		if login.Text == "" || password.Text == "" || passphrase.Text == "" {
			return
		}
		err := g.processor.LoginUser(login.Text, password.Text, passphrase.Text)
		if err != nil {
			return
		}
		userLabel.SetText("User logged in: " + login.Text)
		userLabel.Show()
		fullSyncButton.Show()
		LoginLabel.Hide()
		login.Hide()
		passwordLabel.Hide()
		password.Hide()
		passphraseLabel.Hide()
		passphrase.Hide()

	})

	registerButton := widget.NewButton("Register", func() {
		if login.Text == "" || password.Text == "" || passphrase.Text == "" {
			return
		}
		err := g.processor.CreateUser(login.Text, password.Text, passphrase.Text)
		if err != nil {
			return
		}
		userLabel.SetText("User registered and logged in: " + login.Text)
		userLabel.Show()
		fullSyncButton.Show()
		LoginLabel.Hide()
		login.Hide()
		passwordLabel.Hide()
		password.Hide()
		passphraseLabel.Hide()
		passphrase.Hide()
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
