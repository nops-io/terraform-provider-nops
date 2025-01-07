package nops

type Project struct {
	ID            int    `json:"id"`
	Client        int    `json:"client"`
	Arn           string `json:"arn"`
	Bucket        string `json:"bucket"`
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	ExternalID    string `json:"external_id"`
	RoleName      string `json:"role_name"`
}

type NewProject struct {
	Name                     string `json:"name"`
	AccountNumber            string `json:"account_number"`
	MasterPayerAccountNumber string `json:"master_payer_account_number"`
}

type UpdateProject struct {
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
}

type Integration struct {
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

type IntegrationResponse struct {
	Status string `json:"status"`
}

type ComputeCopilotOnboarding struct {
	ClusterArns []string `json:"cluster_arns"`
	RegionName  string   `json:"region_name"`
	Version     string   `json:"version"`
	AccountID   string   `json:"account_id"`
}

type ContainerCostBucketSetup struct {
	Project int64 `json:"project"`
}

type ContainerCostBucket struct {
	ID      int64  `json:"id"`
	Project int64  `json:"project"`
	Bucket  string `json:"bucket"`
	Region  string `json:"region"`
	Status  string `json:"status"`
}
