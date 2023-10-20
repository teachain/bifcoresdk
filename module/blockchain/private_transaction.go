package blockchain

import (
	"encoding/json"
	"github.com/teachain/bifcoresdk/common"
	"github.com/teachain/bifcoresdk/exception"
	"github.com/teachain/bifcoresdk/types/request"
	"github.com/teachain/bifcoresdk/types/response"
	"github.com/teachain/bifcoresdk/utils/http"
)

type BIFPrivateTransactionService interface {
	// Send 私有化交易合约内容
	Send(request.BIFPrivateTransactionSendRequest) response.BIFPrivateTransactionSendResponse
}

// PrivateTransactionService ...
type PrivateTransactionService struct {
	url string
}

func GetPrivateTransactionInstance(url string) *PrivateTransactionService {
	return &PrivateTransactionService{
		url,
	}
}

func (ps *PrivateTransactionService) Send(r request.BIFPrivateTransactionSendRequest) response.BIFPrivateTransactionSendResponse {

	if r.Payload == "" {
		return response.BIFPrivateTransactionSendResponse{
			BIFBaseResponse: exception.INVALID_PRITX_PAYLAOD_ERROR,
		}
	}
	if r.From == "" {
		return response.BIFPrivateTransactionSendResponse{
			BIFBaseResponse: exception.INVALID_PRITX_FROM_ERROR,
		}
	}
	if len(r.To) == 0 {
		return response.BIFPrivateTransactionSendResponse{
			BIFBaseResponse: exception.INVALID_PRITX_TO_ERROR,
		}
	}

	params := make(map[string]interface{})
	params["payload"] = r.Payload
	params["from"] = r.From
	params["to"] = r.To
	transactionSendRequest, err := json.Marshal(params)
	if err != nil {
		return response.BIFPrivateTransactionSendResponse{
			BIFBaseResponse: exception.SYSTEM_ERROR,
		}
	}
	priTxSendUrl := common.PriTxSend(ps.url)
	dataByte, err := http.HttpPost(priTxSendUrl, transactionSendRequest)
	if err != nil {
		return response.BIFPrivateTransactionSendResponse{
			BIFBaseResponse: exception.CONNECTNETWORK_ERROR,
		}
	}

	var res response.BIFPrivateTransactionSendResponse
	err = json.Unmarshal(dataByte, &res)
	if err != nil {
		return response.BIFPrivateTransactionSendResponse{
			BIFBaseResponse: exception.SYSTEM_ERROR,
		}
	}

	return res
}
