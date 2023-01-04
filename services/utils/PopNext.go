package utils

import "container/list"

type PopNextHandler struct {
	items *list.List
}

func (handler *PopNextHandler) Next() interface{} {
	toRemove := handler.items.Back()
	if toRemove == nil {
		return nil
	}
	handler.items.Remove(toRemove)
	return toRemove.Value
}

func PrepareMockReturnValue(items ...interface{}) *PopNextHandler {
	handler := PopNextHandler{
		items: list.New(),
	}

	for _, v := range items {
		handler.items.PushBack(v)
	}

	return &handler
}

func (handler *PopNextHandler) PrepareMockReturnValue(items ...interface{}) {
	if handler.items == nil {
		handler.items = list.New()
	}

	for _, v := range items {
		handler.items.PushBack(v)
	}
}
