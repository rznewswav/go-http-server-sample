package auth

type Login struct {
	Email    string `form:"email" json:"email" xml:"email"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
