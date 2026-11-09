# External resources

A curated list of third-party material for learning ISO 8583 and working with
this library: concept guides, standards references, and tools.

> **Disclaimer**
>
> The resources below are provided by third parties. moov-io does **not**
> control, endorse, or take responsibility for their content, availability, or
> any products or services they may offer. A site being listed here is not a
> recommendation of anything it sells. Links can change or disappear over time.
> If a link is broken, outdated, or no longer appropriate, please
> [open an issue](https://github.com/moov-io/iso8583/issues/new) or a PR.

## What gets listed here

These are guidelines, not hard rules — judgment still applies. A resource is a
good fit when most of the following hold:

- **It helps someone working with ISO 8583 or this library.** Concepts, spec
  details, worked examples, or tooling a reader can actually use.
- **Closer to this package ranks higher.** We prioritize resources that use or
  explain `moov-io/iso8583` or related packages such as
  [moov-io/iso8583-connection](https://github.com/moov-io/iso8583-connection).
  General ISO 8583 material is welcome, but package-specific guidance comes
  first, and moov-io's own resources lead each list.
- **The useful part is free to reach.** No paywall, no signup, no "the rest of
  the article after you buy." A commercial site is fine; gating the content that
  made us link it is not.
- **It is accurate and maintained.** We remove links that rot or go stale.
- **It stands on its own merit.** We list genuinely useful material, not pages
  whose main purpose is SEO or lead generation that happens to mention ISO 8583.
- **Affiliation is disclosed.** If you wrote the resource, run the site, or your
  employer does, say so in the PR. Disclosed self-submissions are welcome — this
  is about honesty, not exclusion.

We list; we do not endorse. Meeting these guidelines gets a link considered, not
guaranteed — maintainers make the final call.

## Guides and articles

- [Mastering ISO 8583 messages with Golang](https://alovak.com/2024/08/15/mastering-iso-8583-messages-with-golang/)
- [Mastering ISO 8583 Message Networking with Golang](https://alovak.com/2024/08/27/mastering-iso-8583-message-networking-with-golang/)

## Standards and references

- [ISO 8583 Terms and Definitions](https://www.iso.org/obp/ui/#iso:std:iso:8583:-1:ed-1:v1:en)

## Tools

- [Annotated ISO 8583 message examples](https://iso8583parser.com/en/articles/iso8583-message-examples) — worked examples that open pre-decoded in a browser-based parser; handy for sanity-checking your own `Pack()` output against an independent decoder.
