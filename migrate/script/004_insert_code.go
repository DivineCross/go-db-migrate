package script

import (
	db "migrate/driver"
)

func init() {
	Register(s004Up, s004Down)
}

func s004Up() {
	db.InsertWhenNotMatched("Code", s004CodeJson, "category", "code")
}

func s004Down() {
}

const s004CodeJson = `[
  {
    "category": "color",
    "code": "1",
    "text": "紅色",
    "sort": {
      "$numberInt": "1"
    },
    "disabled": false
  },
  {
    "category": "color",
    "code": "2",
    "text": "綠色",
    "sort": {
      "$numberInt": "2"
    },
    "disabled": false
  },
  {
    "category": "color",
    "code": "3",
    "text": "藍色",
    "sort": {
      "$numberInt": "3"
    },
    "disabled": false
  },
  {
    "category": "size",
    "code": "1",
    "text": "小",
    "sort": {
      "$numberInt": "1"
    },
    "disabled": false
  },
  {
    "category": "size",
    "code": "2",
    "text": "中",
    "sort": {
      "$numberInt": "2"
    },
    "disabled": false
  },
  {
    "category": "size",
    "code": "3",
    "text": "大",
    "sort": {
      "$numberInt": "3"
    },
    "disabled": false
  },
  {
    "category": "status",
    "code": "1",
    "text": "待處理",
    "sort": {
      "$numberInt": "1"
    },
    "disabled": false
  },
  {
    "category": "status",
    "code": "2",
    "text": "處理中",
    "sort": {
      "$numberInt": "2"
    },
    "disabled": false
  },
  {
    "category": "status",
    "code": "3",
    "text": "已完成",
    "sort": {
      "$numberInt": "3"
    },
    "disabled": false
  }
]`
