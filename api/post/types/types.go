package types

type Post struct {
	PostId      string   `json:"postId" dynamodbav:"postId"`
	Email       string   `json:"email" dynamodbav:"email"`
	Subject     string   `json:"subject" dynamodbav:"subject"`
	Content     string   `json:"content" dynamodbav:"content"`
	Categories  []string `json:"categories" dynamodbav:"categories"`
	Images      []string `json:"images" dynamodbav:"images"`
	CreatedDate string   `json:"createdDate" dynamodbav:"createdDate"`
	UpdatedDate string   `json:"updatedDate" dynamodbav:"updatedDate"`
}

type SearchKeys struct {
	PostId      string `json:"postId" dynamodbav:"postId"`
	UpdatedDate string `json:"updatedDate" dynamodbav:"updatedDate"`
}
