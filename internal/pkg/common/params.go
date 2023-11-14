package common

type QuotaRequest struct {
	ProductId uint64 `json:"product_id" valid:"product_id"`
	Total     int64  `json:"total" valid:"total"`
}

type UserCertQuotaRequest struct {
	UserId uint64          `json:"user_id" valid:"user_id"`
	Quota  []*QuotaRequest `json:"quota" valid:"quota"`
}

type Vendor struct {
	VendorId string `json:"vendor_id" form:"vendor_id" valid:"vendor_id"`
}

type ApplyTemplateRequest struct {
	Vendor
	CertType           string `json:"cert_type" form:"cert_type" valid:"cert_type"`
	CertValidationType string `json:"cert_validation_type" form:"cert_validation_type" valid:"cert_validation_type"`
	ProductCode        string `json:"product_code" form:"product_code" valid:"product_code"`
}

type PreDCVInfoRequest struct {
	Vendor
	CsrStr      string `json:"csr_str" valid:"csr_str"`
	UniqueValue string `json:"unique_value" valid:"unique_value"`
}

type GenerateCSRRequest struct {
	CommonName             string `json:"common_name" valid:"common_name"`             //域名
	Organization           string `json:"organization" valid:"organization"`           //企业/单位名称
	OrganizationalUnit     string `json:"organization_unit" valid:"organization_unit"` //部门
	Country                string `json:"country" valid:"country"`                     //国家
	Province               string `json:"province" valid:"province"`                   //省份
	Locality               string `json:"locality" valid:"locality"`                   //城市
	Postcode               string `json:"postcode" valid:"postcode"`
	Address                string `json:"address" valid:"address"`
	EmailAddress           string `json:"email_address"`               //subject 的邮箱地址
	Email                  string `json:"email" valid:"email"`         //邮箱
	DNSNames               string `json:"dns_names" valid:"dns_names"` //备用域名
	IpAddress              string `json:"ip_address"`
	KeyAlgorithm           string `json:"key_algorithm" valid:"key_algorithm"`
	KeySize                string `json:"key_size" valid:"key_size"`
	SignatureHashAlgorithm string `json:"signature_hash_algorithm" valid:"signature_hash_algorithm"`
}

type ApplyRequest struct {
	Vendor
	Product         *ApplyProduct      `json:"product"`
	Csr             string             `json:"csr" valid:"csr"`
	AdminContact    *ApplyContact      `json:"admin_contact"`
	TechContact     *ApplyContact      `json:"tech_contact"`
	FinanceContact  *ApplyContact      `json:"finance_contact"`
	Organization    *ApplyOrganization `json:"organization"`
	Domains         []ApplyDomain      `json:"domains" valid:"domains"`
	GiftDays        uint8              `json:"gift_days"`
	SupplierId      string             `json:"supplier_id"`
	PrimaryDomain   string             `json:"primary_domain" valid:"primary_domain"`
	CustomerOrderId string             `json:"customer_order_id"`
	UniqueValue     string             `json:"unique_value"`
	PrivateKey      string             `json:"private_key"`
	UserId          uint64             `json:"user_id" valid:"user_id"`
}

type ApplyProduct struct {
	ProductId      uint64 `json:"product_id" valid:"product_id"`
	ProductCode    string `json:"product_code" valid:"product_code"`
	SslType        string `json:"ssl_type"`
	IsMulti        bool   `json:"is_multi"`
	ValidationType string `json:"validation_type"`
	IsWild         bool   `json:"is_wild"`
}

type ApplyContact struct {
	Organization string `json:"organization"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Job          string `json:"job"`
	Country      string `json:"country"`
	State        string `json:"state"`
	City         string `json:"city"`
	Address      string `json:"address"`
	PostCode     string `json:"post_code"`
	Mobile       string `json:"mobile"`
	Email        string `json:"email"`
	IdType       string `json:"id_type"`
	IdNumber     string `json:"id_number"`
}

type ApplyOrganization struct {
	Name       string `json:"name"`
	Country    string `json:"country"`
	State      string `json:"state"`
	City       string `json:"city"`
	PostCode   string `json:"post_code"`
	Address    string `json:"address"`
	Mobile     string `json:"mobile"`
	RegNumber  string `json:"reg_number"`
	IdType     string `json:"id_type"`
	SupplierId string `json:"supplier_id"`
}

type ApplyDomain struct {
	Name      string `json:"name"`
	DcvMethod string `json:"dcv_method"`
	DcvEmail  string `json:"dcv_email"`
	IsVerify  string `json:"is_verify"`
}

type PaginationRequest struct {
	Sort    string `valid:"sort" json:"sort" form:"sort"`
	Order   string `valid:"order" json:"order" form:"order"`
	PerPage int    `valid:"per_page" json:"per_page" form:"per_page"`
	Page    int    `valid:"page" json:"page" form:"page"`
}

type FreeCertListRequest struct {
	PaginationRequest
	UserId            uint64 `form:"user_id" valid:"user_id"`
	SubjectCommonName string `form:"subject_common_name"`
	Status            int8   `form:"status"`
	SerialNum         string `form:"serial_num"`
	BeginTime         string `form:"begin_time"`
	EndTime           string `form:"end_time"`
}

type ServerGropuListRequest struct {
	PaginationRequest
	Name string `json:"name"`
}

type CreateHostingRequest struct {
	Domain      string `json:"domain" valid:"domain"`
	Plan        string `json:"plan" valid:"plan"`
	Area        string `json:"area" valid:"area"`
	Module      string `json:"module" valid:"module"`
	Login       string ` json:"login"`
	Email       string ` json:"email"`
	Company     string ` json:"company"`
	Description string ` json:"description"`
	ExternalId  string ` json:"external_id"` //实例号
}

type UpgradePleskRequest struct {
	ID   string `json:"id" valid:"id"`
	Plan string `json:"plan" valid:"plan"`
}

type SetPleskRequest struct {
	ID string `json:"id" valid:"id"`
}

type ChangePasswordRequest struct {
	ID string `json:"id" valid:"id"`
	// Account  string `json:"account" valid:"account"`
	Password string `json:"password" valid:"password"`
	// NewPassword string `json:"new_password" valid:"new_password"`
}

type UpgradePlanRequest struct {
	Name string `json:"name"`
	Plan Plan   `json:"plan"`
}

type Result struct {
	Status string `xml:"status"`
	ID     string `xml:"id"`
}

type CreateSession struct {
	Result Result `xml:"result"`
}

type Server struct {
	CreateSession CreateSession `xml:"create_session"`
}

type Packet struct {
	Version string `xml:"version,attr"`
	Server  Server `xml:"server"`
}

type Plan struct {
	Name string `json:"name"`
}

type ChangePassword struct {
	// Login    string `json:"login"`
	Password string `json:"password"`
	// Plan     Plan   `json:"plan"`
}

type CreatePleskServerRequest struct {
	Name      string `valid:"name" json:"name"`
	Account   string `valid:"account" json:"account"`
	Password  string `valid:"password" json:"password"`
	Max       string `valid:"max" json:"max"`
	Service   string `valid:"service" json:"service"`
	Ip        string `valid:"ip" json:"ip"`
	Url       string `valid:"url" json:"url"`
	AreaGroup string `valid:"server_group" json:"server_group"`
	Type      string `valid:"type" json:"type"`
}

type UpdatePleskServerRequest struct {
	Id        string `valid:"id" json:"id"`
	Name      string `valid:"name" json:"name"`
	Account   string `valid:"account" json:"account"`
	Password  string `valid:"password" json:"password"`
	Max       string `valid:"max" json:"max"`
	Service   string `valid:"service" json:"service"`
	Ip        string `valid:"ip" json:"ip"`
	Url       string `valid:"url" json:"url"`
	AreaGroup string `valid:"server_group" json:"server_group"`
	Type      string `valid:"type" json:"type"`
}

type CreatePlanRequest struct {
	Name       string `valid:"name" json:"name"`
	Plan       string `valid:"plan" json:"plan"`
	MaxData    string `valid:"max_data" json:"max_data"`
	MaxTraffic string `valid:"max_traffic" json:"max_traffic"`
}

type UpdatePlanRequest struct {
	Id         string `valid:"id" json:"id"`
	Name       string `valid:"name" json:"name"`
	Plan       string `valid:"plan" json:"plan"`
	MaxData    string `valid:"max_data" json:"max_data"`
	MaxTraffic string `valid:"max_traffic" json:"max_traffic"`
}

type DeletePlanRequest struct {
	Id string `valid:"id" json:"id"`
}

type DeleteServerRequest struct {
	Id string `valid:"id" json:"id"`
}

type CreatePleskDomainRequest struct {
	Domain   string `json:"domain"`
	Password string `valid:"password" json:"password"`
	Max      string `valid:"max" json:"max"`
	Service  string `valid:"service" json:"service"`
	Ip       string `valid:"ip" json:"ip"`
	Url      string `valid:"url" json:"url"`
}
