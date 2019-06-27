package account_voucher

	import (
		"net/http"

		"github.com/hexya-erp/hexya/src/controllers"
		"github.com/hexya-erp/hexya/src/models"
		"github.com/hexya-erp/hexya/src/models/types"
		"github.com/hexya-erp/hexya/src/models/types/dates"
		"github.com/hexya-erp/pool/h"
		"github.com/hexya-erp/pool/q"
	)
	
//import odoo.addons.decimal_precision as dp
func init() {
h.AccountVoucher().DeclareModel()




h.AccountVoucher().Methods().DefaultJournal().DeclareMethod(
`DefaultJournal`,
func(rs m.AccountVoucherSet)  {
//        voucher_type = self._context.get('voucher_type', 'sale')
//        company_id = self._context.get(
//            'company_id', self.env.user.company_id.id)
//        domain = [
//            ('type', '=', voucher_type),
//            ('company_id', '=', company_id),
//        ]
//        return self.env['account.journal'].search(domain, limit=1)
})
h.AccountVoucher().AddFields(map[string]models.FieldDefinition{
"VoucherType": models.SelectionField{
Selection: types.Selection{
"sale": "Sale",
"purchase": "Purchase",
},
String: "Type",
ReadOnly: true,
//states={'draft': [('readonly', False)]}
//oldname="type"
},
"Name": models.CharField{
String: "Payment Reference",
ReadOnly: true,
//states={'draft': [('readonly', False)]}
Default: models.DefaultValue(""),
},
"Date": models.DateField{
String: "Bill Date",
ReadOnly: true,
Index: true,
//states={'draft': [('readonly', False)]}
NoCopy: true,
Default: func (env models.Environment) interface{} { return odoo.fields.Date.context_today },
},
"AccountDate": models.DateField{
String: "Accounting Date",
ReadOnly: true,
Index: true,
//states={'draft': [('readonly', False)]}
Help: "Effective date for accounting entries",
NoCopy: true,
Default: func (env models.Environment) interface{} { return odoo.fields.Date.context_today },
},
"JournalId": models.Many2OneField{
RelationModel: h.AccountJournal(),
String: "Journal",
Required: true,
ReadOnly: true,
//states={'draft': [('readonly', False)]}
Default: models.DefaultValue(_default_journal),
},
"PaymentJournalId": models.Many2OneField{
RelationModel: h.AccountJournal(),
String: "Payment Method",
ReadOnly: true,
Stored: false,
//states={'draft': [('readonly', False)]}
Filter: q.Type().In(%!s(<nil>)),
Compute: h.AccountVoucher().Methods().ComputePaymentJournalId(),
//inverse='_inverse_payment_journal_id'
},
"AccountId": models.Many2OneField{
RelationModel: h.AccountAccount(),
String: "Account",
Required: true,
ReadOnly: true,
//states={'draft': [('readonly', False)]}
Filter: q.Deprecated().Equals(False).And().InternalType().Equals(%!s(<nil>)),
},
"LineIds": models.One2ManyField{
RelationModel: h.AccountVoucherLine(),
ReverseFK: "",
String: "Voucher Lines",
ReadOnly: true,
NoCopy: false,
//states={'draft': [('readonly', False)]}
},
"Narration": models.TextField{
String: "Notes",
ReadOnly: true,
//states={'draft': [('readonly', False)]}
},
"CurrencyId": models.Many2OneField{
RelationModel: h.Currency(),
Compute: h.AccountVoucher().Methods().GetJournalCurrency(),
String: "Currency",
ReadOnly: true,
Required: true,
Default: func (env models.Environment) interface{} { return self._get_currency() },
},
"CompanyId": models.Many2OneField{
RelationModel: h.Company(),
String: "Company",
Required: true,
ReadOnly: true,
//states={'draft': [('readonly', False)]}
Related: `JournalId.CompanyId`,
Default: func (env models.Environment) interface{} { return self._get_company() },
},
"State": models.SelectionField{
Selection: types.Selection{
"draft": "Draft",
"cancel": "Cancelled",
"proforma": "Pro-forma",
"posted": "Posted",
},
String: "Status",
ReadOnly: true,
//track_visibility='onchange'
NoCopy: true,
Default: models.DefaultValue("draft"),
Help: " * The 'Draft' status is used when a user is encoding a" + 
"new and unconfirmed Voucher." + 
" * The 'Pro-forma' status is used when the voucher does" + 
"not have a voucher number." + 
" * The 'Posted' status is used when user create voucher,a" + 
"voucher number is generated and voucher entries are created in account." + 
" * The 'Cancelled' status is used when user cancel voucher.",
},
"Reference": models.CharField{
String: "Bill Reference",
ReadOnly: true,
//states={'draft': [('readonly', False)]}
Help: "The partner reference of this document.",
NoCopy: true,
},
"Amount": models.MonetaryField{
String: "Total",
Stored: true,
ReadOnly: true,
Compute: h.AccountVoucher().Methods().ComputeTotal(),
},
"TaxAmount": models.MonetaryField{
ReadOnly: true,
Stored: true,
Compute: h.AccountVoucher().Methods().ComputeTotal(),
},
"TaxCorrection": models.MonetaryField{
ReadOnly: true,
//states={'draft': [('readonly', False)]}
Help: "In case we have a rounding problem in the tax, use this" + 
"field to correct it",
},
"Number": models.CharField{
ReadOnly: true,
NoCopy: true,
},
"MoveId": models.Many2OneField{
RelationModel: h.AccountMove(),
String: "Journal Entry",
NoCopy: true,
},
"PartnerId": models.Many2OneField{
RelationModel: h.Partner(),
String: "Partner",
//change_default=1
ReadOnly: true,
//states={'draft': [('readonly', False)]}
},
"Paid": models.BooleanField{
Compute: h.AccountVoucher().Methods().CheckPaid(),
Help: "The Voucher has been totally paid.",
},
"PayNow": models.SelectionField{
Selection: types.Selection{
"pay_now": "Pay Directly",
"pay_later": "Pay Later",
},
String: "Payment",
Index: true,
ReadOnly: true,
//states={'draft': [('readonly', False)]}
Default: models.DefaultValue("pay_later"),
},
"DateDue": models.DateField{
String: "Due Date",
ReadOnly: true,
Index: true,
//states={'draft': [('readonly', False)]}
},
})
h.AccountVoucher().Methods().CheckPaid().DeclareMethod(
`CheckPaid`,
func(rs h.AccountVoucherSet) h.AccountVoucherData {
//        self.paid = any([((line.account_id.internal_type, 'in', ('receivable',
//                                                                 'payable')) and line.reconciled) for line in self.move_id.line_ids])
})
h.AccountVoucher().Methods().GetCurrency().DeclareMethod(
`GetCurrency`,
func(rs m.AccountVoucherSet)  {
//        journal = self.env['account.journal'].browse(
//            self.env.context.get('default_journal_id', False))
//        if journal.currency_id:
//            return journal.currency_id.id
//        return self.env.user.company_id.currency_id.id
})
h.AccountVoucher().Methods().GetCompany().DeclareMethod(
`GetCompany`,
func(rs m.AccountVoucherSet)  {
//        return self._context.get('company_id', self.env.user.company_id.id)
})
h.AccountVoucher().Methods().NameGet().Extend(
`NameGet`,
func(rs m.AccountVoucherSet)  {
//        return [(r.id, (r.number or _('Voucher'))) for r in self]
})
h.AccountVoucher().Methods().GetJournalCurrency().DeclareMethod(
`GetJournalCurrency`,
func(rs h.AccountVoucherSet) h.AccountVoucherData {
//        self.currency_id = self.journal_id.currency_id.id or self.company_id.currency_id.id
})
h.AccountVoucher().Methods().ComputePaymentJournalId().DeclareMethod(
`ComputePaymentJournalId`,
func(rs h.AccountVoucherSet) h.AccountVoucherData {
//        for voucher in self:
//            if voucher.pay_now != 'pay_now':
//                continue
//            domain = [
//                ('type', 'in', ('bank', 'cash')),
//                ('company_id', '=', voucher.company_id.id),
//            ]
//            if voucher.account_id and voucher.account_id.internal_type == 'liquidity':
//                field = 'default_debit_account_id' if voucher.voucher_type == 'sale' else 'default_credit_account_id'
//                domain.append((field, '=', voucher.account_id.id))
//            voucher.payment_journal_id = self.env['account.journal'].search(
//                domain, limit=1)
})
h.AccountVoucher().Methods().InversePaymentJournalId().DeclareMethod(
`InversePaymentJournalId`,
func(rs m.AccountVoucherSet)  {
//        for voucher in self:
//            if voucher.pay_now != 'pay_now':
//                continue
//            if voucher.voucher_type == 'sale':
//                voucher.account_id = voucher.payment_journal_id.default_debit_account_id
//            else:
//                voucher.account_id = voucher.payment_journal_id.default_credit_account_id
})
h.AccountVoucher().Methods().ComputeTotal().DeclareMethod(
`ComputeTotal`,
func(rs h.AccountVoucherSet) h.AccountVoucherData {
//        for voucher in self:
//            total = 0
//            tax_amount = 0
//            for line in voucher.line_ids:
//                tax_info = line.tax_ids.compute_all(
//                    line.price_unit, voucher.currency_id, line.quantity, line.product_id, voucher.partner_id)
//                total += tax_info.get('total_included', 0.0)
//                tax_amount += sum([t.get('amount', 0.0)
//                                   for t in tax_info.get('taxes', False)])
//            voucher.amount = total + voucher.tax_correction
//            voucher.tax_amount = tax_amount
})
h.AccountVoucher().Methods().OnchangeDate().DeclareMethod(
`OnchangeDate`,
func(rs m.AccountVoucherSet)  {
//        self.account_date = self.date
})
h.AccountVoucher().Methods().OnchangePartnerId().DeclareMethod(
`OnchangePartnerId`,
func(rs m.AccountVoucherSet)  {
//        if self.pay_now != 'pay_now':
//            if self.partner_id:
//                self.account_id = self.partner_id.property_account_receivable_id \
//                    if self.voucher_type == 'sale' else self.partner_id.property_account_payable_id
//            else:
//                self.account_id = self.journal_id.default_debit_account_id \
//                    if self.voucher_type == 'sale' else self.journal_id.default_credit_account_id
})
h.AccountVoucher().Methods().ProformaVoucher().DeclareMethod(
`ProformaVoucher`,
func(rs m.AccountVoucherSet)  {
//        self.action_move_line_create()
})
h.AccountVoucher().Methods().ActionCancelDraft().DeclareMethod(
`ActionCancelDraft`,
func(rs m.AccountVoucherSet)  {
//        self.write({'state': 'draft'})
})
h.AccountVoucher().Methods().CancelVoucher().DeclareMethod(
`CancelVoucher`,
func(rs m.AccountVoucherSet)  {
//        for voucher in self:
//            voucher.move_id.button_cancel()
//            voucher.move_id.unlink()
//        self.write({'state': 'cancel', 'move_id': False})
})
h.AccountVoucher().Methods().Unlink().Extend(
`Unlink`,
func(rs m.AccountVoucherSet)  {
//        for voucher in self:
//            if voucher.state not in ('draft', 'cancel'):
//                raise UserError(
//                    _('Cannot delete voucher(s) which are already opened or paid.'))
//        return super(AccountVoucher, self).unlink()
})
h.AccountVoucher().Methods().FirstMoveLineGet().DeclareMethod(
`FirstMoveLineGet`,
func(rs m.AccountVoucherSet, move_id interface{}, company_currency interface{}, current_currency interface{})  {
//        debit = credit = 0.0
//        if self.voucher_type == 'purchase':
//            credit = self._convert_amount(self.amount)
//        elif self.voucher_type == 'sale':
//            debit = self._convert_amount(self.amount)
//        if debit < 0.0:
//            debit = 0.0
//        if credit < 0.0:
//            credit = 0.0
//        sign = debit - credit < 0 and -1 or 1
//        move_line = {
//            'name': self.name or '/',
//            'debit': debit,
//            'credit': credit,
//            'account_id': self.account_id.id,
//            'move_id': move_id,
//            'journal_id': self.journal_id.id,
//            'partner_id': self.partner_id.commercial_partner_id.id,
//            'currency_id': company_currency != current_currency and current_currency or False,
//            'amount_currency': (sign * abs(self.amount)  # amount < 0 for refunds
//                                if company_currency != current_currency else 0.0),
//            'date': self.account_date,
//            'date_maturity': self.date_due,
//            'payment_id': self._context.get('payment_id'),
//        }
//        return move_line
})
h.AccountVoucher().Methods().AccountMoveGet().DeclareMethod(
`AccountMoveGet`,
func(rs m.AccountVoucherSet)  {
//        if self.number:
//            name = self.number
//        elif self.journal_id.sequence_id:
//            if not self.journal_id.sequence_id.active:
//                raise UserError(
//                    _('Please activate the sequence of selected journal !'))
//            name = self.journal_id.sequence_id.with_context(
//                ir_sequence_date=self.date).next_by_id()
//        else:
//            raise UserError(_('Please define a sequence on the journal.'))
//        move = {
//            'name': name,
//            'journal_id': self.journal_id.id,
//            'narration': self.narration,
//            'date': self.account_date,
//            'ref': self.reference,
//        }
//        return move
})
h.AccountVoucher().Methods().ConvertAmount().DeclareMethod(
`
        This function convert the amount given in company
currency. It takes either the rate in the voucher (if the
        payment_rate_currency_id is relevant) either the
rate encoded in the system.
        :param amount: float. The amount to convert
        :param voucher: id of the voucher on which we want
the conversion
        :param context: to context to use for the conversion.
It may contain the key 'date' set to the voucher date
            field in order to select the good rate to use.
        :return: the amount in the currency of the voucher's company
        :rtype: float
        `,
func(rs m.AccountVoucherSet, amount interface{})  {
//        for voucher in self:
//            return voucher.currency_id.compute(amount, voucher.company_id.currency_id)
})
h.AccountVoucher().Methods().VoucherPayNowPaymentCreate().DeclareMethod(
`VoucherPayNowPaymentCreate`,
func(rs m.AccountVoucherSet)  {
//        if self.voucher_type == 'sale':
//            payment_methods = self.journal_id.inbound_payment_method_ids
//            payment_type = 'inbound'
//            partner_type = 'customer'
//            sequence_code = 'account.payment.customer.invoice'
//        else:
//            payment_methods = self.journal_id.outbound_payment_method_ids
//            payment_type = 'outbound'
//            partner_type = 'supplier'
//            sequence_code = 'account.payment.supplier.invoice'
//        name = self.env['ir.sequence'].with_context(
//            ir_sequence_date=self.date).next_by_code(sequence_code)
//        return {
//            'name': name,
//            'payment_type': payment_type,
//            'payment_method_id': payment_methods and payment_methods[0].id or False,
//            'partner_type': partner_type,
//            'partner_id': self.partner_id.commercial_partner_id.id,
//            'amount': self.amount,
//            'currency_id': self.currency_id.id,
//            'payment_date': self.date,
//            'journal_id': self.payment_journal_id.id,
//            'company_id': self.company_id.id,
//            'communication': self.name,
//            'state': 'reconciled',
//        }
})
h.AccountVoucher().Methods().VoucherMoveLineCreate().DeclareMethod(
`
        Create one account move line, on the given account
move, per voucher line where amount is not 0.0.
        It returns Tuple with tot_line what is total of
difference between debit and credit and
        a list of lists with ids to be reconciled with
this format (total_deb_cred,list_of_lists).

        :param voucher_id: Voucher id what we are working with
        :param line_total: Amount of the first line, which
correspond to the amount we should totally split among
all voucher lines.
        :param move_id: Account move wher those lines will be joined.
        :param company_currency: id of currency of the
company to which the voucher belong
        :param current_currency: id of currency of the voucher
        :return: Tuple build as (remaining amount not allocated
on voucher lines, list of account_move_line created in this method)
        :rtype: tuple(float, list of int)
        `,
func(rs m.AccountVoucherSet, line_total interface{}, move_id interface{}, company_currency interface{}, current_currency interface{})  {
//        for line in self.line_ids:
//            # create one move line per voucher line where amount is not 0.0
//            if not line.price_subtotal:
//                continue
//            line_subtotal = line.price_subtotal
//            if self.voucher_type == 'sale':
//                line_subtotal = -1 * line.price_subtotal
//            # convert the amount set on the voucher line into the currency of the voucher's company
//            # this calls res_curreny.compute() with the right context,
//            # so that it will take either the rate on the voucher if it is relevant or will use the default behaviour
//            amount = self._convert_amount(line.price_unit*line.quantity)
//            move_line = {
//                'journal_id': self.journal_id.id,
//                'name': line.name or '/',
//                'account_id': line.account_id.id,
//                'move_id': move_id,
//                'partner_id': self.partner_id.commercial_partner_id.id,
//                'analytic_account_id': line.account_analytic_id and line.account_analytic_id.id or False,
//                'quantity': 1,
//                'credit': abs(amount) if self.voucher_type == 'sale' else 0.0,
//                'debit': abs(amount) if self.voucher_type == 'purchase' else 0.0,
//                'date': self.account_date,
//                'tax_ids': [(4, t.id) for t in line.tax_ids],
//                'amount_currency': line_subtotal if current_currency != company_currency else 0.0,
//                'currency_id': company_currency != current_currency and current_currency or False,
//                'payment_id': self._context.get('payment_id'),
//            }
//            self.env['account.move.line'].with_context(
//                apply_taxes=True).create(move_line)
//        return line_total
})
h.AccountVoucher().Methods().ActionMoveLineCreate().DeclareMethod(
`
        Confirm the vouchers given in ids and create the
journal entries for each of them
        `,
func(rs m.AccountVoucherSet)  {
//        for voucher in self:
//            local_context = dict(
//                self._context, force_company=voucher.journal_id.company_id.id)
//            if voucher.move_id:
//                continue
//            company_currency = voucher.journal_id.company_id.currency_id.id
//            current_currency = voucher.currency_id.id or company_currency
//            # we select the context to use accordingly if it's a multicurrency case or not
//            # But for the operations made by _convert_amount, we always need to give the date in the context
//            ctx = local_context.copy()
//            ctx['date'] = voucher.account_date
//            ctx['check_move_validity'] = False
//            # Create a payment to allow the reconciliation when pay_now = 'pay_now'.
//            if self.pay_now == 'pay_now' and self.amount > 0:
//                ctx['payment_id'] = self.env['account.payment'].create(
//                    self.voucher_pay_now_payment_create()).id
//            # Create the account move record.
//            move = self.env['account.move'].create(voucher.account_move_get())
//            # Get the name of the account_move just created
//            # Create the first line of the voucher
//            move_line = self.env['account.move.line'].with_context(ctx).create(voucher.with_context(
//                ctx).first_move_line_get(move.id, company_currency, current_currency))
//            line_total = move_line.debit - move_line.credit
//            if voucher.voucher_type == 'sale':
//                line_total = line_total - \
//                    voucher._convert_amount(voucher.tax_amount)
//            elif voucher.voucher_type == 'purchase':
//                line_total = line_total + \
//                    voucher._convert_amount(voucher.tax_amount)
//            # Create one move line per voucher line where amount is not 0.0
//            line_total = voucher.with_context(ctx).voucher_move_line_create(
//                line_total, move.id, company_currency, current_currency)
//
//            # Add tax correction to move line if any tax correction specified
//            if voucher.tax_correction != 0.0:
//                tax_move_line = self.env['account.move.line'].search(
//                    [('move_id', '=', move.id), ('tax_line_id', '!=', False)], limit=1)
//                if len(tax_move_line):
//                    tax_move_line.write({'debit': tax_move_line.debit + voucher.tax_correction if tax_move_line.debit > 0 else 0,
//                                         'credit': tax_move_line.credit + voucher.tax_correction if tax_move_line.credit > 0 else 0})
//
//            # We post the voucher.
//            voucher.write({
//                'move_id': move.id,
//                'state': 'posted',
//                'number': move.name
//            })
//            move.post()
//        return True
})
h.AccountVoucher().Methods().TrackSubtype().DeclareMethod(
`TrackSubtype`,
func(rs m.AccountVoucherSet, init_values interface{})  {
//        if 'state' in init_values:
//            return 'account_voucher.mt_voucher_state_change'
//        return super(AccountVoucher, self)._track_subtype(init_values)
})
h.AccountVoucherLine().DeclareModel()


h.AccountVoucherLine().AddFields(map[string]models.FieldDefinition{
"Name": models.TextField{
String: "Description",
Required: true,
},
"Sequence": models.IntegerField{
Default: models.DefaultValue(10),
Help: "Gives the sequence of this line when displaying the voucher.",
},
"VoucherId": models.Many2OneField{
RelationModel: h.AccountVoucher(),
String: "Voucher",
Required: true,
OnDelete: `cascade`,
},
"ProductId": models.Many2OneField{
RelationModel: h.ProductProduct(),
String: "Product",
OnDelete: `set null`,
Index: true,
},
"AccountId": models.Many2OneField{
RelationModel: h.AccountAccount(),
String: "Account",
Required: true,
Filter: q.Deprecated().Equals(False),
Help: "The income or expense account related to the selected product.",
},
"PriceUnit": models.FloatField{
String: "Unit Price",
Required: true,
//digits=dp.get_precision('Product Price')
//oldname='amount'
},
"PriceSubtotal": models.MonetaryField{
String: "Amount",
Stored: true,
ReadOnly: true,
Compute: h.AccountVoucherLine().Methods().ComputeSubtotal(),
},
"Quantity": models.FloatField{
//digits=dp.get_precision('Product Unit of Measure')
Required: true,
Default: models.DefaultValue(1),
},
"AccountAnalyticId": models.Many2OneField{
RelationModel: h.AccountAnalyticAccount(),
String: "Analytic Account",
},
"CompanyId": models.Many2OneField{
RelationModel: h.Company(),
Related: `VoucherId.CompanyId`,
String: "Company",
Stored: true,
ReadOnly: true,
},
"TaxIds": models.Many2ManyField{
RelationModel: h.AccountTax(),
String: "Tax",
Help: "Only for tax excluded from price",
},
"CurrencyId": models.Many2OneField{
RelationModel: h.Currency(),
Related: `VoucherId.CurrencyId`,
},
})
h.AccountVoucherLine().Methods().ComputeSubtotal().DeclareMethod(
`ComputeSubtotal`,
func(rs h.AccountVoucherLineSet) h.AccountVoucherLineData {
//        self.price_subtotal = self.quantity * self.price_unit
//        if self.tax_ids:
//            taxes = self.tax_ids.compute_all(self.price_unit, self.voucher_id.currency_id,
//                                             self.quantity, product=self.product_id, partner=self.voucher_id.partner_id)
//            self.price_subtotal = taxes['total_excluded']
})
h.AccountVoucherLine().Methods().OnchangeLineDetails().DeclareMethod(
`OnchangeLineDetails`,
func(rs m.AccountVoucherLineSet)  {
//        if not self.voucher_id or not self.product_id or not self.voucher_id.partner_id:
//            return
//        onchange_res = self.product_id_change(
//            self.product_id.id,
//            self.voucher_id.partner_id.id,
//            self.price_unit,
//            self.company_id.id,
//            self.voucher_id.currency_id.id,
//            self.voucher_id.voucher_type)
//        for fname, fvalue in onchange_res['value'].iteritems():
//            setattr(self, fname, fvalue)
})
h.AccountVoucherLine().Methods().GetAccount().DeclareMethod(
`GetAccount`,
func(rs m.AccountVoucherLineSet, product interface{}, fpos interface{}, typeName interface{})  {
//        accounts = product.product_tmpl_id.get_product_accounts(fpos)
//        if type == 'sale':
//            return accounts['income']
//        return accounts['expense']
})
h.AccountVoucherLine().Methods().ProductIdChange().DeclareMethod(
`ProductIdChange`,
func(rs m.AccountVoucherLineSet, product_id interface{}, partner_id interface{}, price_unit interface{}, company_id interface{}, currency_id interface{}, typeName interface{})  {
//        context = self._context
//        company_id = company_id if company_id is not None else context.get(
//            'company_id', False)
//        company = self.env['res.company'].browse(company_id)
//        currency = self.env['res.currency'].browse(currency_id)
//        if not partner_id:
//            raise UserError(_("You must first select a partner!"))
//        part = self.env['res.partner'].browse(partner_id)
//        if part.lang:
//            self = self.with_context(lang=part.lang)
//        product = self.env['product.product'].browse(product_id)
//        fpos = part.property_account_position_id
//        account = self._get_account(product, fpos, type)
//        values = {
//            'name': product.partner_ref,
//            'account_id': account.id,
//        }
//        if type == 'purchase':
//            values['price_unit'] = price_unit or product.standard_price
//            taxes = product.supplier_taxes_id or account.tax_ids
//            if product.description_purchase:
//                values['name'] += '\n' + product.description_purchase
//        else:
//            values['price_unit'] = price_unit or product.lst_price
//            taxes = product.taxes_id or account.tax_ids
//            if product.description_sale:
//                values['name'] += '\n' + product.description_sale
//        values['tax_ids'] = taxes.ids
//        if company and currency:
//            if company.currency_id != currency:
//                if type == 'purchase':
//                    values['price_unit'] = price_unit or product.standard_price
//                values['price_unit'] = values['price_unit'] * currency.rate
//        return {'value': values, 'domain': {}}
})
}