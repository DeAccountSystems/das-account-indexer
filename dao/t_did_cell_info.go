package dao

import (
	"das-account-indexer/tables"
	"gorm.io/gorm"
)

func (d *DbDao) CreateDidCellRecordsInfos(outpoint string, didCellInfo tables.TableDidCellInfo, recordsInfos []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id = ?", didCellInfo.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}

		if len(recordsInfos) > 0 {
			if err := tx.Create(&recordsInfos).Error; err != nil {
				return err
			}
		}

		if err := tx.Select("outpoint", "block_number").
			Where("outpoint = ?", outpoint).
			Updates(didCellInfo).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *DbDao) EditDidCellOwner(outpoint string, didCellInfo tables.TableDidCellInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("outpoint", "block_number", "args", "lock_code_hash").
			Where("outpoint = ?", outpoint).
			Updates(didCellInfo).Error; err != nil {
			return err
		}
		return nil

	})
}

func (d *DbDao) DidCellRecycle(outpoint, accountId string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id=?", accountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("outpoint = ? ", outpoint).Delete(&tables.TableDidCellInfo{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *DbDao) QueryDidCell(args string, didType tables.DidCellStatus, limit, offset int) (didList []tables.TableDidCellInfo, err error) {
	sql := d.db.Where(" args= ?", args)
	timestamp := tables.GetDidCellRecycleExpiredAt()
	if didType == tables.DidCellStatusNormal {
		sql.Where("expired_at > ", timestamp)
	} else if didType == tables.DidCellStatusExpired {
		sql.Where("expired_at <= ", timestamp)
	}
	if limit > 0 {
		err = sql.Limit(limit).Offset(offset).Find(&didList).Error
	} else {
		err = sql.Find(&didList).Error
	}
	return
}

func (d *DbDao) QueryDidCellTotal(args string, didType tables.DidCellStatus) (count int64, err error) {
	sql := d.db.Model(tables.TableDidCellInfo{}).Where(" args= ?", args)
	timestamp := tables.GetDidCellRecycleExpiredAt()
	if didType == tables.DidCellStatusNormal {
		sql.Where("expired_at > ", timestamp)
	} else if didType == tables.DidCellStatusExpired {
		sql.Where("expired_at <= ", timestamp)
	}
	err = sql.Count(&count).Error
	return
}

func (d *DbDao) GetDidCellByAccountId(accountId string) (info tables.TableDidCellInfo, err error) {
	err = d.db.Where("account_id=?", accountId).
		Order("expired_at DESC").Limit(1).Find(&info).Error
	return
}

func (d *DbDao) GetAccountInfoByOutpoint(outpoint string) (acc tables.TableDidCellInfo, err error) {
	err = d.db.Where(" outpoint= ? ", outpoint).Find(&acc).Error
	return
}
