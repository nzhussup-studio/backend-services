package model

type PersonalData struct {
	WorkExperience []*WorkExperience `json:"work_experience"`
	Education      []*Education      `json:"education"`
	Projects       []*Project        `json:"projects"`
	Skills         []*Skill          `json:"skills"`
	Certificates   []*Certificate    `json:"certificates"`
}
