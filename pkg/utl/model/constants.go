package model

const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
	BUCKET                = "BUCKET"
	BUCKET_PRIVATE_PREFIX = "BUCKET_PRIVATE_PREFIX"
	BUCKET_PUBLIC_PREFIX  = "BUCKET_PUBLIC_PREFIX"
	AWS_REGION            = "AWS_REGION"
	SES_SENDER            = "SES_SENDER"
	FQDN                  = "FQDN"
	COMPANY               = "COMPANY"
	ADMIN_EMAIL           = "ADMIN_EMAIL"
	ADMIN_PASSWORD        = "ADMIN_PASSWORD"
	ENVIRONMENT           = "ENVIRONMENT"
	JWT_SECRET            = "JWT_SECRET"
	DB_HOST               = "DB_HOST"
	DB_PORT               = "DB_PORT"
	DB_USER               = "DB_USER"
	DB_PASSWORD           = "DB_PASSWORD"
	DB_NAME               = "DB_NAME"
	APP_PORT              = "APP_PORT"
	CDN                   = "CDN"

	bulkInsert = 100

	AdminUser  = "5acf6e54-0e38-421e-b9ee-2cad164bcfe2"
	AdminEmail = "system@system.com"

	MinuteCron = "minute"
	HourlyCron = "hourly"
	DailyCron  = "daily"

	GeneralFile               = "Documento"
	ImageFile                 = "Imagem"
	ProfileImage              = "Imagem de Perfil"
	CoverPhoto                = "Imagem de Capa"
	IDScan                    = "Bilhete de Identidade"
	PassportScan              = "Passaporte"
	DiplomaScan               = "Diploma"
	CertificateScan           = "Certificado"
	CertificateValidationScan = "Certificado Homologado"
	EnrollmentProof           = "Declaração de Frequência"
	ResidenceProof            = "Certificado de Residência"
	HealthCertificate         = "Declaração de Saúde"
	ProjectProposal           = "Proposta de Projecto"
	PovertyDeclaration        = "Declaração de Pobreza"
	RPEProof                  = "Comprovativo de RPE"
	RecommendationLetter      = "Carta de Recomendação"
	ReturnAgreement           = "Compromisso de Retorno"
	WorkDeclaration           = "Declaração de Trabalho"
	IESAppraisal              = "Parecer da IES"
	MilitaryDocument          = "Documento Militar"

	GenderMale      = "Masculino"
	GenderFemale    = "Feminino"
	GenderNonBinary = "Não-Binário"

	RelationFather = "Pai"
	RelationMother = "Mãe"
	RelationSpouse = "Cônjuge"

	HighSchool     = "Ensino Médio"
	Grad           = "Licenciatura"
	GradOther      = "Graduação"
	PostGrad       = "Pós-Graduação"
	Masters        = "Mestrado"
	Doctorate      = "Doutoramento"
	Specialization = "Especialidade"

	NormalCourse   = "Normal"
	PriorityCourse = "Prioritário"

	GradeNA          = "Não Aplicável"
	GradeDistinction = "Bom com distinção"
	GradeVeryGood    = "Muito Bom"
	GradeExcelent    = "Excelente"

	StatusActive    = "Activo"
	StatusSuspended = "Suspenso"
	StatusInactive  = "Inactivo"
	StatusCanceled  = "Cancelado"

	StatusDraft                = "Rascunho"
	StatusOpen                 = "Aberta"
	StatusClosed               = "Encerrada"
	StatusConcludedAppointment = "Concluído"
	StatusCompleted            = "Concluída"
	StatusTrash                = "Lixeira"
	StatusPublished            = "Publicado"
	StatusScholarshipPublished = "Publicada"

	StatusValidated           = "Validada"
	StatusApproved            = "Aprovado"
	StatusCanceledApplication = "Cancelada"
	StatusApprovedApplication = "Aprovada"
	StatusAwarded             = "Atribuída"
	StatusDuplicate           = "Duplicada"
	StatusNeedsReview         = "Precisa de Revisão"
	StatusRejected            = "Rejeitado"
	StatusRejectedApplication = "Rejeitada"
	StatusPending             = "Pendente"
	StatusError               = "Erro"
	StatusProcessing          = "A Processar"

	FileS3    = "S3"
	FileLocal = "Local"

	OrgSchool  = "Escola"
	OrgSponsor = "Patrocinador"
	OrgMain    = "Principal"

	LanguageWeak     = "Fraco"
	LanguageMedium   = "Suficiente"
	LanguageGood     = "Bom"
	LanguageExcelent = "Excelente"

	InternalScholarship = "Interna"
	ExternalScholarship = "Externa"
	MeritScholarship    = "Externa de Mérito"

	ResourceScholarship   = "Bolsa"
	ResourcePost          = "Artigo"
	ResourceUser          = "Utilizador"
	ResourceApplication   = "Candidatura"
	ResourceAppointment   = "Audiência"
	ResourceRole          = "Função"
	ResourceFile          = "Documento"
	ResourceProvinceQuota = "Quota Provincial"
	ResourceStatistic     = "Estatística"
	ResourceNotification  = "Notificação"
	ResourceRecipient     = "Recipiente"
	ResourceCity          = "Cidade"
	ResourceCourse        = "Curso"
	ResourceLanguages     = "Línguas"
	ResourceFamily        = "Família"
	ResourceOrganization  = "Organização"

	GlobalStat   = "Global"
	ProvinceStat = "Província"

	TemplateAction  = "Acção"
	TemplateMessage = "Mensagem"
	TemplateContact = "Contacto"

	//File Access Type
	AccessPublic  = "Public"
	AccessPrivate = "Private"

	//Post Categories
	CategoryNews         = "Notícia"
	CategoryAnnouncement = "Anúncio"
	CategoryGeneral      = "Geral"
	Epoch                = "1970-01-01 00:00:00"
)

var (
	assetsDir     string
	DocumentTypes = List{GeneralFile, ProfileImage, CoverPhoto, IDScan, PassportScan, DiplomaScan, CertificateScan, CertificateValidationScan, EnrollmentProof, ResidenceProof, HealthCertificate, ProjectProposal, PovertyDeclaration, RPEProof, RecommendationLetter, ReturnAgreement, WorkDeclaration, IESAppraisal, MilitaryDocument}

	requiredDocuments = map[string]List{
		InternalScholarship + Grad:           List{IDScan, CertificateScan, EnrollmentProof},
		InternalScholarship + Masters:        List{IDScan, CertificateScan, CertificateValidationScan, EnrollmentProof},
		InternalScholarship + Specialization: List{IDScan, CertificateScan, CertificateValidationScan, EnrollmentProof},
		InternalScholarship + Doctorate:      List{IDScan, CertificateScan, CertificateValidationScan, EnrollmentProof},

		ExternalScholarship + Grad:           List{IDScan, DiplomaScan, CertificateScan, CertificateValidationScan},
		ExternalScholarship + Masters:        List{IDScan, DiplomaScan, CertificateScan, CertificateValidationScan},
		ExternalScholarship + Specialization: List{IDScan, DiplomaScan, CertificateScan, CertificateValidationScan},
		ExternalScholarship + Doctorate:      List{IDScan, DiplomaScan, CertificateScan, CertificateValidationScan},

		MeritScholarship + Masters:        List{IDScan, PassportScan, DiplomaScan, CertificateScan, CertificateValidationScan, RecommendationLetter, ReturnAgreement},
		MeritScholarship + Specialization: List{IDScan, PassportScan, DiplomaScan, CertificateScan, CertificateValidationScan, RecommendationLetter, ReturnAgreement},
		MeritScholarship + Doctorate:      List{IDScan, PassportScan, DiplomaScan, CertificateScan, CertificateValidationScan, RecommendationLetter, ReturnAgreement, ProjectProposal},
	}

	applicationCSVHeaders = []string{
		//Identity
		"Ref", "Estado", "Nome", "Bilhete de Identidade", "Passaporte", "Regime", "Nota", "Sexo", "Telefone", "Email", "Skype",

		//Scholarship
		"Ref. Bolsa", "Bolsa", "Nível da Bolsa", "Tipo da Bolsa",

		//Birth and Residence
		"Data de Nascimento", "Província de Origem", "País de Origem", "Província de Residência", "País de Residência",

		//Bank
		"Banco", "Conta Bancária", "Titular da Conta", "IBAN",

		//Employment
		"Empregador", "Endereço do Empregador", "Profissão", "Salário (AKZ/Mês)",

		//Family
		"Membros de Família", "Membros Trabalhadores", "Menores",
		"Nome do Pai", "Naturalidade do Pai", "Profissão do Pai", "Empregador do Pai", "Salário do Pai (AKZ/Mês)",
		"Nome da Mãe", "Naturalidade da Mãe", "Profissão da Mãe", "Empregador da Mãe", "Salário da Mãe (AKZ/Mês)",
		"Nome do Cônjuge", "Naturalidade do Cônjuge", "Profissão do Cônjuge", "Empregador do Cônjuge", "Salário do Cônjuge (AKZ/Mês)",

		//Education
		"Grau Enterior", "Instituição Anterior", "Departamento/UO Anterior", "Curso Anterior", "País do Curso Anterior", "Província do Curso Anterior",
		"Data de Graduação", "Média Final", "Avaliação Qualitativa",

		//Application
		"Instituição", "Unidade Orgânica", "Opção 1", "Opção 2", "Opção 3", "País de Ensino", "Província de Ensino", "Ano de frequência", "Média",

		//Languages
		"Francês", "Espanhol", "Inglês", "Italiano", "Alemão", "Russo", "Mandarin", "Outra Língua",

		//Documents
		"Cópia do BI", "Cópia do Passaporte", "Diploma", "Certificado", "Certificado Homologado", "Declaração de Frequência", "Certificado de Residência",
		"Declaração de Saúde", "Proposta de Projecto", "Declaração de Pobreza", "Comprovativo de RPE", "Carta de Recomendação", "Compromisso de Retorno",
		"Declaração de Trabalho", "Parecer da IES", "Documento Militar",
	}
)
