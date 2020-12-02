
// protocole name

NAME= Base ISO8583 (1987) example

// protocol encoding definition

MTI_Enc= CHAR
BMP_Enc= HEX
LEN_Enc= CHAR
NUM_Enc= CHAR
CHR_Enc= ASCII
BIN_Enc= NONE
TRK_Enc= CHAR

// data elements definition
// field number, type, size, size len, subfield, description

FIELD=   2;    N;     19; 2; N; Primary Account Number
FIELD=   3;   AN;      6; 0; N; Processing Code
FIELD=   4;    N;     12; 0; N; Transaction Amount
FIELD=   7;    N;     10; 0; N; Transmission date & time
FIELD=  11;    N;      6; 0; N; Systems trace audit number
FIELD=  12;    N;      6; 0; N; Time, Local Tranaction
FIELD=  13;    N;      4; 0; N; Date, local Tranaction
FIELD=  14;    N;      4; 0; N; Expiration date
FIELD=  15;    N;      4; 0; N; Settlement date
FIELD=  18;    N;      4; 0; N; Merchant type, or merchant category code
FIELD=  32;    N;     11; 2; N; Acquiring institution identification code
FIELD=  35;    Z;     37; 2; N; Track 2 data
FIELD=  37;   AN;     12; 0; N; Retrieval reference number
FIELD=  39;   AN;      2; 0; N; Response code
FIELD=  42;  ANS;     15; 0; N; Card acceptor identification code
FIELD=  48;  ANS;    999; 3; N; Additional Data - Private
FIELD=  49;    N;      3; 0; N; Currency Code
FIELD=  55;  ANS;    999; 3; N; ICC data – EMV having multiple tags
FIELD=  60;  ANS;    999; 3; N; Reserved (national)
FIELD=  63;  ANS;    999; 3; N; Reserved (private)
FIELD=  70;    N;      3; 0; N; Network management information code
FIELD=  73;    N;      6; 0; N; Action date
FIELD=  90;    N;     42; 0; N; Original Data Elements

// message definition
// MTI, mandatory et optional masks, description

MSG=0100;722000000000000000000000000000000000000000000000;000800010000000000000000000000000000000000000000;Authorization request
MSG=0110;722000010200000000000000000000000000000000000000;000000000000000000000000000000000000000000000000;Authorization request response
MSG=0420;722000010000000000000000000000000000000000000000;000000000000000000000000000000000000000000000000;Reversal advice
MSG=0440;722000010200000000000000000000000000000000000000;000000000000000000000000000000000000000000000000;Reversal notification
MSG=0450;722000000000000000000000000000000000000000000000;000000000000000000000000000000000000000000000000;Reversal notification acknowledgement
