package ungo

type LabeledBitmap struct {
	bitmap *Bitmap
	labels []string
}

func NewLabeledBitmap(labels []string) *LabeledBitmap {
	size := uint(1)
	for len(labels)/8 >= int(size) {
		size += 1
	}

	return &LabeledBitmap{
		bitmap: NewBitmap(size),
		labels: labels,
	}
}

func (b *LabeledBitmap) GetBitmap() *Bitmap {
	return b.bitmap
}

func (b *LabeledBitmap) Test(bit uint) bool {
	return b.bitmap.Test(bit)
}

func (b *LabeledBitmap) Set(bit uint) {
	b.bitmap.Set(bit)
}

func (b *LabeledBitmap) Clear(bit uint) {
	b.bitmap.Clear(bit)
}

func (b *LabeledBitmap) Size() uint {
	return b.bitmap.Size()
}

func (b *LabeledBitmap) Count() uint {
	return b.bitmap.Count()
}

func (b *LabeledBitmap) Label(bit uint) string {
	return b.labels[bit]
}

func (b *LabeledBitmap) Get(label string) bool {
	bit := uint(0)
	for _, l := range b.labels {
		if l == label {
			return b.bitmap.Test(bit)
		}
		bit++
	}
	return false
}

func (b *LabeledBitmap) SetLabel(label string) {
	bit := uint(0)
	for _, l := range b.labels {
		if l == label {
			b.bitmap.Set(bit)
			return
		}
		bit++
	}
}

func (b *LabeledBitmap) ClearLabel(label string) {
	bit := uint(0)
	for _, l := range b.labels {
		if l == label {
			b.bitmap.Clear(bit)
			return
		}
		bit++
	}
}
