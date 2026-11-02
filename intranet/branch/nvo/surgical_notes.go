package nvo

import (
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

const surgicalNoteEnum = "enum('0','1','2','3','4','5','6','7','8','9','10')"

// SurgicalNoteModel maps the existing surgical_notes table.
type SurgicalNoteModel struct {
	ID                       uint64         `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	NoteType                 *string        `json:"note_type" gorm:"column:note_type;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	Sucursal                 *string        `json:"sucursal" gorm:"column:sucursal;type:varchar(255)"`
	PatientName              *string        `json:"patient_name" gorm:"column:patient_name;type:varchar(255)"`
	PatientRoom              *string        `json:"patient_room" gorm:"column:patient_room;type:varchar(255)"`
	PatientBirthday          *time.Time     `json:"patient_birthday" gorm:"column:patient_birthday;type:date"`
	PatientAge               *int32         `json:"patient_age" gorm:"column:patient_age;type:int"`
	PatientSex               *string        `json:"patient_sex" gorm:"column:patient_sex;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	OnChargeService          *string        `json:"on_charge_service" gorm:"column:on_charge_service;type:varchar(255)"`
	SurgeryDatetime          *time.Time     `json:"surgery_datetime" gorm:"column:surgery_datetime;type:datetime"`
	SurgeryDatetimeEnd       *time.Time     `json:"surgery_datetime_end" gorm:"column:surgery_datetime_end;type:datetime"`
	PreDiagnosis             *string        `json:"pre_diagnosis" gorm:"column:pre_diagnosis;type:longtext"`
	PostDiagnosis            *string        `json:"post_diagnosis" gorm:"column:post_diagnosis;type:longtext"`
	PreSurgicalPlan          *string        `json:"pre_surgical_plan" gorm:"column:pre_surgical_plan;type:longtext"`
	PostSurgicalPlan         *string        `json:"post_surgical_plan" gorm:"column:post_surgical_plan;type:longtext"`
	SurgicalTechnique        *string        `json:"surgical_technique" gorm:"column:surgical_technique;type:longtext"`
	SurgicalFindings         *string        `json:"surgical_findings" gorm:"column:surgical_findings;type:longtext"`
	BleedingCuantification   *string        `json:"bleeding_cuantification" gorm:"column:bleeding_cuantification;type:varchar(255)"`
	IncidentsOrAccidents     *string        `json:"incidents_or_accidents" gorm:"column:incidents_or_accidents;type:longtext"`
	GauzeCount               *string        `json:"gauze_count" gorm:"column:gauze_count;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	CompleteInstrumental     *string        `json:"complete_instrumental" gorm:"column:complete_instrumental;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	SurgicalInterventionType *string        `json:"surgical_intervention_type" gorm:"column:surgical_intervention_type;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	ReIntervention           *string        `json:"re_intervention" gorm:"column:re_intervention;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	ReInterventionReason     *string        `json:"re_intervention_reason" gorm:"column:re_intervention_reason;type:varchar(255)"`
	AmbulatorySurgery        *string        `json:"ambulatory_surgery" gorm:"column:ambulatory_surgery;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	AmbulatorySurgeryReason  *string        `json:"ambulatory_surgery_reason" gorm:"column:ambulatory_surgery_reason;type:varchar(255)"`
	SurgicalRisk             *string        `json:"surgical_risk" gorm:"column:surgical_risk;type:enum('0','1','2','3','4','5','6','7','8','9','10')"`
	PlannedAnesthesia        *string        `json:"planned_anesthesia" gorm:"column:planned_anesthesia;type:varchar(255)"`
	SurgeryMaterial          *string        `json:"surgery_material" gorm:"column:surgery_material;type:longtext"`
	PreSurgeryCareplan       *string        `json:"pre_surgery_careplan" gorm:"column:pre_surgery_careplan;type:longtext"`
	PostSurgeryCareplan      *string        `json:"post_surgery_careplan" gorm:"column:post_surgery_careplan;type:longtext"`
	Forecast                 *string        `json:"forecast" gorm:"column:forecast;type:longtext"`
	PreSurgeryNote           *string        `json:"pre_surgery_note" gorm:"column:pre_surgery_note;type:longtext"`
	AuxiliarStudies          *string        `json:"auxiliar_studies" gorm:"column:auxiliar_studies;type:longtext"`
	Transfusions             *string        `json:"transfusions" gorm:"column:transfusions;type:varchar(255)"`
	Surgeon                  *string        `json:"surgeon" gorm:"column:surgeon;type:varchar(255)"`
	Anesthesiologist         *string        `json:"anesthesiologist" gorm:"column:anesthesiologist;type:varchar(255)"`
	Instrumentalist          *string        `json:"instrumentalist" gorm:"column:instrumentalist;type:varchar(255)"`
	FirstHelper              *string        `json:"first_helper" gorm:"column:first_helper;type:varchar(255)"`
	SecondHelper             *string        `json:"second_helper" gorm:"column:second_helper;type:varchar(255)"`
	Circulating              *string        `json:"circulating" gorm:"column:circulating;type:varchar(255)"`
	OtherMembers             *string        `json:"other_members" gorm:"column:other_members;type:varchar(255)"`
	PostSurgicalState        *string        `json:"post_surgical_state" gorm:"column:post_surgical_state;type:varchar(255)"`
	SurgicalPiecesBiopsy     *string        `json:"surgical_pieces_biopsy" gorm:"column:surgical_pieces_biopsy;type:varchar(255)"`
	OtherFindings            *string        `json:"other_findings" gorm:"column:other_findings;type:varchar(255)"`
	DoctorData               *string        `json:"doctor_data" gorm:"column:doctor_data;type:varchar(255)"`
	DeletedAt                gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp"`
	CreatedAt                *time.Time     `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt                *time.Time     `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
}

// TableName binds the model to the existing Laravel table.
func (SurgicalNoteModel) TableName() string {
	return "surgical_notes"
}

// CreateSurgicalNotesTable creates the table when explicitly called.
// GORM handles database differences, so this helper works with SQLite tests
// and the PostgreSQL production setup.
func CreateSurgicalNotesTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&SurgicalNoteModel{}); err != nil {
		return fmt.Errorf("create surgical_notes table: %w", err)
	}
	return nil
}

func surgicalNoteUpdateMap(note SurgicalNoteModel) map[string]interface{} {
	updates := make(map[string]interface{})
	v := reflect.ValueOf(note)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "ID" || field.Name == "DeletedAt" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" {
			continue
		}
		value := v.Field(i)
		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}
		if value.Kind() == reflect.Ptr && !value.IsNil() {
			updates[field.Tag.Get("json")] = value.Elem().Interface()
			continue
		}
		if !value.IsZero() {
			updates[field.Tag.Get("json")] = value.Interface()
		}
	}
	return updates
}
