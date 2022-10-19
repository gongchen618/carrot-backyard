package model

func init() {
	m := GetModel()
	defer m.Close()
	Token, _ = m.GetToken()
}
