package validate

import (
	"github.com/graniticio/granitic/test"
	"github.com/graniticio/granitic/types"
	"testing"
)

func TestUnsetObjDetection(t *testing.T) {

	ob := NewObjectValidatorBuilder("DEF", nil)

	ov, err := ob.parseRule("CP", []string{"REQ:MISSING"})

	test.ExpectNil(t, err)

	sub := new(Parent)
	vc := new(validationContext)
	vc.Subject = sub

	r, err := ov.Validate(vc)
	test.ExpectNil(t, err)
	c := r.ErrorCodes

	test.ExpectInt(t, len(c), 1)

	ov, err = ob.parseRule("CM", []string{"REQ:MISSING"})

	test.ExpectNil(t, err)

	r, err = ov.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)

	ov, err = ob.parseRule("CV", []string{"REQ:MISSING"})

	test.ExpectNil(t, err)

	r, err = ov.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

}

func TestSetObjDetection(t *testing.T) {

	ob := NewObjectValidatorBuilder("DEF", nil)

	ov, err := ob.parseRule("CP", []string{"REQ:MISSING"})

	test.ExpectNil(t, err)

	sub := new(Parent)
	sub.CP = new(Child)
	vc := new(validationContext)
	vc.Subject = sub

	r, err := ov.Validate(vc)
	test.ExpectNil(t, err)
	c := r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.CM = make(map[string]interface{})
	ov, err = ob.parseRule("CM", []string{"REQ:MISSING"})

	test.ExpectNil(t, err)

	r, err = ov.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	ov, err = ob.parseRule("CV", []string{"REQ:MISSING"})

	test.ExpectNil(t, err)

	r, err = ov.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

}

func TestObjectMExFieldDetection(t *testing.T) {
	vb := NewObjectValidatorBuilder("DEF", nil)

	bv, err := vb.parseRule("CP", []string{"MEX:setField1,setField2:BAD_MEX"})

	test.ExpectNil(t, err)

	sub := new(Parent)
	vc := new(validationContext)
	vc.Subject = sub
	vc.KnownSetFields = types.NewOrderedStringSet([]string{})

	sub.CP = new(Child)

	r, err := bv.Validate(vc)
	test.ExpectNil(t, err)
	c := r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	vc.KnownSetFields.Add("ignoreField")

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	vc.KnownSetFields.Add("setField1")

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "BAD_MEX")

	vc.KnownSetFields = types.NewOrderedStringSet([]string{})
	vc.KnownSetFields.Add("setField2")

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "BAD_MEX")

}

func TestInvalidTypeHandling(t *testing.T) {
	ob := NewObjectValidatorBuilder("DEF", nil)

	ov, err := ob.parseRule("S", []string{"REQ:MISSING"})

	test.ExpectNil(t, err)

	sub := new(Parent)
	sub.CP = new(Child)
	vc := new(validationContext)
	vc.Subject = sub

	_, err = ov.Validate(vc)
	test.ExpectNotNil(t, err)

}

type Parent struct {
	CP *Child
	CV Child
	CM map[string]interface{}
	S  string
}

type Child struct {
}
