package domain

type Admin struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password_hash string `json:"password_hash"`
	Role          string `json:"role"`
	Created_at    string `json:"created_at"`
	Updated_at    string `json:"updated_at"`
	Created_by    string `json:"created_by"`
	Updated_by    string `json:"updated_by"`
}
