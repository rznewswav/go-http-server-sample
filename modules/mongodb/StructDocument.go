package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document[T any] struct {
	ID  primitive.ObjectID `bson:"_id"`
	Doc *T
}

func CompileDocument[T any](objectId interface{}, document *T) Document[T] {
	doc := Document[T]{}
	doc.ID = objectId.((primitive.ObjectID))
	doc.Doc = document
	return doc
}
func ExtractDocument[T any](doc Document[T]) (primitive.ObjectID, *T) {
	return doc.ID, doc.Doc
}
