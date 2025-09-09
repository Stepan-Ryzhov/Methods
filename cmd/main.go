package main

import (
	"image/color"

	db "methodi_razrabotki/internal/database"
	lg "methodi_razrabotki/internal/login"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct {
	defaultTheme fyne.Theme
	variant      fyne.ThemeVariant
	buttonColor  color.Color
	textColor    color.Color
}

func (m MyTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameButton:
		return m.buttonColor
	case theme.ColorNameForeground:
		return m.textColor
	default:
		return m.defaultTheme.Color(name, m.variant)
	}
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return m.defaultTheme.Font(style)
}

func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return m.defaultTheme.Icon(name)
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return m.defaultTheme.Size(name)
}

func main() {
	db.Init()
	myApp := app.New()
	myWindow := myApp.NewWindow("Интернет магазин электроники Axiom Technology")

	buttonColor := color.NRGBA{R: 204, G: 204, B: 204, A: 255}
	textColor := color.Black

	myTheme := MyTheme{
		defaultTheme: theme.DefaultTheme(),
		variant:      theme.VariantLight,
		buttonColor:  buttonColor,
		textColor:    textColor,
	}

	myApp.Settings().SetTheme(myTheme)

	content := container.NewWithoutLayout()
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.SetFixedSize(true)

	lg.Start(myApp, myWindow, content)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
