package data

import (
  "fmt"
  "strconv"
)

type ActionCompareOp string

const (
	Equal ActionCompareOp = "="
	NotEqual ActionCompareOp = "!="
	Less ActionCompareOp = "<"
	LessEqual ActionCompareOp = "<="
  Greater ActionCompareOp = ">"
	GreaterEqual ActionCompareOp = ">="
)

type ActionMod interface {
  GetMod(c Character) ActionModAddSub
}

type Action struct {
  ID string `json:"id"`
  Name string `json:"name"`
  Description string `json:"desc"`
  Tags []string `json:"tags"`
  Mods []ActionMod `json:"mod"`
}

type actionModBase struct {
  AttributeID string `json:"id"`
}

type ActionModAddSub struct {
  actionModBase
  Modifier int16 `json:"mod"`
}

func (a ActionModAddSub) GetMod(c Character) ActionModAddSub {
  return a
}

type ActionModMultDiv struct {
  actionModBase
  Multiplier int16 `json:"mult"`
}

func (a ActionModMultDiv) GetMod(c Character) ActionModAddSub {
  attr := c.GetAttribute(a.AttributeID)
  if attr.ShortName == "N/A" {
    panic(fmt.Sprint("Attribute", a.AttributeID, "not found"))
  }
  return ActionModAddSub{actionModBase: actionModBase{AttributeID: a.AttributeID}, Modifier: attr.Value * a.Multiplier}
}

type ActionModCompare struct {
  actionModBase
  Condition struct {
    AttributeID string `json:"id"`
    Op ActionCompareOp `json:"op"`
    Value string `json:"value"`
  } `json:"condition"`
  Action ActionMod `json:"action"`
}

func (a ActionModCompare) GetMod(c Character) ActionModAddSub {
  attrCond := c.GetAttributeStrict(a.Condition.AttributeID)
  compVal, err := strconv.ParseInt(a.Condition.Value, 10, 16)
  if err != nil {
    
  }
  compVal = int64(c.GetAttributeStrict(a.Condition.Value).Value)
  switch(a.Condition.Op) {
    case Equal: {
      if attrCond.Value == int16(compVal) {
        return a.Action.GetMod(c)
      }
    }
    case NotEqual: {
      if attrCond.Value != int16(compVal) {
        return a.Action.GetMod(c)
      }
    }
    case Greater: {
      if attrCond.Value > int16(compVal) {
        return a.Action.GetMod(c)
      }
    }
    case GreaterEqual: {
      if attrCond.Value >= int16(compVal) {
        return a.Action.GetMod(c)
      }
    }
    case Less: {
      if attrCond.Value < int16(compVal) {
        return a.Action.GetMod(c)
      }
    }
    case LessEqual: {
      if attrCond.Value <= int16(compVal) {
        return a.Action.GetMod(c)
      }
    }
  }
  return ActionModAddSub{actionModBase: actionModBase{AttributeID: a.AttributeID}, Modifier: 0}
}