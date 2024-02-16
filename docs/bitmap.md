# Bitmaps

A bitmap is a sequence of bits that represents whether or not potential fields after the bitmap are present. The location of these fields in the bitmap correspond to their location in the overall message. A bit set to `1` indicates that a field is present, while a bit set to `0` indicates that a field is not present. For example, a bitmap of `00011001` means that fields 4, 5, and 8 are present, while fields 1, 2, 3, 6, and 7 are not.

An ISO 8583 message contains at least one bitmap, called the "primary bitmap", which immediately follows the Message Type Indicator (MTI). The primary bitmap is 8 bytes long and provides indicators for fields 1-64. We can think of the primary bitmap as a mandatory "field 0". There may also be a secondary bitmap at field 1 (the first field following the primary bitmap). This secondary bitmap provides indicators for fields 65-128. A tertiary bitmap may exist, but this is very rare.

Bitmaps are often represented by hex characters. For example, `0x4210000000000000` in hex corresponds to fields 2, 7, and 12 being present. The equivalent binary bitmap would be `0b0100001000010000000000000000000000000000000000000000000000000000`.

Similarly, let's look at an example with a secondary bitmap. Say the primary bitmap is `0xF000000000000000` and the secondary bitmap is `0x3000000000000000`. We can conclude that fields 1 (the secondary bitmap itself), 2, 3, 4, 67, and 68 are present.
