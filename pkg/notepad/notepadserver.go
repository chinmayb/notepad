package notepad

import "context"

type NotePad struct {
}

type NotePadServer interface {
	List(context.Context) ([]NotePad, error)
	Create(context.Context, *NotePad) (NotePad, error)
	Get(context.Context, string) (NotePad, error)
	Update(context.Context, *NotePad) (NotePad, error)
}

type NotePadS struct {
}

func NewNotePadServer() NotePadServer {
	return NotePadS{}
}

func (n NotePadS) List(ctx context.Context) ([]NotePad, error) {
	panic("implement me")
}

func (n NotePadS) Create(ctx context.Context, pad *NotePad) (NotePad, error) {
	panic("implement me")
}

func (n NotePadS) Get(ctx context.Context, s string) (NotePad, error) {
	panic("implement me")
}

func (n NotePadS) Update(ctx context.Context, pad *NotePad) (NotePad, error) {
	panic("implement me")
}
