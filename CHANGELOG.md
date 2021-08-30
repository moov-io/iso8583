## v0.5.0 (Released 2021-08-30)

BREAKING CHANGES

- refactor field.Spec and Composite field type to support string tag encoding
- amend prefixer interface and implement BER-TLV encoder and prefixer
- amend prefixers and fields to accommodate for change in Prefixer interface

ADDITIONS

- implement bertlv prefixer and encoding/bertlv tests
