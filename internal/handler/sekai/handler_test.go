package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"context"
	"log"
	"testing"
)

func TestRegisterCommandHandler(t *testing.T) {

	RegisterSekaiCommandHandler()
	v, e := handler.Dispatch(context.Background(), handler.Event{
		Message: "/cn查谱面 虾",
	})
	log.Println(v, e)
	v, e = handler.Dispatch(context.Background(), handler.Event{
		Message: "/card 1",
	})
	log.Println(v, e)
}
