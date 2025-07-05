package controllers

import (
	"sort"
	"strconv"
	"strings"

	"go_binance_futures/models"
	"go_binance_futures/spot"
	"go_binance_futures/utils"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)


type SpotController struct {
	web.Controller
}

func (ctrl *SpotController) Get() {
	paramsSort := ctrl.GetString("sort")
	symbol_type := ctrl.GetString("symbol_type")
	symbol := ctrl.GetString("symbol")
	enable := ctrl.GetString("enable")
	pin := ctrl.GetString("pin")
	
	o := orm.NewOrm()
	var symbols []models.SpotSymbols
	query := o.QueryTable("spot_symbols")
	if symbol_type != "" {
		query = query.Filter("type", symbol_type)
	}
	if symbol != "" {
		query = query.Filter("symbol__contains", symbol)
	}
	if enable != "" {
		query = query.Filter("enable", enable)
	}
	if pin != "" {
		query = query.Filter("pin", 1)
	}
	_, err := query.OrderBy("-Pin", "ID").All(&symbols)
	if err != nil {
		ctrl.Ctx.Resp(utils.ResJson(400, nil, "error"))
	}
	
	if strings.HasPrefix(paramsSort, "percent_change") {
		sort.SliceStable(symbols, func(i, j int) bool {
			base := symbols[i].Pin >= symbols[j].Pin
			if paramsSort == "percent_change+" {
				return base && symbols[i].PercentChange >= symbols[j].PercentChange // 涨幅从大到小排序
			} else if paramsSort == "percent_change-" {
				return base && symbols[i].PercentChange < symbols[j].PercentChange // 涨幅从小到大排序
			} else {
				return true
			}
		})
	}

	ctrl.Ctx.Resp(map[string]interface{} {
		"code": 200,
		"data": symbols,
		"msg": "success",
	})
}

func (ctrl *SpotController) Edit() {
	id := ctrl.Ctx.Input.Param(":id")
	var symbols models.SpotSymbols
	o := orm.NewOrm()
	o.QueryTable(symbols.TableName()).Filter("Id", id).One(&symbols)
	
	ctrl.BindJSON(&symbols)
	
	_, err := o.Update(&symbols) // _ 是受影响的条数
    if err != nil {
        // 处理错误
		ctrl.Ctx.Resp(utils.ResJson(400, nil, err.Error()))
		return
    }
	ctrl.Ctx.Resp(map[string]interface{} {
		"code": 200,
		"data": symbols,
		"msg": "success",
	})
}

func (ctrl *SpotController) Delete() {
	id := ctrl.Ctx.Input.Param(":id")
	symbols := new(models.SpotSymbols)
	intId, _ := strconv.ParseInt(id, 10, 64)
	symbols.ID = intId
	o := orm.NewOrm()
	
	_, err := o.Delete(symbols)
    if err != nil {
        // 处理错误
		ctrl.Ctx.Resp(utils.ResJson(400, nil, err.Error()))
		return
    }
	ctrl.Ctx.Resp(map[string]interface{} {
		"code": 200,
		"msg": "success",
	})
}

func (ctrl *SpotController) Post() {
	symbols := new(models.SpotSymbols)
	ctrl.BindJSON(&symbols)
	symbols.PercentChange = 0
	symbols.Close = "0"
	symbols.Open = "0"
	symbols.Low = "0"
	symbols.High = "0"
	
	symbols.StepSize = "0.1"
	symbols.TickSize = "0.1"
	symbols.Usdt = "10"
	symbols.Profit = "100"
	symbols.Loss = "100"
	
	o := orm.NewOrm()
	id, err := o.Insert(symbols)
	
	go func() {
		spot.UpdateSymbolsTradePrecision() // 更新合约交易精度
	}()
	
    if err != nil {
        // 处理错误
		ctrl.Ctx.Resp(utils.ResJson(400, nil, err.Error()))
		return
    }
	symbols.ID = id
	
	ctrl.Ctx.Resp(map[string]interface{} {
		"code": 200,
		"data": symbols,
		"msg": "success",
	})
}

func (ctrl *SpotController) UpdateEnable() {
	flag := ctrl.Ctx.Input.Param(":flag")
	
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE spot_symbols SET enable = ?", flag).Exec()
	if err != nil {
		// 处理错误
		ctrl.Ctx.Resp(utils.ResJson(400, nil, err.Error()))
		return
	}
	ctrl.Ctx.Resp(utils.ResJson(200, nil))
}

func (ctrl *SpotController) BatchEdit() {
	params := new(BatchEditParams)
	ctrl.BindJSON(&params)
	query := "UPDATE spot_symbols SET"
	
	if params.Usdt != "" {
		query += " usdt = " + params.Usdt + ","
	}
	if (params.Profit != "") {
		query += " profit = " + params.Profit + ","	
	}
	if (params.Loss != "") {
		query += " loss = " + params.Loss + ","
	}
	if (params.StrategyType != "") {
		query += " strategy_type = '" + params.StrategyType + "',"
	}
	if (params.StrategyTemplateId != 0) {
		var template models.StrategyTemplates
		orm.NewOrm().QueryTable("strategy_templates").Filter("Id", params.StrategyTemplateId).One(&template)
		query += " technology = '" + template.Technology + "',"
		query += " strategy = '" + template.Strategy + "',"
	}
	
	if strings.HasSuffix(query, ",") {
		query = strings.TrimSuffix(query, ",")
		// query += " WHERE enable = 1" // 只更新开启的交易对
		_, err := orm.NewOrm().Raw(query).Exec()
		if err != nil {
			// 处理错误
			ctrl.Ctx.Resp(utils.ResJson(400, nil, err.Error()))
			return
		}
	}
	
	// go func() {
	// 	spot.UpdateSymbolsTradePrecision() // 更新合约交易精度
	// }()
	
	ctrl.Ctx.Resp(map[string]interface{} {
		"code": 200,
		"msg": "success",
	})
}
