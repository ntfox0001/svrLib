package payDataStruct

type IapPayDataReq struct {
	Receipt_data             string `json:"receipt-data"`
	Password                 string `json:"password"`
	Exclude_old_transactions string `json:"exclude-old-transactions"`
}

type IapPayDataResp struct {
	Status                      int        `json:"status"`
	Receipt                     IapReceipt `json:"receipt"`
	Latest_receipt              string     `json:"latest_receipt"`
	Latest_receipt_info         string     `json:"latest_receipt_info"`
	Latest_expired_receipt_info string     `json:"latest_expired_receipt_info"`
	Pending_renewal_info        string     `json:"pending_renewal_info"`
	Is_retryable                string     `json:"is-retryable"`
}

type IapReceipt struct {
	Bundle_id                    string            `json:"bundle_id"`
	Application_version          string            `json:"application_version"`
	Original_application_version string            `json:"original_application_version"`
	Receipt_creation_date        string            `json:"receipt_creation_date"`
	Expiration_date              string            `json:"expiration_date"`
	In_app                       []IapInAppReceipt `json:"in_app"`
}

type IapInAppReceipt struct {
	Quantity                    int    `json:"quantity,string"`
	Product_id                  string `json:"product_id"`
	Transaction_id              string `json:"transaction_id"`
	Original_transaction_id     string `json:"original_transaction_id"`
	Purchase_date               string `json:"purchase_date"`
	Original_purchase_date      string `json:"original_purchase_date"`
	Expires_date                string `json:"expires_date"`
	Expiration_intent           string `json:"expiration_intent"`
	Is_in_billing_retry_period  string `json:"is_in_billing_retry_period"`
	Is_trial_period             string `json:"is_trial_period"`
	Is_in_intro_offer_period    string `json:"is_in_intro_offer_period"`
	Cancellation_date           string `json:"cancellation_date"`
	Cancellation_reason         string `json:"cancellation_reason"`
	App_item_id                 string `json:"app_item_id"`
	Version_external_identifier string `json:"version_external_identifier"`
	Web_order_line_item_id      string `json:"web_order_line_item_id"`
	Auto_renew_status           int    `json:"auto_renew_status"`
	Auto_renew_product_id       string `json:"auto_renew_product_id"`
	Price_consent_status        int    `json:"price_consent_status"`
}
