# Message Type Indicator

An ISO 8583 message starts with four numeric digits that define general information about the transaction, called a Message Type Indicator (MTI). Each digit defines a different piece of information, and the combination of all four describes the message in great detail. From an MTI alone, you can determine the ISO 8583 version, message class, message function, and message origin.

## Version

The **first digit** indicates the ISO 8583 version. There are currently three versions: 1987, 1993, and 2007 (with 2007 being the least common).

| Code (#xxx) | Meaning         |
|-------------|-----------------|
| 0           | ISO 8583:1987   |
| 1           | ISO 8583:1993   |
| 2           | ISO 8583:2003   |
| 3-7         | Reserved by ISO |
| 8           | National use    |
| 9           | Private use     |

## Message Class

The **second digit** indicates the overall message classification, or purpose.

| Code (x#xx) | Meaning                |
|-------------|------------------------|
| 0           | Reserved by ISO        |
| 1           | Authorization          |
| 2           | Financial              |
| 3           | File Actions           |
| 4           | Reversal or Chargeback |
| 5           | Reconciliation         |
| 6           | Administrative         |
| 7           | Fee Collection         |
| 8           | Network Management     |
| 9           | Reserved by ISO        |

## Message Function

The **third digit** defines how the message should flow within the payment system. Requests are end-to-end, while advices, notifications, and instructions are point-to-point. A request may be rejected by the reciever (the issuer), while advices and notifications must be accepted by their receiver. It should be noted values `6` and `7` are only used in ISO8583:2007.

| Code (xx#x) | Meaning                      |
|-------------|------------------------------|
| 0           | Request (from acquirer to issuer)  |
| 1           | Request Response             |
| 2           | Advice                       |
| 3           | Advice Response              |
| 4           | Notification                 |
| 5           | Notification Acknowledgement |
| 6           | Instruction                  |
| 7           | Instruction Acknowledgement  |
| 8-9         | Reserved by ISO              |

## Message Origin

The **fourth digit** indicates the message's source within the payment chain. Repeats occur after timeouts.

| Code (xxx#) | Meaning         |
|-------------|-----------------|
| 0           | Acquirer        |
| 1           | Acquirer Repeat |
| 2           | Issuer          |
| 3           | Issuer Repeat   |
| 4           | Other           |
| 5           | Other Repeat    |
| 6-9         | Reserved by ISO |

## Summary

Deciphering an MTI code is very straightforward. For code `1100`, you can derive the following:

| Digit | Value | Meaning       |
|-------|-------|---------------|
| 1     | 1     | Version 1993  |
| 2     | 1     | Authorization |
| 3     | 0     | Request       |
| 4     | 0     | Acquirer      |

Thus, `1100` is an authorization request originating from the acquirer and using the 1993 version of ISO 8583.

There are some codes that are invalid under normal circumstances. For example, `1102` would represent an authorization request originating from an issuer, which isn't possible. Certain codes are also significantly more common than others — check out this [reference guide for common MTI codes and their standard meanings](http://www.fintrnmsgtool.com/iso-mti-code.html).
