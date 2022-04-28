package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
)

func (b *BlockParser) ActionRedeclareReverseRecord(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameReverseRecordCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersionTx err: %s", err.Error())
		return
	} else if !isCV {
		return
	}

	log.Info("das tx:", req.Action, req.TxHash)

	account := string(req.Tx.OutputsData[0])
	ownerHex, _, err := b.DasCore.Daf().ArgsToHex(req.Tx.Outputs[0].Lock.Args)
	if err != nil {
		resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
		return
	}

	reverseInfo := tables.TableReverseInfo{
		BlockNumber:    req.BlockNumber,
		BlockTimestamp: req.BlockTimestamp,
		Outpoint:       common.OutPoint2String(req.TxHash, 0),
		AlgorithmId:    ownerHex.DasAlgorithmId,
		ChainType:      ownerHex.ChainType,
		Address:        ownerHex.AddressHex,
		Account:        account,
		Capacity:       req.Tx.Outputs[0].Capacity,
	}
	lastOutpoint := common.OutPointStruct2String(req.Tx.Inputs[0].PreviousOutput)

	if err := b.DbDao.UpdateReverseInfo(&reverseInfo, lastOutpoint); err != nil {
		resp.Err = fmt.Errorf("UpdateReverseInfo err: %s", err.Error())
		return
	}
	return
}
