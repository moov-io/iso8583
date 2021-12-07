package iso8583

// MesssageTypeIndicator message type indicator is a four-digit numeric field which indicates the overall function of the ISO 8583:1987 message
type MesssageTypeIndicator string

const (
	// AuthorizationRequest is a request from a point-of-sale terminal for authorization for a cardholder purchase
	AuthorizationRequest MesssageTypeIndicator = "0100"

	// AuthorizationResponse is a request response to a point-of-sale terminal for authorization for a cardholder purchase
	AuthorizationResponse MesssageTypeIndicator = "0110"

	// AuthorizationAdvice is when the point-of-sale device breaks down and you have to sign a voucher
	AuthorizationAdvice MesssageTypeIndicator = "0120"

	// AuthorizationAdviceRepeat used to repeat if the advice times out
	AuthorizationAdviceRepeat MesssageTypeIndicator = "0121"

	// IssuerResponseToAuthorizationAdvice is a confirmation of receipt of authorization advice
	IssuerResponseToAuthorizationAdvice MesssageTypeIndicator = "0130"

	// AuthorizationPositiveAcknowledgement indicates that an Authorization Response was received
	AuthorizationPositiveAcknowledgement MesssageTypeIndicator = "0180"

	// AuthorizationNegativeAcknowledgement indicates that an Authorization Response or Reversal Response was late or invalid
	AuthorizationNegativeAcknowledgement MesssageTypeIndicator = "0190"

	// AcquirerFinancialRequest is a request for funds, typically from an ATM or pinned point-of-sale device
	AcquirerFinancialRequest MesssageTypeIndicator = "0200"

	// IssuerResponseToFinancialRequest is a issuer response to request for funds
	IssuerResponseToFinancialRequest MesssageTypeIndicator = "0210"

	// AcquirerFinancialAdvice is used to complete transaction initiated with authorization request. e.g. Checkout at a hotel.
	AcquirerFinancialAdvice MesssageTypeIndicator = "0220"

	// AcquirerFinancialAdviceRepeat is used if the advice times out
	AcquirerFinancialAdviceRepeat MesssageTypeIndicator = "0221"

	// IssuerResponseToFinancialAdvice is a confirmation of receipt of financial advice
	IssuerResponseToFinancialAdvice MesssageTypeIndicator = "0230"

	// BatchUpload is a file update/transfer advice
	BatchUpload MesssageTypeIndicator = "0320"

	// BatchUploadResponse is a file update/transfer advice response
	BatchUploadResponse MesssageTypeIndicator = "0330"

	// AcquirerReversalRequest is used to reverse a transaction
	AcquirerReversalRequest MesssageTypeIndicator = "0400"

	// AcquirerReversalResponse is a response to a reversal request
	AcquirerReversalResponse MesssageTypeIndicator = "0410"

	// AcquirerReversalAdvice
	AcquirerReversalAdvice MesssageTypeIndicator = "0420"

	// AcquirerReversalAdviceResponse
	AcquirerReversalAdviceResponse MesssageTypeIndicator = "0430"

	// BatchSettlementResponse is a card acceptor reconciliation request response
	BatchSettlementResponse MesssageTypeIndicator = "0510"

	// AdministrativeRequest is a message delivering administrative data, often free-form and potentially indicating a failure message
	AdministrativeRequest MesssageTypeIndicator = "0600"

	// AdministrativeResponse is a response to an administrative request
	AdministrativeResponse MesssageTypeIndicator = "0610"

	// AdministrativeAdvice is an administrative request with stronger delivery guarantees
	AdministrativeAdvice MesssageTypeIndicator = "0620"

	// AdministrativeAdviceResponse is a response to an administrative advice
	AdministrativeAdviceResponse MesssageTypeIndicator = "0630"

	// NetworkManagementRequest is used in hypercom terminals initialize request. Echo test, logon, logoff etc
	NetworkManagementRequest MesssageTypeIndicator = "0800"

	// NetworkManagementResponse is a hypercom terminals initialize response. Echo test, logon, logoff etc.
	NetworkManagementResponse MesssageTypeIndicator = "0810"

	// NetworkManagementAdvice is a key change
	NetworkManagementAdvice MesssageTypeIndicator = "0820"
)
