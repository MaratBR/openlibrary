package i18n

import "strings"

type keyDef struct {
	keys map[string]*keyDef
}

func (k *keyDef) addInner(key string) *keyDef {
	if k.keys == nil {
		k.keys = make(map[string]*keyDef)
	}

	def := &keyDef{}
	k.keys[key] = def
	return def
}

func (def *keyDef) mergeWith(other *keyDef) {
	if def.keys == nil {
		def.keys = other.keys
		return
	}

	for k, v := range other.keys {
		if existing, ok := def.keys[k]; ok {
			existing.mergeWith(v)
		} else {
			def.keys[k] = v
		}
	}
}

func (def *keyDef) path(key string) *keyDef {
	if key == "" {
		return def
	}
	parts := strings.Split(key, ".")

	target := def
	for _, k := range parts {
		if target.keys == nil {
			return nil
		}
		if newTarget, ok := target.keys[k]; ok {
			target = newTarget
		} else {
			return nil
		}
	}

	return target
}

func (def *keyDef) walkFullKeys(cb func(fullKey string)) {
	def.internalWalk(cb, "")
}

func (def *keyDef) internalWalk(cb func(fullKey string), prefix string) {
	if def.keys == nil {
		return
	}

	for k, v := range def.keys {
		if v.keys == nil {
			cb(prefix + k)
		} else {
			v.internalWalk(cb, prefix+k+".")
		}
	}
}
