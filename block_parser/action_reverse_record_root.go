package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"gorm.io/gorm"
)

func (b *BlockParser) ActionReverseRecordRoot(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameReverseRecordRootCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err)
		return
	} else if !isCV {
		return
	}
	log.Info("ActionReverseRecordRoot:", req.BlockNumber, req.TxHash)

	smtBuilder := witness.NewReverseSmtBuilder()
	txReverseSmtRecord, err := smtBuilder.FromTx(req.Tx)
	if err != nil {
		resp.Err = err
		return
	}

	if err := b.DbDao.Transaction(func(tx *gorm.DB) error {
		for idx, v := range txReverseSmtRecord {
			outpoint := common.OutPoint2String(req.TxHash, uint(idx))
			algorithmId := common.DasAlgorithmId(v.SignType)
			reverseInfo := &tables.TableReverseInfo{
				BlockNumber:    req.BlockNumber,
				BlockTimestamp: req.BlockTimestamp,
				Outpoint:       outpoint,
				AlgorithmId:    algorithmId,
				ChainType:      algorithmId.ToChainType(),
				Address:        common.Bytes2Hex(v.Address),
				Account:        v.NextAccount,
				ReverseType:    tables.ReverseTypeSmt,
			}

			if v.PrevAccount != "" {
				if err := tx.Where("address=? and reverse_type=?", v.Address, tables.ReverseTypeSmt).Delete(&tables.TableReverseInfo{}).Error; err != nil {
					return err
				}
			}
			if v.Action == witness.ReverseSmtRecordActionUpdate {
				if err := tx.Create(reverseInfo).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		resp.Err = err
		return
	}
	return
}
