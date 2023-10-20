package models

type Transaction struct {
	Number      uint64 `gorm:"column:number;index"`    //当前交易所在区块的高度
	Hash        string `gorm:"column:hash;index"`      //当前交易哈希
	ConfirmTime uint64 `gorm:"column:confirm_time"`    //交易时间
	FileHash    string `gorm:"column:file_hash;index"` //存证的文件哈希，可以为空
	Data        string `gorm:"column:data"`            //交易的数据
}
