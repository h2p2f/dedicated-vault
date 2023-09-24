package gui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	GetDataByType(dataType string) ([]models.Data, error)
	FullSync() error
}

type GraphicApp struct {
	processor   Processor
	config      *config.ClientConfig
	guiApp      fyne.App
	mainWindow  fyne.Window
	notLoggedIn func()
	lostData    func()
}

func NewGraphicApp(proc Processor, conf *config.ClientConfig) *GraphicApp {
	return &GraphicApp{
		processor: proc,
		config:    conf,
		guiApp:    app.New(),
	}
}

func (g *GraphicApp) Run(ctx context.Context) {
	guiApp := g.guiApp
	g.mainWindow = guiApp.NewWindow("Dedicated Vault")
	g.mainWindow.Resize(fyne.NewSize(800, 500))
	g.mainWindow.CenterOnScreen()

	g.notLoggedIn = func() {
		if g.config.User == "" {
			dialog.ShowInformation("Error", "You are not logged in", g.mainWindow)
		}
	}
	g.lostData = func() {
		dialog.ShowInformation("Error", "You lost your data", g.mainWindow)
	}

	img := canvas.NewImageFromFile("img/logo.png")
	img.FillMode = canvas.ImageFillOriginal

	userBox := g.settingsTab()
	crList, crBox := g.credentialTab()
	ccList, ccBox := g.creditCardTab()
	txList, txBox := g.textTab()
	biList, biBox := g.binaryTab()

	//_ = settingsBox

	//settingsTabItem := container.NewTabItem("Settings", container.NewGridWithColumns(2,
	//	userBox,
	//))
	settingsTabItem := container.NewTabItem("Settings", container.New(
		layout.NewGridLayoutWithColumns(3),
		container.New(layout.NewCenterLayout()),
		userBox,
		container.New(layout.NewCenterLayout()),
	))

	credentialTabItem := container.NewTabItem("Credentials", container.NewGridWithColumns(2,
		crList,
		crBox,
	))
	creditCardTabItem := container.NewTabItem("Credit Cards", container.NewGridWithColumns(2,
		ccList,
		ccBox,
	))
	textTabItem := container.NewTabItem("Text Notes", container.NewGridWithColumns(2,
		txList,
		txBox,
	))
	binaryTabItem := container.NewTabItem("Binary Files", container.NewGridWithColumns(2,
		biList,
		biBox,
	))

	appTabs := container.NewAppTabs(
		settingsTabItem,
		credentialTabItem,
		creditCardTabItem,
		textTabItem,
		binaryTabItem,
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
	g.mainWindow.SetContent(appTabs)
	g.mainWindow.Show()

	g.mainWindow.SetMaster()
	guiApp.Run()
}
