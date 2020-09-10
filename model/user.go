package model

func GetAtuh() {
	Atuh := db.First(&user)
	return Atuh
}
