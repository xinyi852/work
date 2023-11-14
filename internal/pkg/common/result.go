package common

type UserCertQuotaResult struct {
	ProductId   uint64 `json:"product_id"`
	TotalQuota  int64  `json:"total_quota"`
	RemainQuota int64  `json:"remain_quota"`
}

type CertificateStatResult struct {
	TotalNum        int64 `json:"total_num"`
	IssuedNum       int64 `json:"issued_num"`
	ApplyNum        int64 `json:"apply_num"`
	ExpiringSoonNum int64 `json:"expiring_soon_num"`
	ExpiredNum      int64 `json:"expired_num"`
}

type GenerateCSRResult struct {
	CSRPem     string `json:"csr"`
	PrivatePem string `json:"private_key"`
}
