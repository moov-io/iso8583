# Data Elements

The bulk of an ISO 8583 message consists of various data elements. Some are general purpose, while others may be specific to certain countries or systems. The maximum lengths and overall collection of data fields varies from version to version, with the 1987 standard supporting up to 128 data elements. A data type is commonly defined by a value type, length type, and maximum length value.

## Value Types

ISO 8583 often uses various abbreviations to represent data value types. The table below lists all possible abbreviations and their definitions.

| Abbreviation | Meaning                                                                                                          |
|--------------|------------------------------------------------------------------------------------------------------------------|
| a            | Alpha, including blanks                                                                                          |
| n            | Numeric                                                                                                          |
| s            | Special characters                                                                                               |
| an           | Alphanumeric                                                                                                     |
| as           | Alpha and special characters                                                                                     |
| ns           | Numeric and special characters                                                                                   |
| ans          | Alphanumeric and special characters                                                                              |
| b            | Binary                                                                                                           |
| x+n          | First byte is either 'C' (positive or credit value) or 'D' (negative or debit value), followed by numeric digits |
| z            | Tracks 2 and 3 code set as defined in ISO/IEC 7813 and ISO/IEC 4909 respectively                                 |

## Length Types

A field may be fixed or variable length. ISO 8583 uses 'LVAR' notation to indicate variable lengths, where each 'L' represents a digit for the max length. For example, LLVAR allows for a max length up to 99 digits and LLLVAR allows a max length up to 999 digits. You may see dots '..' that are equivalent to the number of 'L's when looking at data field abbreviations.

## Data Type Examples

| Abbreviation | Meaning                                                  |
|--------------|----------------------------------------------------------|
| a 8          | Alpha with fixed length of 8 digits                      |
| an..35       | Alphanumeric with LLVAR variable length of 35 digits max |
| b...999      | Binary with LLLVAR variable length of 999 digits max     |
