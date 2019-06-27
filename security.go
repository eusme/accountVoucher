package account_voucher

import (
	"github.com/hexya-erp/pool/h"
)

//vars

var ()

//rights
func init() {
	h.AccountVoucher().Methods().Load().AllowGroup(GroupAccountManager)
	h.AccountVoucherLine().Methods().Load().AllowGroup(GroupAccountManager)
	h.AccountVoucher().Methods().AllowAllToGroup(GroupAccountInvoice)
	h.AccountVoucherLine().Methods().AllowAllToGroup(GroupAccountInvoice)
}
