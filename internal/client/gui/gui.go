// Package: gui
// in this file we have main logic for gui
package gui

import (
	"context"
	"time"

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
)

// Processor is an interface for processing data
type Processor interface {
	CreateUser(ctx context.Context, userName, password, passphrase string) error
	LoginUser(ctx context.Context, userName, password, passphrase string) error
	ChangePassword(ctx context.Context, userName, password, newPassword string) error
	SaveData(ctx context.Context, data models.Data) error
	ChangeData(ctx context.Context, data models.Data) error
	DeleteData(ctx context.Context, data models.Data) error
	GetDataByType(dataType string) ([]models.Data, error)
	FullSync(ctx context.Context) error
}

// Updater is an interface for updating data
type Updater interface {
	FullSync(ctx context.Context, p Processor) error
}

// GraphicApp is a struct for gui application
type GraphicApp struct {
	processor   Processor
	updater     Updater
	config      *config.ClientConfig
	guiApp      fyne.App
	mainWindow  fyne.Window
	notLoggedIn func()
	lostData    func()
	dialogErr   func(err error)
}

// NewGraphicApp creates a new GraphicApp
func NewGraphicApp(proc Processor, conf *config.ClientConfig) *GraphicApp {
	return &GraphicApp{
		processor: proc,
		config:    conf,
		guiApp:    app.New(),
	}
}

// Run launches the main gui logic
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

	g.dialogErr = func(err error) {
		dialog.ShowInformation("Error", err.Error(), g.mainWindow)
	}

	img := canvas.NewImageFromFile("img/logo.png")
	img.FillMode = canvas.ImageFillOriginal

	userBox := g.settingsTab(ctx)
	crList, crBox := g.credentialTab(ctx)
	ccList, ccBox := g.creditCardTab(ctx)
	txList, txBox := g.textTab(ctx)
	biList, biBox := g.binaryTab(ctx)

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
					widget.NewLabel("Version:"+g.config.Version),
					widget.NewLabel("Build date: "+g.config.BuildDate),
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
