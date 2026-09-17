package nodeGRPC

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// distinct mutates m in place, setting exactly one field to a value derived from seed, and
// returns m for convenient use inline at fixture-construction call sites. It exists so this
// package's "*_Success"/"*_DrainsStream" tests exercise real request/response marshaling instead
// of trivially satisfying proto.Equal comparisons regardless of wire-level bugs — a zero-valued
// fixture (&X{}) round-trips through GRPC identically whether or not a wrapper drops or corrupts
// a field, since every field is already at its zero value on both sides (production-readiness
// finding 7).
//
// Rather than hand-writing a specific non-zero field per message type (there are ~80 distinct
// request/response types across this package's tests), distinct walks the message's own
// descriptor via protoreflect and sets the first field it finds that it knows how to mutate,
// preferring (in order): a top-level scalar field, then a singular nested-message field
// (recursing into it), then a repeated scalar field (appending one value), then a repeated
// message field (appending one recursively-distinguished element). This keeps every fixture
// change mechanical and self-maintaining as the underlying protos evolve, rather than requiring
// ~80 hand-picked field names to be kept in sync by hand.
//
// seed should differ between the request and response fixtures of the same test (and between
// distinct items of the same message type within a single streaming test), so the resulting
// values are actually distinguishable from one another, catching both field-dropping bugs and
// item-ordering bugs.
//
// A small number of message types used in these tests are genuinely all-zero-field (e.g.
// GetVersionRequest, GetIdentityRequest, TransactionEventRequest) — for those, distinct is a
// no-op and the fixture legitimately stays {}; there's nothing to distinguish.
func distinct[T proto.Message](m T, seed int) T {
	setDistinctField(m.ProtoReflect(), seed)
	return m
}

// setDistinctField implements the walk described on distinct's doc comment. It returns true if
// it found and set a field, so callers recursing into nested messages know whether that nested
// message ended up distinguishable (and therefore worth keeping/setting on the parent).
func setDistinctField(m protoreflect.Message, seed int) bool {
	fields := m.Descriptor().Fields()

	// Pass 1: singular scalar fields on this message.
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if fd.IsList() || fd.IsMap() || fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
			continue
		}
		if val, ok := scalarValue(fd, seed); ok {
			m.Set(fd, val)
			return true
		}
	}

	// Pass 2: singular nested-message fields - recurse into the first one that yields a
	// distinguishable value.
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if fd.IsList() || fd.IsMap() {
			continue
		}
		if fd.Kind() != protoreflect.MessageKind && fd.Kind() != protoreflect.GroupKind {
			continue
		}
		sub := m.NewField(fd)
		if setDistinctField(sub.Message(), seed) {
			m.Set(fd, sub)
			return true
		}
	}

	// Pass 3: repeated scalar fields - append one distinguishing element.
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if !fd.IsList() || fd.IsMap() {
			continue
		}
		if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
			continue
		}
		if val, ok := scalarValue(fd, seed); ok {
			list := m.Mutable(fd).List()
			list.Append(val)
			return true
		}
	}

	// Pass 4: repeated message fields - append one recursively-distinguished element.
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if !fd.IsList() || fd.IsMap() {
			continue
		}
		if fd.Kind() != protoreflect.MessageKind && fd.Kind() != protoreflect.GroupKind {
			continue
		}
		list := m.Mutable(fd).List()
		elem := list.NewElement()
		if setDistinctField(elem.Message(), seed) {
			list.Append(elem)
			return true
		}
	}

	return false
}

// scalarValue returns a value derived from seed appropriate for fd's kind, and false if fd's
// kind isn't a scalar this helper knows how to mutate (message/group are handled by the caller;
// an enum with only a zero value has nothing meaningful to set it to).
func scalarValue(fd protoreflect.FieldDescriptor, seed int) (protoreflect.Value, bool) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(true), true
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(int32(seed)), true
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(int64(seed)), true
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(uint32(seed)), true
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(uint64(seed)), true
	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(float32(seed)), true
	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(float64(seed)), true
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(fmt.Sprintf("distinct-%d", seed)), true
	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes([]byte(fmt.Sprintf("distinct-%d", seed))), true
	case protoreflect.EnumKind:
		values := fd.Enum().Values()
		if values.Len() > 1 {
			idx := seed % (values.Len() - 1)
			if idx < 0 {
				idx += values.Len() - 1
			}
			return protoreflect.ValueOfEnum(values.Get(idx + 1).Number()), true
		}
		return protoreflect.Value{}, false
	default:
		return protoreflect.Value{}, false
	}
}
