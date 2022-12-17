package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type ClickImage struct {
	widget.BaseWidget
	Image    *canvas.Image
	OnTapped func() `json:"-"`
}

func NewClickImage() *ClickImage {
	clickImage := &ClickImage{}
	clickImage.ExtendBaseWidget(clickImage)
	return clickImage
}

var _ fyne.WidgetRenderer = (*ClickImageRenderer)(nil)

type ClickImageRenderer struct {
	fyne.WidgetRenderer
	image *canvas.Image
	c     *ClickImage
}

func (i *ClickImageRenderer) MinSize() fyne.Size {
	t := float32(100)
	return fyne.NewSize(t, t)
}

func (r *ClickImageRenderer) Refresh() {
	//r.image.Image = r.c.Image
	canvas.Refresh(r.c)
}

func (i *ClickImage) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)

	r := &ClickImageRenderer{
		WidgetRenderer: widget.NewSimpleRenderer(i.Image),
		image:          i.Image,
		c:              i,
	}

	return r
}

func (t *ClickImage) Tapped(e *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped()
	}
}

func (t *ClickImage) TappedSecondary(_ *fyne.PointEvent) {
}
