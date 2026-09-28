package script

func init() {
	Register(s006Up, s006Down)
}

func s006Up() {
	db.InsertWhenNotMatched("Pricing", s006PricingJson, "pricingId", "productId")
}

func s006Down() {
}

const s006PricingJson = `[
  {
    "pricingId": "P_001",
    "startAt": {
      "$date": {
        "$numberLong": "1790812800000"
      }
    },
    "productId": "MX100",
    "litePrice": {
      "$numberDouble": "299.0"
    },
    "plusPrice": {
      "$numberDouble": "599.0"
    },
    "proPrice": {
      "$numberDouble": "999.0"
    }
  },
  {
    "pricingId": "P_001",
    "startAt": {
      "$date": {
        "$numberLong": "1790812800000"
      }
    },
    "productId": "MX200",
    "litePrice": {
      "$numberDouble": "499.0"
    },
    "plusPrice": {
      "$numberDouble": "899.0"
    },
    "proPrice": {
      "$numberDouble": "1499.0"
    }
  },
  {
    "pricingId": "P_001",
    "startAt": {
      "$date": {
        "$numberLong": "1790812800000"
      }
    },
    "productId": "MX300",
    "litePrice": {
      "$numberDouble": "799.0"
    },
    "plusPrice": {
      "$numberDouble": "1299.0"
    },
    "proPrice": {
      "$numberDouble": "1999.0"
    }
  },
  {
    "pricingId": "P_002",
    "startAt": {
      "$date": {
        "$numberLong": "1796083200000"
      }
    },
    "productId": "MX100",
    "litePrice": {
      "$numberDouble": "299.0"
    },
    "plusPrice": {
      "$numberDouble": "599.0"
    },
    "proPrice": {
      "$numberDouble": "999.0"
    }
  },
  {
    "pricingId": "P_002",
    "startAt": {
      "$date": {
        "$numberLong": "1796083200000"
      }
    },
    "productId": "MX200",
    "litePrice": {
      "$numberDouble": "499.0"
    },
    "plusPrice": {
      "$numberDouble": "899.0"
    },
    "proPrice": {
      "$numberDouble": "1499.0"
    }
  },
  {
    "pricingId": "P_002",
    "startAt": {
      "$date": {
        "$numberLong": "1796083200000"
      }
    },
    "productId": "MX300",
    "litePrice": {
      "$numberDouble": "799.0"
    },
    "plusPrice": {
      "$numberDouble": "1299.0"
    },
    "proPrice": {
      "$numberDouble": "1999.0"
    }
  },
  {
    "pricingId": "P_003",
    "startAt": {
      "$date": {
        "$numberLong": "1801440000000"
      }
    },
    "productId": "MX100",
    "litePrice": {
      "$numberDouble": "299.0"
    },
    "plusPrice": {
      "$numberDouble": "599.0"
    },
    "proPrice": {
      "$numberDouble": "999.0"
    }
  },
  {
    "pricingId": "P_003",
    "startAt": {
      "$date": {
        "$numberLong": "1801440000000"
      }
    },
    "productId": "MX200",
    "litePrice": {
      "$numberDouble": "499.0"
    },
    "plusPrice": {
      "$numberDouble": "899.0"
    },
    "proPrice": {
      "$numberDouble": "1499.0"
    }
  },
  {
    "pricingId": "P_003",
    "startAt": {
      "$date": {
        "$numberLong": "1801440000000"
      }
    },
    "productId": "MX300",
    "litePrice": {
      "$numberDouble": "799.0"
    },
    "plusPrice": {
      "$numberDouble": "1299.0"
    },
    "proPrice": {
      "$numberDouble": "1999.0"
    }
  }
]`
