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
	mainWindow.Resize(fyne.NewSize(800, 400))
	mainWindow.CenterOnScreen()

	img := canvas.NewImageFromFile("img/logo.png")
	img.FillMode = canvas.ImageFillOriginal

	//Settings tab
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
	// End of Settings tab

	// Credentials tab
	credentialsList := widget.NewList(
		func() int {
			return 20
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText("template")
		},
	)
	crLabel := widget.NewLabel("Credentials details:")
	crMeta := widget.NewLabel("Meta:")
	crMetaDetails := widget.NewLabel("")
	crLogin := widget.NewLabel("Login:")
	crLoginDetails := widget.NewLabel("")
	crPassword := widget.NewLabel("Password:")
	crPasswordDetails := widget.NewLabel("")
	addButton := widget.NewButton("Add", func() {})

	removeButton := widget.NewButton("Remove", func() {})

	editButton := widget.NewButton("Edit", func() {})

	credentialsDetailsBox := container.NewVBox(
		addButton, removeButton, editButton,
		crLabel,
		crMeta, crMetaDetails,
		crLogin, crLoginDetails,
		crPassword, crPasswordDetails,
	)
	// End of Credentials tab

	// Credit Cards tab
	creditCardList := widget.NewList(
		func() int {
			return 20
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText("template")
		},
	)

	ccLabel := widget.NewLabel("Credit Card details:")
	ccMeta := widget.NewLabel("Meta:")
	ccMetaDetails := widget.NewLabel("")
	ccNumber := widget.NewLabel("Number:")
	ccNumberDetails := widget.NewLabel("")
	ccNameOnCard := widget.NewLabel("Name on card:")
	ccNameOnCardDetails := widget.NewLabel("")
	ccExpireDate := widget.NewLabel("Expire date:")
	ccExpireDateDetails := widget.NewLabel("")

	ccAddButton := widget.NewButton("Add", func() {})
	ccRemoveButton := widget.NewButton("Remove", func() {})
	ccEditButton := widget.NewButton("Edit", func() {})

	creditCardDetailsBox := container.NewVBox(
		ccAddButton, ccRemoveButton, ccEditButton,
		ccLabel,
		ccMeta, ccMetaDetails,
		ccNumber, ccNumberDetails,
		ccNameOnCard, ccNameOnCardDetails,
		ccExpireDate, ccExpireDateDetails,
	)
	// End of Credit Cards tab

	// Text Notes tab
	textList := widget.NewList(
		func() int {
			return 20
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText("template")
		},
	)

	textLabel := widget.NewLabel("Text Note details:")
	textMeta := widget.NewLabel("Meta:")
	textMetaDetails := widget.NewLabel("")
	textText := widget.NewLabel("Text:")
	textTextDetails := widget.NewRichTextWithText("")
	textTextDetails.Resize(fyne.NewSize(200, 200))
	textAddButton := widget.NewButton("Add", func() {})
	textRemoveButton := widget.NewButton("Remove", func() {})
	textEditButton := widget.NewButton("Edit", func() {})
	textDetailsBox := container.NewVBox(
		textAddButton, textRemoveButton, textEditButton,
		textLabel,
		textMeta, textMetaDetails,
		textText, textTextDetails,
	)
	// End of Text Notes tab

	// Binary Files tab
	binaryList := widget.NewList(
		func() int {
			return 20
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText("template")
		},
	)

	binaryLabel := widget.NewLabel("Binary File details:")
	binaryMeta := widget.NewLabel("Meta:")
	binaryMetaDetails := widget.NewLabel("")
	binaryName := widget.NewLabel("Name:")
	binaryNameDetails := widget.NewLabel("")
	binaryAddButton := widget.NewButton("Add", func() {})
	binaryRemoveButton := widget.NewButton("Remove", func() {})
	binaryEditButton := widget.NewButton("Edit", func() {})
	binarySaveButton := widget.NewButton("Save to disk", func() {})

	binaryDetailsBox := container.NewVBox(
		binaryAddButton, binaryRemoveButton, binaryEditButton, binarySaveButton,
		binaryLabel,
		binaryMeta, binaryMetaDetails,
		binaryName, binaryNameDetails,
	)

	// End of Binary Files tab

	appTabs := container.NewAppTabs(
		container.NewTabItem("Settings", container.NewVBox(
			widget.NewLabel("Settings"),
			registerContainer,
		)),
		container.NewTabItem("Credentials", container.NewGridWithColumns(2,
			credentialsList,
			credentialsDetailsBox,
		)),
		container.NewTabItem("Credit Cards", container.NewGridWithColumns(2,
			creditCardList,
			creditCardDetailsBox,
		)),
		container.NewTabItem("Text Notes", container.NewGridWithColumns(2,
			textList,
			textDetailsBox,
		)),
		container.NewTabItem("Binary Files", container.NewGridWithColumns(2,
			binaryList,
			binaryDetailsBox,
		)),
	)
	appTabs.SetTabLocation(container.TabLocationLeading)

	if drv, ok := guiApp.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.Resize(fyne.NewSize(800, 400))
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
