package nops

type Project struct {
	ID            int    `json:"id"`
	Client        int    `json:"client"`
	Arn           string `json:"arn"`
	Bucket        string `json:"bucket"`
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	ExternalID    string `json:"external_id"`
}

type NewProject struct {
	Name                     string `json:"name"`
	AccountNumber            string `json:"account_number"`
	MasterPayerAccountNumber string `json:"master_payer_account_number"`
}

type Notification struct {
	RoleArn            string             `json:"role_arn"`
	BucketName         string             `json:"bucket_name"`
	AccountNumber      string             `json:"account_number"`
	ExternalID         string             `json:"external_id"`
	RequestType        string             `json:"RequestType"`
	ResourceProperties ResourceProperties `json:"ResourceProperties"`
}

type ResourceProperties struct {
	ServiceBucket string `json:"ServiceBucket"`
	AWSAccountID  string `json:"AWSAccountID"`
	RoleArn       string `json:"RoleArn"`
	ExternalID    string `json:"ExternalID"`
}

type NotificationResponse struct {
	Status string `json:"status"`
}
