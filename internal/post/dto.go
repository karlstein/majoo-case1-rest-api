package post

type CreatePostRequest struct {
    Title   string `json:"title" binding:"required,min=1,max=255"`
    Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
    Title   *string `json:"title" binding:"omitempty,min=1,max=255"`
    Content *string `json:"content" binding:"omitempty,min=1"`
}


