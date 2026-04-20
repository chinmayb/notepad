---
title: "RFC 9364: DNS Security Extensions (DNSSEC)"
source: https://datatracker.ietf.org/doc/html/rfc9364
author:
  - "[[P. Hoffman ICANN]]"
published: 2023-02-01
created: 2026-04-13
description: DNS Security Extensions (DNSSEC) (RFC 9364, )
tags:
  - clippings
---




| RFC 9364 | DNSSEC                | February 2023 |
| -------- | --------------------- | ------------- |
| Hoffman  | Best Current Practice | \[Page\]      |

## RFC 9364

## Abstract

This document describes the DNS Security Extensions (commonly called "DNSSEC") that are specified in RFCs 4033, 4034, and 4035, as well as a handful of others. One purpose is to introduce all of the RFCs in one place so that the reader can understand the many aspects of DNSSEC. This document does not update any of those RFCs. A second purpose is to state that using DNSSEC for origin authentication of DNS data is the best current practice. A third purpose is to provide a single reference for other documents that want to refer to DNSSEC.[¶](#section-abstract-1)

## Status of This Memo

This memo documents an Internet Best Current Practice.[¶](#section-boilerplate.1-1)

This document is a product of the Internet Engineering Task Force (IETF). It represents the consensus of the IETF community. It has received public review and has been approved for publication by the Internet Engineering Steering Group (IESG). Further information on BCPs is available in Section 2 of RFC 7841.[¶](#section-boilerplate.1-2)

Information about the current status of this document, any errata, and how to provide feedback on it may be obtained at [https://www.rfc-editor.org/info/rfc9364](https://www.rfc-editor.org/info/rfc9364).[¶](#section-boilerplate.1-3)

## 1.

The core specification for what we know as DNSSEC (the combination of \[[RFC4033](https://datatracker.ietf.org/doc/html/rfc4033)\], \[[RFC4034](https://datatracker.ietf.org/doc/html/rfc4034)\], and \[[RFC4035](https://datatracker.ietf.org/doc/html/rfc4035)\]) describes a set of protocols that provide origin authentication of DNS data. \[[RFC6840](https://datatracker.ietf.org/doc/html/rfc6840)\] updates and extends those core RFCs but does not fundamentally change the way that DNSSEC works.[¶](#section-1-1)

This document lists RFCs that should be considered by someone creating an implementation of, or someone deploying, DNSSEC as it is currently standardized. Although an effort was made to be thorough, the reader should not assume this list is comprehensive. It uses terminology from those documents without defining that terminology. It also points to the relevant IANA registry groups that relate to DNSSEC. It does not, however, point to standards that rely on zones needing to be signed by DNSSEC, such as DNS-Based Authentication of Named Entities (DANE) \[[RFC6698](https://datatracker.ietf.org/doc/html/rfc6698)\].[¶](#section-1-2)

### 1.1.

Using the DNSSEC set of protocols is the best current practice for adding origin authentication of DNS data. To date, no Standards Track RFCs offer any other method for such origin authentication of data in the DNS.[¶](#section-1.1-1)

More than 15 years after the DNSSEC specification was published, it is still not widely deployed. Recent estimates are that fewer than 10% of the domain names used for websites are signed, and only around a third of queries to recursive resolvers are validated. However, this low level of deployment does not affect whether using DNSSEC is a best current practice; it just indicates that the value of deploying DNSSEC is often considered lower than the cost. Nonetheless, the significant deployment of DNSSEC beneath some top-level domains (TLDs) and the near-universal deployment of DNSSEC for the TLDs in the DNS root zone demonstrate that DNSSEC is applicable for implementation by both ordinary and highly sophisticated domain owners.[¶](#section-1.1-2)

### 1.2.

Developers of validating resolvers and authoritative servers, as well as operators of validating resolvers and authoritative servers, need to know the parts of the DNSSEC protocol that would affect them. They should read the DNSSEC core documents and probably at least be familiar with the extensions. Developers will probably need to be very familiar with the algorithm documents as well.[¶](#section-1.2-1)

As a side note, some of the DNSSEC-related RFCs have significant errata, so reading the RFCs should also include looking for the related errata.[¶](#section-1.2-2)

## 2.

What we refer to as "DNSSEC" is the third iteration of the DNSSEC specification; \[[RFC2065](https://datatracker.ietf.org/doc/html/rfc2065)\] was the first, and \[[RFC2535](https://datatracker.ietf.org/doc/html/rfc2535)\] was the second. Earlier iterations have not been deployed on a significant scale. Throughout this document, "DNSSEC" means the protocol initially defined in \[[RFC4033](https://datatracker.ietf.org/doc/html/rfc4033)\], \[[RFC4034](https://datatracker.ietf.org/doc/html/rfc4034)\], and \[[RFC4035](https://datatracker.ietf.org/doc/html/rfc4035)\].[¶](#section-2-1)

The three initial core documents generally cover different topics:[¶](#section-2-2)

- \[[RFC4033](https://datatracker.ietf.org/doc/html/rfc4033)\] is an overview of DNSSEC, including how it might change the resolution of DNS queries.[¶](#section-2-3.1)
- \[[RFC4034](https://datatracker.ietf.org/doc/html/rfc4034)\] specifies the DNS resource records used in DNSSEC. It obsoletes many RFCs about earlier versions of DNSSEC.[¶](#section-2-3.2)
- \[[RFC4035](https://datatracker.ietf.org/doc/html/rfc4035)\] covers the modifications to the DNS protocol incurred by DNSSEC. These include signing zones, serving signed zones, resolving in light of DNSSEC, and authenticating DNSSEC-signed data.[¶](#section-2-3.3)

At the time this set of core documents was published, someone could create a DNSSEC implementation of signing software, of a DNSSEC-aware authoritative server, and/or of a DNSSEC-aware recursive resolver from the three core documents, plus a few older RFCs specifying the cryptography used. Those two older documents are the following:[¶](#section-2-4)

- \[[RFC2536](https://datatracker.ietf.org/doc/html/rfc2536)\] defines how to use the DSA signature algorithm (although it refers to other documents for the details). DSA was thinly implemented and can safely be ignored by DNSSEC implementations.[¶](#section-2-5.1)
- \[[RFC3110](https://datatracker.ietf.org/doc/html/rfc3110)\] defines how to use the RSA signature algorithm (although refers to other documents for the details). RSA is still among the most popular signing algorithms for DNSSEC.[¶](#section-2-5.2)

It is important to note that later RFCs update the core documents. As just one example, \[[RFC9077](https://datatracker.ietf.org/doc/html/rfc9077)\] changes how TTL values are calculated in DNSSEC processing.[¶](#section-2-6)

### 2.1.

As with any major protocol, developers and operators discovered issues with the original core documents over the years. \[[RFC6840](https://datatracker.ietf.org/doc/html/rfc6840)\] is an omnibus update to the original core documents and thus itself has become a core document. In addition to covering new requirements from new DNSSEC RFCs, it describes many important security and interoperability issues that arose during the deployment of the initial specifications, particularly after the DNS root was signed in 2010. It also lists some errors in the examples of the core specifications.[¶](#section-2.1-1)

\[[RFC6840](https://datatracker.ietf.org/doc/html/rfc6840)\] brings a few additions into the core of DNSSEC. It makes NSEC3 \[[RFC5155](https://datatracker.ietf.org/doc/html/rfc5155)\] as much a part of DNSSEC as NSEC is. It also makes the SHA-256 and SHA-512 hash functions defined in \[[RFC4509](https://datatracker.ietf.org/doc/html/rfc4509)\] and \[[RFC5702](https://datatracker.ietf.org/doc/html/rfc5702)\] part of the core.[¶](#section-2.1-2)

## 3.

Current cryptographic algorithms typically weaken over time as computing power improves and new cryptoanalysis emerges. Two new signing algorithms have been adopted by the DNSSEC community: Elliptic Curve Digital Signature Algorithm (ECDSA) \[[RFC6605](https://datatracker.ietf.org/doc/html/rfc6605)\] and Edwards-curve Digital Signature Algorithm (EdDSA) \[[RFC8080](https://datatracker.ietf.org/doc/html/rfc8080)\]. ECDSA and EdDSA have become very popular signing algorithms in recent years. The GOST signing algorithm \[[GOST-SIGN](https://datatracker.ietf.org/doc/html/draft-ietf-dnsop-rfc5933-bis-13)\] was also adopted but has seen very limited use, likely because it is a national algorithm specific to a very small number of countries.[¶](#section-3-1)

Implementation developers who want to know which algorithms to implement in DNSSEC software should refer to \[[RFC8624](https://datatracker.ietf.org/doc/html/rfc8624)\]. Note that this specification is only about what algorithms should and should not be included in implementations, i.e., it is not advice about which algorithms zone operators should or should not use for signing, nor which algorithms recursive resolver operators should or should not use for validation.[¶](#section-3-2)

## 4.

The DNSSEC community has extended the DNSSEC core and the cryptographic algorithms, both in terms of describing good operational practices and in new protocols. Some of the RFCs that describe these extensions include the following:[¶](#section-4-1)

- \[[RFC5011](https://datatracker.ietf.org/doc/html/rfc5011)\] describes a method to help resolvers update their DNSSEC trust anchors in an automated fashion. This method was used in 2018 to update the DNS root trust anchor.[¶](#section-4-2.1)
- \[[RFC6781](https://datatracker.ietf.org/doc/html/rfc6781)\] is a compendium of operational practices that may not be obvious from reading just the core specifications.[¶](#section-4-2.2)
- \[[RFC7344](https://datatracker.ietf.org/doc/html/rfc7344)\] describes using the CDS and CDNSKEY resource records to help automate the maintenance of DS records in the parents of signed zones.[¶](#section-4-2.3)
- \[[RFC8078](https://datatracker.ietf.org/doc/html/rfc8078)\] extends \[[RFC7344](https://datatracker.ietf.org/doc/html/rfc7344)\] by showing how to do initial setup of trusted relationships between signed parent and child zones.[¶](#section-4-2.4)
- \[[RFC8198](https://datatracker.ietf.org/doc/html/rfc8198)\] describes how a validating resolver can emit fewer queries in signed zones that use NSEC and NSEC3 for negative caching.[¶](#section-4-2.5)
- \[[RFC9077](https://datatracker.ietf.org/doc/html/rfc9077)\] updates \[[RFC8198](https://datatracker.ietf.org/doc/html/rfc8198)\] with respect to the TTL fields in signed records.[¶](#section-4-2.6)

## 5.

The documents listed above constitute the core of DNSSEC, the additional cryptographic algorithms, and the major extensions to DNSSEC. This section lists some additional documents that someone interested in implementing or operating DNSSEC might find of value:[¶](#section-5-1)

- \[[RFC4470](https://datatracker.ietf.org/doc/html/rfc4470)\] "describes how to construct DNSSEC NSEC resource records that cover a smaller range of names than called for by \[[RFC4034](https://datatracker.ietf.org/doc/html/rfc4034)\]. By generating and signing these records on demand, authoritative name servers can effectively stop the disclosure of zone contents otherwise made possible by walking the chain of NSEC records in a signed zone".[¶](#section-5-2.1)
- \[[RFC6975](https://datatracker.ietf.org/doc/html/rfc6975)\] "specifies a way for validating end-system resolvers to signal to a server which digital signature and hash algorithms they support".[¶](#section-5-2.2)
- \[[RFC7129](https://datatracker.ietf.org/doc/html/rfc7129)\] "provides additional background commentary and some context for the NSEC and NSEC3 mechanisms used by DNSSEC to provide authenticated denial-of-existence responses". This background is particularly important for understanding NSEC and NSEC3 usage.[¶](#section-5-2.3)
- \[[RFC7583](https://datatracker.ietf.org/doc/html/rfc7583)\] "describes the issues surrounding the timing of events in the rolling of a key in a DNSSEC-secured zone".[¶](#section-5-2.4)
- \[[RFC7646](https://datatracker.ietf.org/doc/html/rfc7646)\] "defines Negative Trust Anchors (NTAs), which can be used to mitigate DNSSEC validation failures by disabling DNSSEC validation at specified domains".[¶](#section-5-2.5)
- \[[RFC7958](https://datatracker.ietf.org/doc/html/rfc7958)\] "describes the format and publication mechanisms IANA has used to distribute the DNSSEC trust anchors".[¶](#section-5-2.6)
- \[[RFC8027](https://datatracker.ietf.org/doc/html/rfc8027)\] "describes problems that a Validating DNS resolver, stub-resolver, or application might run into within a non-compliant infrastructure".[¶](#section-5-2.7)
- \[[RFC8145](https://datatracker.ietf.org/doc/html/rfc8145)\] "specifies two different ways for validating resolvers to signal to a server which keys are referenced in their chain of trust".[¶](#section-5-2.8)
- \[[RFC8499](https://datatracker.ietf.org/doc/html/rfc8499)\] contains lists of terminology used when talking about DNS; Sections [10](https://www.rfc-editor.org/rfc/rfc8499#section-10) and [11](https://www.rfc-editor.org/rfc/rfc8499#section-11) cover DNSSEC.[¶](#section-5-2.9)
- \[[RFC8509](https://datatracker.ietf.org/doc/html/rfc8509)\] "specifies a mechanism that will allow an end user and third parties to determine the trusted key state for the root key of the resolvers that handle that user's DNS queries".[¶](#section-5-2.10)
- \[[RFC8901](https://datatracker.ietf.org/doc/html/rfc8901)\] "presents deployment models that accommodate this scenario \[when each DNS provider independently signs zone data with their own keys\] and describes these key-management requirements".[¶](#section-5-2.11)
- \[[RFC9276](https://datatracker.ietf.org/doc/html/rfc9276)\] "provides guidance on setting NSEC3 parameters based on recent operational deployment experience".[¶](#section-5-2.12)

There will certainly be other RFCs related to DNSSEC that are published after this one.[¶](#section-5-3)

## 6.

IANA already has three registry groups that relate to DNSSEC:[¶](#section-6-1)

- [DNSSEC algorithm numbers](https://www.iana.org/assignments/dns-sec-alg-numbers) [¶](#section-6-2.1)
- [DNSSEC NSEC3 parameters](https://www.iana.org/assignments/dnssec-nsec3-parameters) [¶](#section-6-2.2)
- [DNSSEC DS RRtype digest algorithms](https://www.iana.org/assignments/ds-rr-types) [¶](#section-6-2.3)

The rules for the DNSSEC algorithm registry were set in the core RFCs and updated by \[[RFC6014](https://datatracker.ietf.org/doc/html/rfc6014)\], \[[RFC6725](https://datatracker.ietf.org/doc/html/rfc6725)\], and \[[RFC9157](https://datatracker.ietf.org/doc/html/rfc9157)\].[¶](#section-6-3)

This document does not update or create any registry groups or registries.[¶](#section-6-4)

## 7.

All of the security considerations from all of the RFCs referenced in this document apply here.[¶](#section-7-1)

## 8.

### 8.1.

\[RFC3110\]

Eastlake 3rd, D., "RSA/SHA-1 SIGs and RSA KEYs in the Domain Name System (DNS)", RFC 3110, DOI 10.17487/RFC3110, May 2001, < [https://www.rfc-editor.org/info/rfc3110](https://www.rfc-editor.org/info/rfc3110) >.

\[RFC4033\]

Arends, R., Austein, R., Larson, M., Massey, D., and S. Rose, "DNS Security Introduction and Requirements", RFC 4033, DOI 10.17487/RFC4033, March 2005, < [https://www.rfc-editor.org/info/rfc4033](https://www.rfc-editor.org/info/rfc4033) >.

\[RFC4034\]

Arends, R., Austein, R., Larson, M., Massey, D., and S. Rose, "Resource Records for the DNS Security Extensions", RFC 4034, DOI 10.17487/RFC4034, March 2005, < [https://www.rfc-editor.org/info/rfc4034](https://www.rfc-editor.org/info/rfc4034) >.

\[RFC4035\]

Arends, R., Austein, R., Larson, M., Massey, D., and S. Rose, "Protocol Modifications for the DNS Security Extensions", RFC 4035, DOI 10.17487/RFC4035, March 2005, < [https://www.rfc-editor.org/info/rfc4035](https://www.rfc-editor.org/info/rfc4035) >.

\[RFC4509\]

Hardaker, W., "Use of SHA-256 in DNSSEC Delegation Signer (DS) Resource Records (RRs)", RFC 4509, DOI 10.17487/RFC4509, May 2006, < [https://www.rfc-editor.org/info/rfc4509](https://www.rfc-editor.org/info/rfc4509) >.

\[RFC5155\]

Laurie, B., Sisson, G., Arends, R., and D. Blacka, "DNS Security (DNSSEC) Hashed Authenticated Denial of Existence", RFC 5155, DOI 10.17487/RFC5155, March 2008, < [https://www.rfc-editor.org/info/rfc5155](https://www.rfc-editor.org/info/rfc5155) >.

\[RFC5702\]

Jansen, J., "Use of SHA-2 Algorithms with RSA in DNSKEY and RRSIG Resource Records for DNSSEC", RFC 5702, DOI 10.17487/RFC5702, October 2009, < [https://www.rfc-editor.org/info/rfc5702](https://www.rfc-editor.org/info/rfc5702) >.

\[RFC6840\]

Weiler, S., Ed. and D. Blacka, Ed., "Clarifications and Implementation Notes for DNS Security (DNSSEC)", RFC 6840, DOI 10.17487/RFC6840, February 2013, < [https://www.rfc-editor.org/info/rfc6840](https://www.rfc-editor.org/info/rfc6840) >.

### 8.2.

\[GOST-SIGN\]

Belyavsky, D., Dolmatov, V., Ed., and B. Makarenko, Ed., "Use of GOST 2012 Signature Algorithms in DNSKEY and RRSIG Resource Records for DNSSEC", Work in Progress, Internet-Draft, draft-ietf-dnsop-rfc5933-bis-13, 30 November 2022, < [https://datatracker.ietf.org/doc/html/draft-ietf-dnsop-rfc5933-bis-13](https://datatracker.ietf.org/doc/html/draft-ietf-dnsop-rfc5933-bis-13) >.

\[RFC2065\]

Eastlake 3rd, D. and C. Kaufman, "Domain Name System Security Extensions", RFC 2065, DOI 10.17487/RFC2065, January 1997, < [https://www.rfc-editor.org/info/rfc2065](https://www.rfc-editor.org/info/rfc2065) >.

\[RFC2535\]

Eastlake 3rd, D., "Domain Name System Security Extensions", RFC 2535, DOI 10.17487/RFC2535, March 1999, < [https://www.rfc-editor.org/info/rfc2535](https://www.rfc-editor.org/info/rfc2535) >.

\[RFC2536\]

Eastlake 3rd, D., "DSA KEYs and SIGs in the Domain Name System (DNS)", RFC 2536, DOI 10.17487/RFC2536, March 1999, < [https://www.rfc-editor.org/info/rfc2536](https://www.rfc-editor.org/info/rfc2536) >.

\[RFC4470\]

Weiler, S. and J. Ihren, "Minimally Covering NSEC Records and DNSSEC On-line Signing", RFC 4470, DOI 10.17487/RFC4470, April 2006, < [https://www.rfc-editor.org/info/rfc4470](https://www.rfc-editor.org/info/rfc4470) >.

\[RFC5011\]

StJohns, M., "Automated Updates of DNS Security (DNSSEC) Trust Anchors", STD 74, RFC 5011, DOI 10.17487/RFC5011, September 2007, < [https://www.rfc-editor.org/info/rfc5011](https://www.rfc-editor.org/info/rfc5011) >.

\[RFC6014\]

Hoffman, P., "Cryptographic Algorithm Identifier Allocation for DNSSEC", RFC 6014, DOI 10.17487/RFC6014, November 2010, < [https://www.rfc-editor.org/info/rfc6014](https://www.rfc-editor.org/info/rfc6014) >.

\[RFC6605\]

Hoffman, P. and W.C.A. Wijngaards, "Elliptic Curve Digital Signature Algorithm (DSA) for DNSSEC", RFC 6605, DOI 10.17487/RFC6605, April 2012, < [https://www.rfc-editor.org/info/rfc6605](https://www.rfc-editor.org/info/rfc6605) >.

\[RFC6698\]

Hoffman, P. and J. Schlyter, "The DNS-Based Authentication of Named Entities (DANE) Transport Layer Security (TLS) Protocol: TLSA", RFC 6698, DOI 10.17487/RFC6698, August 2012, < [https://www.rfc-editor.org/info/rfc6698](https://www.rfc-editor.org/info/rfc6698) >.

\[RFC6725\]

Rose, S., "DNS Security (DNSSEC) DNSKEY Algorithm IANA Registry Updates", RFC 6725, DOI 10.17487/RFC6725, August 2012, < [https://www.rfc-editor.org/info/rfc6725](https://www.rfc-editor.org/info/rfc6725) >.

\[RFC6781\]

Kolkman, O., Mekking, W., and R. Gieben, "DNSSEC Operational Practices, Version 2", RFC 6781, DOI 10.17487/RFC6781, December 2012, < [https://www.rfc-editor.org/info/rfc6781](https://www.rfc-editor.org/info/rfc6781) >.

\[RFC6975\]

Crocker, S. and S. Rose, "Signaling Cryptographic Algorithm Understanding in DNS Security Extensions (DNSSEC)", RFC 6975, DOI 10.17487/RFC6975, July 2013, < [https://www.rfc-editor.org/info/rfc6975](https://www.rfc-editor.org/info/rfc6975) >.

\[RFC7129\]

Gieben, R. and W. Mekking, "Authenticated Denial of Existence in the DNS", RFC 7129, DOI 10.17487/RFC7129, February 2014, < [https://www.rfc-editor.org/info/rfc7129](https://www.rfc-editor.org/info/rfc7129) >.

\[RFC7344\]

Kumari, W., Gudmundsson, O., and G. Barwood, "Automating DNSSEC Delegation Trust Maintenance", RFC 7344, DOI 10.17487/RFC7344, September 2014, < [https://www.rfc-editor.org/info/rfc7344](https://www.rfc-editor.org/info/rfc7344) >.

\[RFC7583\]

Morris, S., Ihren, J., Dickinson, J., and W. Mekking, "DNSSEC Key Rollover Timing Considerations", RFC 7583, DOI 10.17487/RFC7583, October 2015, < [https://www.rfc-editor.org/info/rfc7583](https://www.rfc-editor.org/info/rfc7583) >.

\[RFC7646\]

Ebersman, P., Kumari, W., Griffiths, C., Livingood, J., and R. Weber, "Definition and Use of DNSSEC Negative Trust Anchors", RFC 7646, DOI 10.17487/RFC7646, September 2015, < [https://www.rfc-editor.org/info/rfc7646](https://www.rfc-editor.org/info/rfc7646) >.

\[RFC7958\]

Abley, J., Schlyter, J., Bailey, G., and P. Hoffman, "DNSSEC Trust Anchor Publication for the Root Zone", RFC 7958, DOI 10.17487/RFC7958, August 2016, < [https://www.rfc-editor.org/info/rfc7958](https://www.rfc-editor.org/info/rfc7958) >.

\[RFC8027\]

Hardaker, W., Gudmundsson, O., and S. Krishnaswamy, "DNSSEC Roadblock Avoidance", BCP 207, RFC 8027, DOI 10.17487/RFC8027, November 2016, < [https://www.rfc-editor.org/info/rfc8027](https://www.rfc-editor.org/info/rfc8027) >.

\[RFC8078\]

Gudmundsson, O. and P. Wouters, "Managing DS Records from the Parent via CDS/CDNSKEY", RFC 8078, DOI 10.17487/RFC8078, March 2017, < [https://www.rfc-editor.org/info/rfc8078](https://www.rfc-editor.org/info/rfc8078) >.

\[RFC8080\]

Sury, O. and R. Edmonds, "Edwards-Curve Digital Security Algorithm (EdDSA) for DNSSEC", RFC 8080, DOI 10.17487/RFC8080, February 2017, < [https://www.rfc-editor.org/info/rfc8080](https://www.rfc-editor.org/info/rfc8080) >.

\[RFC8145\]

Wessels, D., Kumari, W., and P. Hoffman, "Signaling Trust Anchor Knowledge in DNS Security Extensions (DNSSEC)", RFC 8145, DOI 10.17487/RFC8145, April 2017, < [https://www.rfc-editor.org/info/rfc8145](https://www.rfc-editor.org/info/rfc8145) >.

\[RFC8198\]

Fujiwara, K., Kato, A., and W. Kumari, "Aggressive Use of DNSSEC-Validated Cache", RFC 8198, DOI 10.17487/RFC8198, July 2017, < [https://www.rfc-editor.org/info/rfc8198](https://www.rfc-editor.org/info/rfc8198) >.

\[RFC8499\]

Hoffman, P., Sullivan, A., and K. Fujiwara, "DNS Terminology", BCP 219, RFC 8499, DOI 10.17487/RFC8499, January 2019, < [https://www.rfc-editor.org/info/rfc8499](https://www.rfc-editor.org/info/rfc8499) >.

\[RFC8509\]

Huston, G., Damas, J., and W. Kumari, "A Root Key Trust Anchor Sentinel for DNSSEC", RFC 8509, DOI 10.17487/RFC8509, December 2018, < [https://www.rfc-editor.org/info/rfc8509](https://www.rfc-editor.org/info/rfc8509) >.

\[RFC8624\]

Wouters, P. and O. Sury, "Algorithm Implementation Requirements and Usage Guidance for DNSSEC", RFC 8624, DOI 10.17487/RFC8624, June 2019, < [https://www.rfc-editor.org/info/rfc8624](https://www.rfc-editor.org/info/rfc8624) >.

\[RFC8901\]

Huque, S., Aras, P., Dickinson, J., Vcelak, J., and D. Blacka, "Multi-Signer DNSSEC Models", RFC 8901, DOI 10.17487/RFC8901, September 2020, < [https://www.rfc-editor.org/info/rfc8901](https://www.rfc-editor.org/info/rfc8901) >.

\[RFC9077\]

van Dijk, P., "NSEC and NSEC3: TTLs and Aggressive Use", RFC 9077, DOI 10.17487/RFC9077, July 2021, < [https://www.rfc-editor.org/info/rfc9077](https://www.rfc-editor.org/info/rfc9077) >.

\[RFC9157\]

Hoffman, P., "Revised IANA Considerations for DNSSEC", RFC 9157, DOI 10.17487/RFC9157, December 2021, < [https://www.rfc-editor.org/info/rfc9157](https://www.rfc-editor.org/info/rfc9157) >.

\[RFC9276\]

Hardaker, W. and V. Dukhovni, "Guidance for NSEC3 Parameter Settings", BCP 236, RFC 9276, DOI 10.17487/RFC9276, August 2022, < [https://www.rfc-editor.org/info/rfc9276](https://www.rfc-editor.org/info/rfc9276) >.

## Acknowledgements

The DNS world owes a depth of gratitude to the authors and other contributors to the core DNSSEC documents and to the notable DNSSEC extensions.[¶](#appendix-A-1)

In addition, the following people made significant contributions to early draft versions of this document: and.[¶](#appendix-A-2)

## Author's Address

Paul Hoffman

ICANN

Email: [paul.hoffman@icann.org](mailto:paul.hoffman@icann.org)