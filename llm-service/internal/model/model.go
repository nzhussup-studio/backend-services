package model

import "time"

type PersonalData struct {
	WorkExperience []*WorkExperience `json:"work_experience"`
	Education      []*Education      `json:"education"`
	Projects       []*Project        `json:"projects"`
	Skills         []*Skill          `json:"skills"`
	Certificates   []*Certificate    `json:"certificates"`
}

type WorkExperience struct {
	ID           int    `json:"id"`
	Company      string `json:"company"`
	Location     string `json:"location"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Position     string `json:"position"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"displayOrder"`
	TechStack    string `json:"techStack"`
}

type Skill struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	SkillNames   string `json:"skillNames"`
	DisplayOrder int    `json:"displayOrder"`
}

type Project struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	TechStack    string `json:"techStack"`
	DisplayOrder int    `json:"displayOrder"`
}

type Education struct {
	ID           int        `json:"id"`
	Institution  string     `json:"institution"`
	Location     string     `json:"location"`
	StartDate    time.Time  `json:"startDate"`
	EndDate      *time.Time `json:"endDate,omitempty"`
	Degree       string     `json:"degree"`
	Thesis       string     `json:"thesis"`
	Description  *string    `json:"description,omitempty"`
	DisplayOrder int        `json:"displayOrder"`
}

type Certificate struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	DisplayOrder int    `json:"displayOrder"`
}

type SummarizerAPIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
}
