/*
Package graphdb provides an embeddable hypergraph database.

Basics

Here be overview.

Features

Here be features.
*/

package graphdb
/*
Grammar:

Tokens starting with a lower case letter are terminals; int(n)
and uint(n) represent the signed/unsigned encodings of the value n.

GobStream:
	DelimitedMessage*
DelimitedMessage:
	uint(lengthOfMessage) Message
Message:
	TypeSequence TypedValue
FieldValue:
	builtinValue | ArrayValue | MapValue | SliceValue | StructValue | InterfaceValue
InterfaceValue:
	NilInterfaceValue | NonNilInterfaceValue
NilInterfaceValue:
	uint(0)
NonNilInterfaceValue:
	ConcreteTypeName TypeSequence InterfaceContents
ConcreteTypeName:
	uint(lengthOfName) [already read=n] name
InterfaceContents:
	int(concreteTypeId) DelimitedValue
DelimitedValue:
	uint(length) Value
ArrayValue:
	uint(n) FieldValue*n [n elements]
MapValue:
	uint(n) (FieldValue FieldValue)*n  [n (key, value) pairs]
SliceValue:
	uint(n) FieldValue*n [n elements]
StructValue:
	(uint(fieldDelta) FieldValue)*
*/

/*
For implementers and the curious, here is an encoded example.  Given
	type Point struct {X, Y int}
and the value
	p := Point{22, 33}
the bytes transmitted that encode p will be:
	1f ff 81 03 01 01 05 50 6f 69 6e 74 01 ff 82 00
	01 02 01 01 58 01 04 00 01 01 59 01 04 00 00 00
	07 ff 82 01 2c 01 42 00
*/
