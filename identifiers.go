package irmago

import "strings"

type metaObjectIdentifier string

// SchemeManagerIdentifier identifies a scheme manager. Equal to its ID. For example "irma-demo".
type SchemeManagerIdentifier struct {
	metaObjectIdentifier
}

// IssuerIdentifier identifies an inssuer. For example "irma-demo.RU".
type IssuerIdentifier struct {
	metaObjectIdentifier
}

// CredentialTypeIdentifier identifies a credentialtype. For example "irma-demo.RU.studentCard".
type CredentialTypeIdentifier struct {
	metaObjectIdentifier
}

// AttributeTypeIdentifier identifies an attribute. For example "irma-demo.RU.studentCard.studentID".
type AttributeTypeIdentifier struct {
	metaObjectIdentifier
}

type CredentialIdentifier struct {
	Type  CredentialTypeIdentifier
	Index int
	Count int
}

type AttributeIdentifier struct {
	Type  AttributeTypeIdentifier
	Index int
	Count int
}

func (oi metaObjectIdentifier) Parent() string {
	str := string(oi)
	return str[:strings.LastIndex(str, ".")]
}

func (oi metaObjectIdentifier) Name() string {
	str := string(oi)
	return str[strings.LastIndex(str, ".")+1:]
}

func (oi metaObjectIdentifier) String() string {
	return string(oi)
}

// NewSchemeManagerIdentifier converts the specified identifier to a SchemeManagerIdentifier.
func NewSchemeManagerIdentifier(id string) SchemeManagerIdentifier {
	return SchemeManagerIdentifier{metaObjectIdentifier(id)}
}

// NewIssuerIdentifier converts the specified identifier to a IssuerIdentifier.
func NewIssuerIdentifier(id string) IssuerIdentifier {
	return IssuerIdentifier{metaObjectIdentifier(id)}
}

// NewCredentialTypeIdentifier converts the specified identifier to a CredentialTypeIdentifier.
func NewCredentialTypeIdentifier(id string) CredentialTypeIdentifier {
	return CredentialTypeIdentifier{metaObjectIdentifier(id)}
}

// NewAttributeTypeIdentifier converts the specified identifier to a AttributeTypeIdentifier.
func NewAttributeTypeIdentifier(id string) AttributeTypeIdentifier {
	return AttributeTypeIdentifier{metaObjectIdentifier(id)}
}

// SchemeManagerIdentifier returns the scheme manager identifer of the issuer.
func (id IssuerIdentifier) SchemeManagerIdentifier() SchemeManagerIdentifier {
	return NewSchemeManagerIdentifier(id.Parent())
}

// IssuerIdentifier returns the IssuerIdentifier of the credential identifier.
func (id CredentialTypeIdentifier) IssuerIdentifier() IssuerIdentifier {
	return NewIssuerIdentifier(id.Parent())
}

// CredentialTypeIdentifier returns the CredentialTypeIdentifier of the attribute identifier.
func (id AttributeTypeIdentifier) CredentialTypeIdentifier() CredentialTypeIdentifier {
	return NewCredentialTypeIdentifier(id.Parent())
}

func (id AttributeTypeIdentifier) IsCredential() bool {
	return strings.Count(id.String(), ".") == 2
}