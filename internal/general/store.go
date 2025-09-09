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

func Forming(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	StoreSidebar(user, app, window, content)

	z := canvas.NewText("Формирование заказов:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)
}

func StoreSidebar(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	img := canvas.NewImageFromFile("1.jpg")
	img.Resize(fyne.NewSize(200, 200))
	img.Move(fyne.NewPos(50, 600))
	content.Add(img)

	profileBtn := widget.NewButton("Профиль", func() { Profile(user, app, window, content) })
	profileBtn.Resize(fyne.NewSize(200, 50))
	profileBtn.Move(fyne.NewPos(50, 20))
	content.Add(profileBtn)

	formedBtn := widget.NewButton("Формирование", func() { Forming(user, app, window, content) })
	formedBtn.Resize(fyne.NewSize(200, 50))
	formedBtn.Move(fyne.NewPos(50, 120))
	content.Add(formedBtn)
}
