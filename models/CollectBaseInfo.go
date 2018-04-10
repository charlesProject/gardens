package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strings"
)

type DtuRowOfDay struct {
	DTU_no       string    `orm:"column(dtu_no)"`
	MeterAddress int       `orm:"column(meter_address)"`
	CollectTime  time.Time `orm:"column(collect_time)"`
	Rows         int       `orm:"column(rows)"`
}

type CollectCountOfMonth struct {
	CollectTime string `orm:"column(collect_time)"`
	Rows        int    `orm:"column(rows)"`
}

type CollectBaseInfo struct {
	CollectTime  time.Time `orm:"column(collect_time)"`
	DTU_no       string    `orm:"column(dtu_no)"`
	MeterAddress int       `orm:"column(meter_address)"`

	A_electricity float64 `orm:"digits(12);decimals(4);column(a_electricity)"`
	B_electricity float64 `orm:"digits(12);decimals(4);column(b_electricity)"`
	C_electricity float64 `orm:"digits(12);decimals(4);column(c_electricity)"`

	A_power_factor     float64 `orm:"digits(12);decimals(4);column(a_power_factor)"`
	B_power_factor     float64 `orm:"digits(12);decimals(4);column(b_power_factor)"`
	C_power_factor     float64 `orm:"digits(12);decimals(4);column(c_power_factor)"`
	Total_power_factor float64 `orm:"digits(12);decimals(4);column(total_power_factor)"`

	Total_p_at_ee    float64 `orm:"digits(12);decimals(4);column(total_p_at_ee)"`
	Total_r_at_ee    float64 `orm:"digits(12);decimals(4);column(total_r_at_ee)"`
	Total_ap_a_ee    float64 `orm:"digits(12);decimals(4);column(total_ap_a_ee)"`
	Total_ap_reat_ee float64 `orm:"digits(12);decimals(4);column(total_ap_reat_ee)"`

	A_voltage float64 `orm:"digits(12);decimals(4);column(a_voltage)"`
	B_voltage float64 `orm:"digits(12);decimals(4);column(b_voltage)"`
	C_voltage float64 `orm:"digits(12);decimals(4);column(c_voltage)"`

	Total_ap_power float64 `orm:"digits(12);decimals(4);column(total_ap_power)"`
	A_ap_power     float64 `orm:"digits(12);decimals(4);column(a_ap_power)"`
	B_ap_power     float64 `orm:"digits(12);decimals(4);column(b_ap_power)"`
	C_ap_power     float64 `orm:"digits(12);decimals(4);column(c_ap_power)"`

	Total_reactive_power float64 `orm:"digits(12);decimals(4);column(total_reactive_power)"`
	A_reactive_power     float64 `orm:"digits(12);decimals(4);column(a_reactive_power)"`
	B_reactive_power     float64 `orm:"digits(12);decimals(4);column(b_reactive_power)"`
	C_reactive_power     float64 `orm:"digits(12);decimals(4);column(c_reactive_power)"`

	Total_active_power float64 `orm:"digits(12);decimals(4);column(total_active_power)"`
	A_active_power     float64 `orm:"digits(12);decimals(4);column(a_active_power)"`
	B_active_power     float64 `orm:"digits(12);decimals(4);column(b_active_power)"`
	C_active_power     float64 `orm:"digits(12);decimals(4);column(c_active_power)"`

	Total_p_reat_ee float64 `orm:"digits(12);decimals(4);column(total_p_reat_ee)"`
	Total_r_reat_ee float64 `orm:"digits(12);decimals(4);column(total_r_reat_ee)"`
	Total_at_ee     float64 `orm:"digits(12);decimals(4);column(total_at_ee)"`

	Frequency float64 `orm:"digits(12);decimals(4);column(frequency)"`

	Uab_line_voltage float64 `orm:"digits(12);decimals(4);column(uab_line_voltage)"`
	Ubc_line_voltage float64 `orm:"digits(12);decimals(4);column(ubc_line_voltage)"`
	Uac_line_voltage float64 `orm:"digits(12);decimals(4);column(uac_line_voltage)"`
}

type CollectBaseInfoQueryParam struct {
	BaseQueryParam
	CollectTime  string
	DTU_no       string
	MeterAddress string
}

func CollectBaseInfoPageList(params *CollectBaseInfoQueryParam) ([] *CollectBaseInfo, int64) {
	if len(strings.TrimSpace(params.CollectTime)) <= 0 {
		return nil, 0
	}

	if len(strings.TrimSpace(params.MeterAddress)) <= 0 {
		return nil, 0
	}

	beginTime := params.CollectTime + " 00:00:00"
	endTime := params.CollectTime + " 23:59:59"

	data := make([] *CollectBaseInfo, 0)
	o := orm.NewOrm()
	o.Using("kxtimingdata")
	sql := fmt.Sprintf(`SELECT collect_time, dtu_no, meter_address, 
             a_electricity, b_electricity, c_electricity, 
			 a_power_factor, b_power_factor, c_power_factor,
			 total_power_factor, total_p_at_ee, total_r_at_ee, total_ap_a_ee, total_ap_reat_ee,
			 a_voltage, b_voltage, c_voltage,
			 total_ap_power, a_ap_power, b_ap_power, c_ap_power,
			 total_reactive_power, a_reactive_power, b_reactive_power, c_reactive_power,
			 total_active_power, a_active_power, b_active_power, c_active_power,
			 total_p_reat_ee, total_r_reat_ee, total_at_ee,
			 frequency,
			 uab_line_voltage, ubc_line_voltage, uac_line_voltage			 
		FROM collect_base_info 
       where collect_time >= '%s' and collect_time <= '%s'
         and dtu_no like '%s%%'
         and meter_address = %s
       `,
		beginTime,
		endTime,
		params.DTU_no,
		params.MeterAddress,
	)
	total, err := o.Raw(sql).QueryRows(&data)
	if err != nil {
		return nil, 0
	}

	sql = sql + fmt.Sprintf(" LIMIT %d, %d", params.Offset, params.Limit)

	_, err = o.Raw(sql).QueryRows(&data)
	if err != nil {
		return nil, 0
	}
	return data, total
}

func CollectBaseInfoDataList(params *CollectBaseInfoQueryParam) [] *CollectBaseInfo {
	params.Limit = -1
	params.Sort = "collect_time"
	params.Order = "asc"
	data, _ := CollectBaseInfoPageList(params)
	return data
}

//采集进度查询
func GetDtuRowForDayList() ([] *DtuRowOfDay, error) {
	data := make([] *DtuRowOfDay, 0)

	o := orm.NewOrm()
	o.Using("kxtimingdata")
	//sql := "call p_dtu_day_rowtotal('collect_base_info', '')"
	sql := fmt.Sprintf(`SELECT dtu_no, 
                                        meter_address, 
                                        max(collect_time) as collect_time, 
                                        count(1) as rows 
                                   FROM collect_base_info_%s 
                               GROUP BY dtu_no, meter_address
                               `, formatTodayUnderline())
	_, err := o.Raw(sql).QueryRows(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func formatTodayUnderline() string {
	return time.Now().Format("2006_01_02")
}

//取今日采集数量
func GetCollectRowsToday() int {
	var rows int
	o := orm.NewOrm()
	o.Using("kxtimingdata")

	sql := "SELECT count(1) as rows FROM collect_base_info_" + formatTodayUnderline()
	err := o.Raw(sql).QueryRow(&rows)
	if err != nil {
		return 0
	}
	return rows
}

//取月采集数量
func GetCollectRowsOfMonth() ([] *CollectCountOfMonth, error) {
	data := make([] *CollectCountOfMonth, 0)

	o := orm.NewOrm()
	o.Using("kxtimingdata")

	d, _ := time.ParseDuration("-24h")
	yesday := time.Now().Add(d).Format("2006-01-02")
	sql := fmt.Sprintf(`call p_month_rowcount('collect_base_info', '%s')`, yesday)

	_, err := o.Raw(sql).QueryRows(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}