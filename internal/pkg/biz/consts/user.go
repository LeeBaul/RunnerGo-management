package consts

const (
	RoleTypeOwner  = 1 // 创建者
	RoleTypeMember = 2 // 成员
	RoleTypeAdmin  = 3
)

var (
	DefaultAvatarMemo = map[int]string{
		0: "https://apipost.oss-cn-beijing.aliyuncs.com/kunpeng/avatar/default1.png",
		1: "https://apipost.oss-cn-beijing.aliyuncs.com/kunpeng/avatar/default2.png",
		2: "https://apipost.oss-cn-beijing.aliyuncs.com/kunpeng/avatar/default3.png",
		3: "https://apipost.oss-cn-beijing.aliyuncs.com/kunpeng/avatar/default4.png",
	}
)
