package comment

type CreateCommentRequest struct {
    Content string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
    Content *string `json:"content" binding:"omitempty,min=1"`
}


