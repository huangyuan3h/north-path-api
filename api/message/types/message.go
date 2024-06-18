package types

type Message struct {
	Id          string `json:"id" dynamodbav:"id"`
	FromEmail   string `json:"fromEmail" dynamodbav:"fromEmail"`
	ToEmail     string `json:"toEmail" dynamodbav:"toEmail"`
	Subject     string `json:"subject" dynamodbav:"subject"`
	Content     string `json:"content" dynamodbav:"content"`
	CreatedDate string `json:"createdDate" dynamodbav:"createdDate"`
}
