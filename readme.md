# Black Dog
## Directory Structure
```mermaid
blackdog/
├── .gitignore                   
├── go.mod                       
├── makefile                     
├── readme.md                    
├── bin/
│   └── blackdog-web             
├── cmd/
│   └── web/
│       └── main.go              # PACKAGE: main (Composition Root)
├── data/
│   ├── blackdog.db              # SQLite Database
│   └── scripts/                 # Database Migration Scripts
├── internal/
│   ├── features/
│   │   └── application/         # PACKAGE: application (Vertical Slice)
│   │       ├── models.go        # Domain Entities
│   │       ├── repository.go    # Repository Interface (Port)
│   │       ├── service.go       # Application Logic (Use Cases)
│   │       ├── handlers.go      # HTTP Handlers (Adapter)
│   │       └── dtos.go          # DTOs
│   └── infra/                   # PACKAGE: infra
│       └── memory_repo.go       # In-Memory Database (Adapter)
└── ui/
    ├── html/                    # Templates
    │   ├── base.tmpl            
    │   ├── pages/
    │   └── partials/
    └── static/                  # Assets
        ├── css/
        └── js/
```
## Module Packages
```mermaid
classDiagram
namespace application {
%% AGGREGATE ROOT
class Application {
+ID int64
+ApplicationNumber string
+CategoryCode string
+StatusCode string
+RequestedAmount int64
+CreatedAt time.Time
+UpdatedAt time.Time
+LoanDetail LoanDetail
+CardDetail CardDetail
+Applicants []Applicant
+Validate() error
+Submit() error
}

%% ENTITIES & VALUE OBJECTS
class Applicant {
+ID int64
+ApplicationID int64
+RoleCode string
+ProductTypeCode string
+FirstName string
+LastName string
+DateOfBirth string
+Identifications []Identification
+Addresses []Address
+Contacts []Contact
+Employments []Employment
+Educations []Education
+Validate() error
}

class Id {
+ID int64
+TypeCode string
+Value string
+IssuingAuthority string
+ExpiryDate string
}

class Address {
+ID int64
+TypeCode string
+Line1 string
+Line2 string
+CityMunicipality string
+Province string
+PostalCode string
}

class Contact {
+ID int64
+TypeCode string
+ContactNumber string
}

class Employment {
+ID int64
+StatusCode string
+EmployerName string
+GrossMonthlyIncome int64
}

class Education {
+ID int64
+LevelCode string
+InstitutionName string
}

class LoanDetail {
+TermMonths int
+InterestRate int
+LoanAmount int64
+Validate() error
}

class CardDetail {
+CardDesignCode string
+CreditLimit int64
+Validate() error
}

%% SERVICES
class Service {
-repo Repository
+CreateApplication(ctx, cmd CreateApplicationDTO)
+GetApplication(ctx, id int64)
}

class Repository {
<<Interface>>
+Save(ctx, app Application) error
+GetByID(ctx, id int64) (Application, error)
}

class Handler {
-svc *Service
-tmpls map[string]*template.Template
+HandleCreate(w, r)
+HandleView(w, r)
}
}

namespace infra_Package {
class InMemoryApplicationRepo {
-store map[int64]*application.Application
-mu sync.RWMutex
+Save(ctx, app Application) error
+GetByID(ctx, id int64) (Application, error)
}
}

%% RELATIONSHIPS
Application *-- Applicant : Contains (Slice)
Application *-- LoanDetail : Has A
Application *-- CardDetail : Has A
Applicant *-- Id : Contains (Slice)
Applicant *-- Address : Contains (Slice)
Applicant *-- Contact : Contains (Slice)
Applicant *-- Employment : Contains (Slice)
Applicant *-- Education : Contains (Slice)
Service --> Repository : Uses
Handler --> Service : Uses
InMemoryApplicationRepo ..|> Repository : Implements
```

