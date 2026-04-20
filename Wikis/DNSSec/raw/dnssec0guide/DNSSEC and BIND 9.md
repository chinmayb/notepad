---
title: DNSSEC and BIND 9
source: https://www.isc.org/dnssec/
author:
  - "[[Internet Systems Consortium]]"
published: 2025-12-30
created: 2026-04-13
description: BIND 9 fully supports DNSSEC and we encourage the use of DNSSEC as a best practice
tags:
  - clippings
ingested:
---
In addition to verifying the integrity of your zone data, the DNSSEC chain of trust can also be used to associate other information with your domain, such as PGP keys, to help improve the security of applications.

#### What is DNSSEC?

DNSSEC (Domain Name System Security Extensions) adds resource records and message header bits which can be used to verify that the requested data matches what the zone administrator put in the zone and has not been altered in transit. DNSSEC doesn’t provide a secure tunnel; it doesn’t encrypt or hide DNS data. It was designed with backwards compatibility in mind. The original standard DNS protocol continues to work the same.

The resource record types are: RRSIG (for digital signature), DNSKEY (the public key), DS (Delegation Signer), and NSEC (pointer to next secure record). The message header bits are: AD (for authenticated data) and CD (checking disabled). A DNSSEC validating resolver uses these records and public key (asymmetric) cryptography to prove the integrity of the DNS data. A private key (specific to a zone) is used to encrypt a hash of a set of resource records; this is the digital signature stored in a RRSIG record.

The corresponding public key is stored in the DNSKEY resource record. The validating resolver uses that DNSKEY to decrypt the RRSIG and then compares the result with the hash of the corresponding resource record set to verify it has not changed. A hash of the public DNSKEY is stored in a DS record. This is stored in the parent zone. The validating resolver retrieves from the parent the DS record and its corresponding signature (RRSIG) and public key (DNSKEY); a hash of that public key is available from its parent. This becomes a chain of trust, also called an authentication chain. The validating resolver is configured with a trust anchor; this is the starting point which refers to a signed zone. The trust anchor is a DNSKEY or DS record and should be securely retrieved from a trusted source (not using DNS).

All the names in the zone have corresponding NSEC records listed in order, creating a chain of all the signed record sets. (Corresponding RRSIG records are also created to verify the NSEC data.) Because there is no gap, NSEC records are used to provide proof of non-existence of an resource record or to authenticate negative replies.

While the basic theory behind DNSSEC hasn’t changed in over a decade, the implementation in BIND 9 has changed quite a bit in recent years, as we have implemented a modern [Key and Signing Policy (KASP)](https://bind9.readthedocs.io/en/latest/chapter5.html#fully-automated-key-and-signing-policy) function, to simplify key and signature management.

### Current on-line documentation on DNSSEC in BIND

#### Preparing for DNSSEC

- Use resolvers that are DNSSEC-capable and configured to do the validation. All versions of BIND 9 are DNSSEC-capable.
- Make sure network devices don’t lose or stop EDNS0 (Extension Mechanisms for DNS) or squash DNSSEC-related traffic. DNSSEC requires EDNS0 to support the larger DNS message sizes and for the DNSSEC OK (DO) EDNS header bit.
- DNSSEC does increase DNS traffic, with more requests and larger responses. For high-volume DNS traffic, prepare for increased bandwidth needs.
- DNSSEC is more sensitive to time issues (i.e. system clocks being really out of sync) than plain DNS; make sure your system clocks are reasonably accurate.
- If hosting DNSSEC-signed zones, make sure your secondaries also support it and have DNSSEC enabled.

#### DNSSEC presentations

- Presentation on the Encrypted DNS call [DNSSEC Multi-Signer Model](https://419.consulting/encrypted-dns/f/the-dnssec-multi-signer-model) - Matthijs Mekking
- Presentation at DNS-OARC41 [Shorter DNSSEC Outages](https://www.isc.org/docs/2023-oarc47-spacek.pdf) - Petr Špaček
- Presentation at ICANN 77 - [DS Automation in BIND 9](https://icann.zoom.us/rec/play/MLZxoabslLzAo5_RJ4wvbWTAB3y4c1-GnskwJac-nqnsZWin-5YUWBT24mlSmS2F7DGwu82fB-2hM2Lv.wBdVDRihaqF1_YDf?startTime%3D1686598390000) - Matthijs Mekking
- Presentation at RIPE 86 [DNSSEC Multi-signer Model in BIND](https://ripe86.ripe.net/archives/video/1079/) - Matthijs Mekking
- Presentation at DNS-OARC [Performance effects of DNSSEC validation on a busy resolver](https://www.youtube.com/watch?t=9847&v=SRw_EHOTWyY&feature=youtu.be) - Petr Špaček

#### Guides and training

- ISC KB: [DNSSEC Key and Signing Policy](https://kb.isc.org/docs/dnssec-key-and-signing-policy)
- ISC KB: [DNSSEC Key States](https://kb.isc.org/docs/dnssec-key-states)
- [DNSVIZ](https://dnsviz.net/) remains the single most useful diagnostic tool for DNSSEC
- The [BIND 9 Administrator Reference Manual (ARM)](https://bind9.readthedocs.io/en/latest/index.html) on ReadTheDocs. Choose the version of BIND 9 that you are running or plan to run and learn about the DNSSEC Key Manager utility, dnssec-keygen, dnssec-signzone, and other DNSSEC features built into BIND 9.
- [ISOC 360 DNSSEC page](https://www.internetsociety.org/deploy360/dnssec/)
- [DELV](https://kb.isc.org/docs/aa-01152), a BIND tool that checks the DNSSEC validation chain

*Older, possibly still useful references*

- [10-part 2020 webinar series](https://youtu.be/L-IXvfp7BdY) on implementing DNSSEC with BIND, covering theory, configuration, operation and troubleshooting. Note that this pre-dated the current [Key and Signing Policy feature (KASP)](https://bind9.readthedocs.io/en/latest/chapter5.html#dnssec-kasp)
- The [BIND 9 DNSSEC Guide](https://bind9.readthedocs.io/en/latest/dnssec-guide.html), now an appendix to the BIND Administrative Reference Guide, is a long-format guide that explains a lot of the background and the “why” of DNSSEC, as well as the “how.” It includes the [KASP (Key and Signing Policy)](https://bind9.readthedocs.io/en/latest/chapter5.html#dnssec-kasp) tool.
- [State of DNSSEC Deployment 2016](https://www.internetsociety.org/resources/doc/2016/state-of-dnssec-deployment-2016/) (report from the Internet Society)
- [“Deploying DNSSEC with BIND 9.7”](https://downloads.isc.org/isc/pubs/pres/NANOG/50/DNSSEC-NANOG50.pdf) given by Alan Clegg at NANOG50
- [DNSSEC goes Mainstream: Deployment Incentives, Experience, Questions](https://www.nanog.org/meetings/nanog46/presentations/Tuesday/Woolf_dnssec_N46.pdf), a presentation by Suzanne Woolf at NANOG, June 2009.