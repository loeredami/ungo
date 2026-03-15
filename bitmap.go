package ungo

type Bitmap struct {
	data []byte
	size uint
}

func NewBitmap(size uint) *Bitmap {
	return &Bitmap{
		data: make([]byte, size),
		size: size,
	}
}

func (b *Bitmap) Set(bit uint) {
	b.data[bit/8] |= 1 << (bit % 8)
}

func (b *Bitmap) Clear(bit uint) {
	b.data[bit/8] &= ^(1 << (bit % 8))
}

func (b *Bitmap) Test(bit uint) bool {
	return b.data[bit/8]&(1<<(bit%8)) != 0
}

func (b *Bitmap) Size() uint {
	return b.size
}

func (b *Bitmap) Count() uint {
	count := uint(0)
	for _, v := range b.data {
		count += uint(v)
	}
	return count
}

func (b *Bitmap) Data() []byte {
	return b.data
}

func (b *Bitmap) Dump(arr []byte) {
	copy(arr, b.data)
}

func (b *Bitmap) Load(arr []byte) {
	copy(b.data, arr)
}

func (b *Bitmap) Reset() {
	for i := range b.data {
		b.data[i] = 0
	}
}

func (b *Bitmap) Clone() *Bitmap {
	return &Bitmap{
		data: append([]byte(nil), b.data...),
		size: b.size,
	}
}
