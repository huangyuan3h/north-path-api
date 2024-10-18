package types

type Post struct {
	PostId       string   `json:"postId" dynamodbav:"postId"`
	Email        string   `json:"email" dynamodbav:"email"`
	Subject      string   `json:"subject" dynamodbav:"subject"`
	Content      string   `json:"content" dynamodbav:"content"`
	Category     string   `json:"category" dynamodbav:"category"`
	Location     string   `json:"location" dynamodbav:"location"`
	Bilibili     string   `json:"bilibili" dynamodbav:"bilibili"`
	Youtube      string   `json:"youtube" dynamodbav:"youtube"`
	Topics       []string `json:"topics" dynamodbav:"topics"`
	Images       []string `json:"images" dynamodbav:"images"`
	UpdatedDate  string   `json:"updatedDate" dynamodbav:"updatedDate"`
	Status       string   `json:"status" dynamodbav:"status"`
	Like         int32    `json:"like" dynamodbav:"like"`
	SortingScore int64    `json:"sortingScore" dynamodbav:"sortingScore"`
}

type SearchKeys struct {
	PostId      string `json:"postId" dynamodbav:"postId"`
	UpdatedDate string `json:"updatedDate" dynamodbav:"updatedDate"`
}
