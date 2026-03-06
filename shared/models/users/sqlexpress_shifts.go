package model

import "time"

type SQLExpressShifts struct {
	SFKey         int       `json:"sf_key" gorm:"primaryKey;column:SF_KEY"`
	SFCode        string    `json:"sf_code" gorm:"column:SF_CODE"`
	SFName        string    `json:"sf_name" gorm:"column:SF_NAME"`
	SFEName       string    `json:"sf_e_name" gorm:"column:SF_E_NAME"`
	SFType        int16     `json:"sf_type" gorm:"column:SF_TYPE"`
	SFDf          int16     `json:"sf_df" gorm:"column:SF_DF"`
	SF1InDay      int16     `json:"sf_1in_day" gorm:"column:SF_1IN_DAY"`
	SF1InTime     time.Time `json:"sf_1in_time" gorm:"column:SF_1IN_TIME"`
	SFLInDay      int16     `json:"sf_lin_day" gorm:"column:SF_LIN_DAY"`
	SFLInTime     time.Time `json:"sf_lin_time" gorm:"column:SF_LIN_TIME"`
	SFInDay       int16     `json:"sf_in_day" gorm:"column:SF_IN_DAY"`
	SFInTime      time.Time `json:"sf_in_time" gorm:"column:SF_IN_TIME"`
	SFOutDay      int16     `json:"sf_out_day" gorm:"column:SF_OUT_DAY"`
	SFOutTime     time.Time `json:"sf_out_time" gorm:"column:SF_OUT_TIME"`
	SFBreak       string    `json:"sf_break" gorm:"column:SF_BREAK"`
	SFBrkiDay     int16     `json:"sf_brki_day" gorm:"column:SF_BRKI_DAY"`
	SFBrkiTime    time.Time `json:"sf_brki_time" gorm:"column:SF_BRKI_TIME"`
	SFBrkoDay     int16     `json:"sf_brko_day" gorm:"column:SF_BRKO_DAY"`
	SFBrkoTime    time.Time `json:"sf_brko_time" gorm:"column:SF_BRKO_TIME"`
	SFNstminDf    int       `json:"sf_nstmin_df" gorm:"column:SF_NSTMIN_DF"`
	SFNstmotDf    int       `json:"sf_nstmot_df" gorm:"column:SF_NSTMOT_DF"`
	SFAbsIndex    int       `json:"sf_abs_index" gorm:"column:SF_ABS_INDEX"`
	SFNstmioDf    int       `json:"sf_nstmio_df" gorm:"column:SF_NSTMIO_DF"`
	SFEnabled     string    `json:"sf_enabled" gorm:"column:SF_ENABLED"`
	SFRemark      string    `json:"sf_remark" gorm:"column:SF_REMARK"`
	SFHrsWork     float64   `json:"sf_hrs_work" gorm:"column:SF_HRS_WORK"`
	SFHrsBreak    float64   `json:"sf_hrs_break" gorm:"column:SF_HRS_BREAK"`
	SFNoStamp     int16     `json:"sf_no_stamp" gorm:"column:SF_NO_STAMP"`
	SFFixTime     string    `json:"sf_fix_time" gorm:"column:SF_FIX_TIME"`
	SFPublHoliday string    `json:"sf_publ_holiday" gorm:"column:SF_PUBL_HOLIDAY"`
	SFGuid        string    `json:"sf_guid" gorm:"column:SF_GUID"`
}

func (SQLExpressShifts) TableName() string {
	return "TMSHIFT"
}

type SQLExpressUser struct {
	EmpAccess       int       `json:"emp_access" gorm:"column:EMP_ACCESS"`
	EmpAddr1        string    `json:"emp_addr_1" gorm:"column:EMP_ADDR_1"`
	EmpAddr2        string    `json:"emp_addr_2" gorm:"column:EMP_ADDR_2"`
	EmpAddr3        string    `json:"emp_addr_3" gorm:"column:EMP_ADDR_3"`
	EmpAddrCountry  string    `json:"emp_addr_country" gorm:"column:EMP_ADDR_COUNTRY"`
	EmpAddrDistrict string    `json:"emp_addr_district" gorm:"column:EMP_ADDR_DISTRICT"`
	EmpAddrProvince string    `json:"emp_addr_province" gorm:"column:EMP_ADDR_PROVINCE"`
	EmpAddrSubDist  string    `json:"emp_addr_sub_district" gorm:"column:EMP_ADDR_SUB_DISTRICT"`
	EmpAlert        int       `json:"emp_alert" gorm:"column:EMP_ALERT"`
	EmpAlertMsg     *string   `json:"emp_alert_msg" gorm:"column:EMP_ALERT_MSG"`
	EmpBirth        time.Time `json:"emp_birth" gorm:"column:EMP_BIRTH"`
	EmpEmail        string    `json:"emp_email" gorm:"column:EMP_EMAIL"`
	EmpEName        string    `json:"emp_e_name" gorm:"column:EMP_E_NAME"`
	EmpGender       int       `json:"emp_gender" gorm:"column:EMP_GENDER"`
	EmpIntl         string    `json:"emp_intl" gorm:"column:EMP_INTL"`
	EmpICard        string    `json:"emp_i_card" gorm:"column:EMP_I_CARD"`
	EmpIExpire      time.Time `json:"emp_i_expire" gorm:"column:EMP_I_EXPIRE"`
	EmpIIssue       string    `json:"emp_i_issue" gorm:"column:EMP_I_ISSUE"`
	EmpKey          int       `json:"emp_key" gorm:"primaryKey;column:EMP_KEY"`
	EmpMarital      string    `json:"emp_marital" gorm:"column:EMP_MARITAL"`
	EmpName         string    `json:"emp_name" gorm:"column:EMP_NAME"`
	EmpNotiID       string    `json:"emp_noti_id" gorm:"column:EMP_NOTI_ID"`
	EmpPost         string    `json:"emp_post" gorm:"column:EMP_POST"`
	EmpRemark       *string   `json:"emp_remark" gorm:"column:EMP_REMARK"`
	EmpScPrx        int       `json:"emp_sc_prx" gorm:"column:EMP_SC_PRX"`
	EmpSlipMsg      *string   `json:"emp_slip_msg" gorm:"column:EMP_SLIP_MSG"`
	EmpSlipPw       *string   `json:"emp_slip_pw" gorm:"column:EMP_SLIP_PW"`
	EmpSurnme       string    `json:"emp_surnme" gorm:"column:EMP_SURNME"`
	EmpTaxID        *string   `json:"emp_tax_id" gorm:"column:EMP_TAX_ID"`
	EmpTel          string    `json:"emp_tel" gorm:"column:EMP_TEL"`
}

func (SQLExpressUser) TableName() string {
	return "EMPFILE"
}

type SQLExpressMasterKey struct {
	TmrBr     int       `json:"tmr_br" gorm:"column:TMR_BR"`
	TmrDate   time.Time `json:"tmr_date" gorm:"column:TMR_DATE"`
	TmrDept   int       `json:"tmr_dept" gorm:"column:TMR_DEPT"`
	TmrDf     int       `json:"tmr_df" gorm:"column:TMR_DF"`
	TmrDfApr  int       `json:"tmr_df_apr" gorm:"column:TMR_DF_APR"`
	TmrDfPr   *int      `json:"tmr_df_pr" gorm:"column:TMR_DF_PR"`
	TmrDfT    int       `json:"tmr_df_t" gorm:"column:TMR_DF_T"`
	TmrEmp    int       `json:"tmr_emp" gorm:"column:TMR_EMP"`
	TmrKey    int       `json:"tmr_key" gorm:"primaryKey;column:TMR_KEY"`
	TmrQty    string    `json:"tmr_qty" gorm:"column:TMR_QTY"`
	TmrQtyApr string    `json:"tmr_qty_apr" gorm:"column:TMR_QTY_APR"`
	TmrQtyT   string    `json:"tmr_qty_t" gorm:"column:TMR_QTY_T"`
	TmrSctn   int       `json:"tmr_sctn" gorm:"column:TMR_SCTN"`
	TmrSf     int       `json:"tmr_sf" gorm:"column:TMR_SF"`
	TmrSite   int       `json:"tmr_site" gorm:"column:TMR_SITE"`
	TmrSs     int       `json:"tmr_ss" gorm:"column:TMR_SS"`
	TmrStt    int       `json:"tmr_stt" gorm:"column:TMR_STT"`
	TmrTev    int       `json:"tmr_tev" gorm:"column:TMR_TEV"`
	TmrTrap   int       `json:"tmr_trap" gorm:"column:TMR_TRAP"`
}

func (SQLExpressMasterKey) TableName() string {
	return "TMRESULT"
}
