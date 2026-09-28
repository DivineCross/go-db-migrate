package script

func init() {
	Register(s005Up, s005Down)
}

func s005Up() {
	db.InsertWhenNotMatched("EventDef", s005EventDefJson, "eventId")
}

func s005Down() {
}

const s005EventDefJson = `[
  {
    "eventId": "E_001",
    "name": "服務啟動完成",
    "level": {
      "$numberInt": "1"
    },
    "shouldRecover": false,
    "message": "服務已完成初始化，開始接受連線。"
  },
  {
    "eventId": "E_002",
    "name": "服務正常停止",
    "level": {
      "$numberInt": "1"
    },
    "shouldRecover": false,
    "message": "服務已完成待處理工作並正常關閉。"
  },
  {
    "eventId": "E_003",
    "name": "磁碟空間不足",
    "level": {
      "$numberInt": "2"
    },
    "shouldRecover": true,
    "message": "磁碟可用空間低於警戒值，請清理檔案或擴充容量。",
    "recoverMessage": "磁碟可用空間已恢復至正常範圍。"
  },
  {
    "eventId": "E_004",
    "name": "設定重新載入失敗",
    "level": {
      "$numberInt": "2"
    },
    "shouldRecover": false,
    "message": "新設定未通過驗證，服務將繼續使用原有設定。"
  },
  {
    "eventId": "E_005",
    "name": "資料庫連線中斷",
    "level": {
      "$numberInt": "3"
    },
    "shouldRecover": true,
    "message": "無法連線至資料庫，資料存取功能暫時無法使用。",
    "recoverMessage": "資料庫連線已恢復，資料存取功能恢復正常。"
  },
  {
    "eventId": "E_006",
    "name": "服務異常終止",
    "level": {
      "$numberInt": "3"
    },
    "shouldRecover": false,
    "message": "服務發生無法處理的錯誤而終止，請檢查錯誤紀錄。"
  }
]`
