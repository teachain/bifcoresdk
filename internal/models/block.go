package models

import "time"

type Block struct {
	ID           uint      `gorm:"primaryKey;column:id"` //id
	CreatedAt    time.Time `gorm:"column:created_at"`    //数据保存时间
	Number       uint64    `gorm:"column:number;index"`  //当前区块高度
	Hash         string    `gorm:"column:hash;index"`    //当前区块哈希
	TxCount      int       `gorm:"column:tx_count"`      //区块里包含的交易数
	ConfirmTime  uint64    `gorm:"column:confirm_time"`  //出块时间
	PreviousHash string    `gorm:"column:previous_hash"` //上一个区块哈希
}
