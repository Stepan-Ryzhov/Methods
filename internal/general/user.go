package general

import (
	//"errors"
	//"fmt"
	"image/color"
	//"strconv"

	models "methodi_razrabotki/internal/models"
	//rep "methodi_razrabotki/internal/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func UserOrders(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Мои заказы:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)
}

func Cart(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Корзина:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)
}

func Catalog(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Каталог:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)
}

func UserSidebar(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	img := canvas.NewImageFromFile("1.jpg")
	img.Resize(fyne.NewSize(200, 200))
	img.Move(fyne.NewPos(50, 600))
	content.Add(img)

	profileBtn := widget.NewButton("Профиль", func() { Profile(user, app, window, content) })
	profileBtn.Resize(fyne.NewSize(200, 50))
	profileBtn.Move(fyne.NewPos(50, 20))
	content.Add(profileBtn)

	cartBtn := widget.NewButton("Корзина", func() { Cart(user, app, window, content) })
	cartBtn.Resize(fyne.NewSize(200, 50))
	cartBtn.Move(fyne.NewPos(50, 120))
	content.Add(cartBtn)

	catalogBtn := widget.NewButton("Каталог", func() { Catalog(user, app, window, content) })
	catalogBtn.Resize(fyne.NewSize(200, 50))
	catalogBtn.Move(fyne.NewPos(50, 220))
	content.Add(catalogBtn)

	orderBtn := widget.NewButton("Мои заказы", func() { UserOrders(user, app, window, content) })
	orderBtn.Resize(fyne.NewSize(200, 50))
	orderBtn.Move(fyne.NewPos(50, 320))
	content.Add(orderBtn)
}
