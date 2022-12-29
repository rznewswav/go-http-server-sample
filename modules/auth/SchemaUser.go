package auth

type SchemaUserContactInfo struct {
	Email       string
	PhoneNumber string `bson:"phoneNumber"`
}

type SchemaUser struct {
	Name        string
	HashedP     string                `bson:"hashedP"`
	ContactInfo SchemaUserContactInfo `bson:"contactInfo"`
}

func (schema *SchemaUser) SchemaName() string {
	return "user"
}
