1. nacos配置改动
   1. 新建dataId=HedgeParamsV2，代替HedgeParams。
      1. 将数据里的float64类型全部改为string类型。
      2. instrument => symbol //品种
      3. max_margin => max_leverage //最大杠杆倍数
      4. max_order => single_max_order_volume //单次最大下单量
      5. max_match_level => max_match_gear //最大成交档位
   2. 新建dataId=CurrencyParamsV2，代替CurrencyParams。
      1. kind类型由string改为int（币种类型，1:法币 2: 公链数字货币 3：稳定币）
      2. user => operator
   3. 新建dataId=SymbolParamsV2，代替SymbolParams。
      1. 将数据里的float64类型全部改为string类型。
      2. symbol_kind类型由string改为int（品种类型，如 1-现货、2-期货等）
      3. bid => underlying
      4. amount_precision => volume_precision //数量精度
      5. sum_precision => amount_precision //金额精度
      6. min_order => single_min_order_volume
      7. max_order => single_max_order_volume
      8. min_money => single_min_order_amount
      9. max_money => single_max_order_amount 
      10. max_match_level => max_match_gear 
      11. otc_min_order => otc_min_order_volume
      12. otc_max_order => otc_max_order_volume
      13. otc_min_price => otc_min_order_amount
      14. otc_max_price => otc_max_order_amount
      15. user => operator
   4. 新建dataId=MemberGroupTradeParamsV2，代替MemberGroupTradeParams。
      1. spread类型由float64改为string //点差，品种tick值的整数倍
      2. taker_fee类型由float64改为string
      3. maker_fee类型由float64改为string
      4. user => operator