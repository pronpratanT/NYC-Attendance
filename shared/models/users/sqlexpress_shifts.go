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
	PRS_WELFARE_D time.Time `gorm:"column:PRS_WELFARE_D"`

	PRS_2FN_NO  string    `gorm:"column:PRS_2FN_NO"`
	PRS_2FN_DD  time.Time `gorm:"column:PRS_2FN_DD"`
	PRS_2FN_DDE time.Time `gorm:"column:PRS_2FN_DDE"`
	PRS_2FN_DDC time.Time `gorm:"column:PRS_2FN_DDC"`

	PRS_FN_SCHEME_CD string `gorm:"column:PRS_FN_SCHEME_CD"`
	PRS_FN_SUBFN1_CD string `gorm:"column:PRS_FN_SUBFN1_CD"`
	PRS_FN_SUBFN2_CD string `gorm:"column:PRS_FN_SUBFN2_CD"`
	PRS_FN_SUBFN3_CD string `gorm:"column:PRS_FN_SUBFN3_CD"`

	PRS_FN_SUBFN1_E float64 `gorm:"column:PRS_FN_SUBFN1_E"`
	PRS_FN_SUBFN2_E float64 `gorm:"column:PRS_FN_SUBFN2_E"`
	PRS_FN_SUBFN3_E float64 `gorm:"column:PRS_FN_SUBFN3_E"`

	PRS_FN_SUBFN1_C float64 `gorm:"column:PRS_FN_SUBFN1_C"`
	PRS_FN_SUBFN2_C float64 `gorm:"column:PRS_FN_SUBFN2_C"`
	PRS_FN_SUBFN3_C float64 `gorm:"column:PRS_FN_SUBFN3_C"`

	PRS_2FN_SCHEME_CD string `gorm:"column:PRS_2FN_SCHEME_CD"`
	PRS_2FN_SUBFN1_CD string `gorm:"column:PRS_2FN_SUBFN1_CD"`
	PRS_2FN_SUBFN2_CD string `gorm:"column:PRS_2FN_SUBFN2_CD"`
	PRS_2FN_SUBFN3_CD string `gorm:"column:PRS_2FN_SUBFN3_CD"`

	PRS_2FN_SUBFN1_E float64 `gorm:"column:PRS_2FN_SUBFN1_E"`
	PRS_2FN_SUBFN2_E float64 `gorm:"column:PRS_2FN_SUBFN2_E"`
	PRS_2FN_SUBFN3_E float64 `gorm:"column:PRS_2FN_SUBFN3_E"`

	PRS_2FN_SUBFN1_C float64 `gorm:"column:PRS_2FN_SUBFN1_C"`
	PRS_2FN_SUBFN2_C float64 `gorm:"column:PRS_2FN_SUBFN2_C"`
	PRS_2FN_SUBFN3_C float64 `gorm:"column:PRS_2FN_SUBFN3_C"`

	PRS_EWF_D  time.Time `gorm:"column:PRS_EWF_D"`
	PRS_EWF_NO string    `gorm:"column:PRS_EWF_NO"`
}

func (SQLExpressUser) TableName() string {
	return "PERSONALINFO"
}
