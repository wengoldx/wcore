// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2019/06/30   zhaixing       Add function from godfs
// -------------------------------------------------------------------

package invar

// Trade status machine:
// ===================================================================
//
// TRADE :  TSUnpaid --- + -> TSRevoked --------------------> TSClosed
//                       |      ^                                 ^
//                       |      |                                 |
//                       + -> TSPayError <- + (max counts 5)      |
//                       |      |        |                        |
//                       |      + ------ +                        |
//                       |      v                                 |
//                       + -> TSPaid -------- + ----------------- +
//                       |                    ^                   |
//                       |                    |                   |
//                       + -> TSCompleted --- +                   |
//                                                                |
// REFUND : TSInProgress -> TSRefund ---------------------------- +
//              |              ^
//              |              |
//              + ------- > TSRefundError <- + (max counts 5)
//                             |             |
//                             + ----------- +
//
// ===================================================================

// Unpid state, can be use as default trade state for generate
// a trade ticket, as status machine it can change to :
//
//	TSUnpaid	 -> TSRevoked   : canceld parment
//				 -> TSPayError  : pay error
//				 -> TSPaid      : success paid
//				 -> TSCompleted : only for dividing payment, to mark dividing completed
const TSUnpaid = "UNPAID"

// Pay error state, as status machine it can change to :
//
//	TSPayError	 -> TSPaid    : success paid
//				 -> TSRevoked : canceld parment
//				 -> self (over 5 times) -> TSClosed
const TSPayError = "PAY_ERROR"

// Revoked state, cancel by user, as status machine it only
// can be changed to closed state:
//
//	TSRevoked	 -> TSClosed : close the trade ticket
const TSRevoked = "REVOKED"

// Paied success state, as status machine it only can be
// changed to closed state :
//
//	TSPaid		 -> TSClosed : close the trade ticket
//
// `WARNING` :
//
// the refund action will generate a new trade ticket and set
// TSInProgress as default.
const TSPaid = "PAID"

// Completed all dividing payments, as status machine it only
// can be changed to closed state :
//
//	TSCompleted	 -> TSClosed : close the trade ticket
const TSCompleted = "COMPLETED"

// Refund in progress state, use as default trade state when generate
// a refund ticket, as status machine it can change to :
//
//	TSInProgress -> TSRefund      : success refund
//				 -> TSRefundError : refund error
const TSInProgress = "REFUND_IN_PROGRESS"

// Refund success state, as status machine it only can be
// changed to closed state :
//
//	TSRefundError	-> TSClosed : close the trade ticket
//					-> self (over 5 times) -> TSClosed
const TSRefundError = "REFUND_ERROR"

// Refund success state, as status machine it only can be
// changed to closed state :
//
//	TSRefund	 -> TSClosed : close the trade ticket
const TSRefund = "REFUND"

// Closed state, as the last state of status machine.
const TSClosed = "CLOSED"
