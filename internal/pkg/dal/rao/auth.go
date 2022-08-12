package rao

type Auth struct {
	Type   string     `json:"type"`
	Kv     *KV        `json:"kv"`
	Bearer *Bearer    `json:"bearer"`
	Basic  *AuthBasic `json:"basic"`
}

type Bearer struct {
	Key string `json:"key"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AuthBasic struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
