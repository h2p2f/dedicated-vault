package gui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
	"time"
)

type Processor interface {
	CreateUser(userName, password, passphrase string) error
	LoginUser(userName, password, passphrase string) error
	ChangePassword(userName, password, newPassword string) error
	SaveData(data models.Data) error
	ChangeData(data models.Data) error
	DeleteData(data models.Data) error
	GetData(uuid string) (*models.Data, error)
	FullSync() error
}

type GUIApp struct {
	processor Processor
	config    *config.ClientConfig
}

func NewGUIApp(proc Processor, conf *config.ClientConfig) *GUIApp {
	return &GUIApp{
		processor: proc,
		config:    conf,
	}
}

func (g *GUIApp) Run(ctx context.Context) {
	guiApp := app.New()
	mainWindow := guiApp.NewWindow("Dedicated Vault")
	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.CenterOnScreen()

	img := canvas.NewImageFromFile("img/logo.png")
	img.FillMode = canvas.ImageFillOriginal
	LoginLabel := widget.NewLabel("Login")
	login := widget.NewEntry()

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

	})
	registerButton := widget.NewButton("Register", func() {
		if login.Text == "" || password.Text == "" || passphrase.Text == "" {
			return
		}
		err := g.processor.CreateUser(login.Text, password.Text, passphrase.Text)
		if err != nil {
			return
		}

	})
	cancelButton := widget.NewButton("Cancel", func() {
		mainWindow.Close()
	})

	authButtonContainer := container.New(
		layout.NewCenterLayout(),
		container.NewHBox(loginButton, registerButton, cancelButton),
	)

	registerContainer := container.New(
		layout.NewCenterLayout(),
		container.NewVBox(
			LoginLabel, login,
			passwordLabel, password,
			passphraseLabel, passphrase,
			authButtonContainer),
	)

	addButton := widget.NewButton("Add", func() {})
	removeButton := widget.NewButton("Remove", func() {})
	editButton := widget.NewButton("Edit", func() {})
	buttonBox := container.NewHBox(addButton, removeButton, editButton)

	appTabs := container.NewAppTabs(
		container.NewTabItem("Settings", container.NewVBox(
			widget.NewLabel("Settings"),
			registerContainer,
		)),
		container.NewTabItem("Credentials", container.NewVBox(
			widget.NewLabel("Credentials"),
			buttonBox,
		)),
		container.NewTabItem("Credit Cards", container.NewVBox(
			widget.NewLabel("Credit Cards"),
			buttonBox,
		)),
		container.NewTabItem("Text Notes", container.NewVBox(
			widget.NewLabel("Notes"),
			buttonBox,
		)),
		container.NewTabItem("Binary Files", container.NewVBox(
			widget.NewLabel("Files"),
			buttonBox,
		)),
	)
	appTabs.SetTabLocation(container.TabLocationLeading)

	if drv, ok := guiApp.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.Resize(fyne.NewSize(800, 600))
		splash.SetContent(
			container.NewHBox(widget.NewLabel("Dedicated Vault"),
				container.NewVBox(
					widget.NewLabel("created by github.com/h2p2f"),
					img,
					widget.NewLabel("Version 0.0.1"),
				),
			))
		splash.CenterOnScreen()
		splash.Show()
		go func() {
			time.Sleep(4 * time.Second)
			splash.Close()
		}()
	}
	mainWindow.SetContent(appTabs)
	mainWindow.Show()

	mainWindow.SetMaster()
	guiApp.Run()
}
