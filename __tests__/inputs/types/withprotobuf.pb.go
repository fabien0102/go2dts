type UserResp struct {
	// user's uuid
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// user's display name
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// user's email address, must be unique in the tenant
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	// will always be empty when returned from the service
	Password string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	// tenant the user belongs to
	TenantId string `protobuf:"bytes,5,opt,name=tenantId,proto3" json:"tenantId,omitempty"`
	// description of the realms the user is a member of
	Realms []*RealmMinimalResp `protobuf:"bytes,6,rep,name=realms,proto3" json:"realms,omitempty"`
	// description of the groups the user is a member of
	Groups []*GroupMinimalResp `protobuf:"bytes,7,rep,name=groups,proto3" json:"groups,omitempty"`
	// list of ssh keys configured for the labs cli
	SshKeys []string `protobuf:"bytes,8,rep,name=sshKeys,proto3" json:"sshKeys,omitempty"`
	// timestamp of the user creation
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,9,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	// timestamp of the last modification to the user instance (name, email, etc)
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,10,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	Image     []byte               `protobuf:"bytes,11,opt,name=image,proto3" json:"image,omitempty"`
	// is the user a tenant admin
	IsAdmin bool `json:"isAdmin,omitempty"`
	// indicates that the user is allowed to login
	IsActive             bool     `protobuf:"varint,13,opt,name=isActive,proto3" json:"isActive,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}