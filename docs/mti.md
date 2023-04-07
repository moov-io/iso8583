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

There are some codes that are invalid under normal circumstances. For example, `1102` would represent an authorization request originating from an issuer, which isn't possible. Certain codes are also significantly more common than others — check out this [reference guide for common MTI codes and their standard meanings](https://www.liquisearch.com/iso_8583/message_type_indicator/examples).

## Example MTI Codes

The different ISO8583 versions have changed MTI codes over the years. Here's a breakdown of the change over time.

| 1987 | 1993 | 2003 | Meaning | Usage
|----|----|----|----|----|
| 0100 | 1100 | 2100 | Authorization Request | Acquirer Send Authorization Request Message to Issuer |
| 0110 | 1110 | 2110 | Authorization Request Response | Issuer Send Authorization Request Response Message back to Acquirer |
| 0120 | 1120 | 2120 | Authorization Advice | Acquirer send Authorization Advice Request Message to Issuer |
| 0121 | 1121 | 2121 | Authorization Advice Repeat | Acquirer send Authorization Advice Repeat Message to Issuer |
| 0130 | 1130 | 2130 | Authorization Advice Response | Issuer send Authorization Advice Response Message back to Acquirer |
| 0200 | 1200 | 2200 | Financial Transaction Request | Acquirer Send Financial Transaction Request Message to Issuer |
| 0210 | 1210 | 2210 | Financial Transaction Request Response | Issuer Send Financial Transaction Response Message back to Acquirer |
| 0220 | 1220 | 2220 | Financial Transaction Advice | Acquirer Send Financial Transaction Advice Request Message to Issuer |
| 0221 | 1221 | 2221 | Financial Transaction Advice Repeat | Acquirer Send Financial Transaction Advice Repeat Message to Issuer |
| 0230 | 1230 | 2230 | Financial Transaction Advice Response | Issuer Send Financial Transaction Advice Repeat Response Message Back to Acquirer |
| 0320 | 1320 | 2320 | Batch Upload Request | Acquirer send Batch Upload Request Message to Issuer |
| 0330 | 1330 | 2330 | Batch Upload Response | Issuer send Batch Upload Response Message to Acquirer |
| 0400 | 1400 | 2400 | Acquirer Reversal Request | Acquirer Send Reversal Request Message to Issuer |
| 0402 | 1402 | 2402 | Card Issuer Reversal Request | Card Issuer Send Reversal Request Message to Acquirer |
| 0410 | 1410 | 2410 | Acquirer Reversal Request Response | Acquirer send reversal request response back to Issuer |
| 0412 | 1412 | 2412 | Card Issuer Reversal Request Response | Card Issuer Send Reversal Request Response Message back to Acquirer |
| 0420 | 1420 | 2420 | Acquirer Reversal Advice | Acquirer Send Reversal Advice message to Card Issuer |
| 0421 | 1421 | 2421 | Acquirer Reversal Advice Repeat | Acquirer Send Reversal Advice Repeat message to Card Issuer |
| 0430 | 1430 | 2430 | Acquirer Reversal Advice Response | Acquirer Send Reversal Advice Response message back to Card Issuer |
| 0432 | 1432 | 2432 | Card Issuer Reversal Advice Response | Card Issuer Send Reversal Advice Response message back to Acquirer |
| 0500 | 1500 | 2500 | Acquirer Reconciliation Request | Acquirer Send Reconciliation Request Message to Card Issuer |
| 0502 | 1502 | 2502 | Issuer Reconciliation Request | Issuer Send Reconciliation Request Message to Acquirer |
| 0510 | 1510 | 2510 | Acquirer Reconciliation Request Response | Acquirer send Reconciliation Request Response message back to Card Issuer |
| 0512 | 1512 | 2512 | Issuer Reconciliation Request Response | Issuer Send Reconciliation Request Response Message back to Acquirer |
| 0520 | 1520 | 2520 | Acquirer Reconciliation Advice | Acquirer send Reconciliation Advice message to Card Issuer |
| 0521 | 1521 | 2521 | Acquirer Reconciliation Advice Repeat | Acquirer send Reconciliation Advice Repeat message to Card Issuer |
| 0522 | 1522 | 2522 | Issuer Reconciliation Advice | Issuer send Reconciliation Advice message to Acquirer |
| 0523 | 1523 | 2523 | Issuer Reconciliation Advice Repeat | Issuer send Reconciliation Advice Repeat message to Acquirer |
| 0530 | 1530 | 2530 | Acquirer Reconciliation Advice Response | Acquirer send Reconciliation Advice Response message to Card Issuer |
| 0532 | 1532 | 2532 | Issuer Reconciliation Advice Response | Issuer send Reconciliation Advice Response message to Acquirer |
| 0604 | 1604 | 2604 | Administrative Request | Acquirer/Card Issuer Send Administrative Request Message |
| 0605 | 1605 | 2605 | Administrative Request Repeat | Acquirer/Card Issuer Send Administrative Request Repeat Message |
| 0614 | 1614 | 2614 | Administrative Request Response | Acquirer/Card Issuer Send Administrative Request Response Message |
| 0624 | 1624 | 2624 | Administrative Advice | Acquirer/Card Issuer Send Administrative Advice Message |
| 0625 | 1625 | 2625 | Administrative Advice Repeat | Acquirer/Card Issuer Send Administrative Advice Repeat Message |
| 0634 | 1634 | 2634 | Administrative Advice Response | Acquirer/Card Issuer Send Administrative Advice Response Message |
| 0644 | 1644 | 2644 | Administrative Notification | Acquirer/Card Issuer Send Administrative Notification Message |
| 0720 | 1720 | 2720 | Acquirer Fee Collection Advice | Acquirer Send Fee Collection Advice Message to Issuer |
| 0721 | 1721 | 2721 | Acquirer Fee Collection Advice Repeat | Acquirer Send Fee Collection Advice Repeat Message to Issuer |
| 0722 | 1722 | 2722 | Issuer Fee Collection Advice | Card Issuer Send Fee Collection Advice Message to Acquirer |
| 0723 | 1723 | 2723 | Issuer Fee Collection Advice Repeat | Card Issuer Send Fee Collection Advice Repeat Message to Acquirer |
| 0730 | 1730 | 2730 | Acquirer Fee Collection Advice response | Acquirer Send Fee Collection Advice Response Message to Card Issuer |
| 0732 | 1732 | 2732 | Issuer Fee Collection Advice response | Card Issuer Send Fee Collection Advice Response Message to Acquirer |
| 0740 | 1740 | 2740 | Acquirer Fee Collection Notification | Acquirer Send Fee Collection Notification Message to Card Issuer |
| 0742 | 1742 | 2742 | Issuer Fee Collection Notification | Card Issuer Send Fee Collection Notification Message to Acquirer |
| 0800 | 1800 | 2800 | Network/Key Management Request | Acquirer/Card Issuer Send Network/Key Management Request Message |
| 0810 | 1810 | 2810 | Network/Key Management Response | Acquirer/Card Issuer Send Network/Key Management Response Message |
| 0820 | 1820 | 2820 | Network Management Advice | Acquirer/Card Issuer Send Network Management Advice Message |
| 0821 | 1821 | 2821 | Network Management Advice Repeat | Acquirer/Card Issuer Send Network Management Advice Repeat Message |
| 0830 | 1830 | 2830 | Network Management Advice Response | Acquirer/Card Issuer Send Network Management Advice Response Message |
