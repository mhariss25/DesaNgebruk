package request

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" `
	NewPassword string `json:"new_password" form:"new_password" `
}
