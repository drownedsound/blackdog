# Black Dog
**Black Dog** is an exercise in building a secure and high-performance web application designed for processing consumer applications for common retail banking products. Currently, the domain model focuses on credit cards and personal loans but other products are in the immediate roadmap. It employs Domain-Driven Design (DDD) to model the problem space but does not attempt to create microservices in the process. Instead, it aims to deliver the system with a monolithic design and fewer moving pieces with the objective of reducing overall system and cognitive complexity.
## Project Structure
```text
blackdog/
├── bin/                 # Compiled binaries
├── cmd/
│   └── web/             # Main entry point (composition root)
├── data/                # Database migration scripts
├── internal/
│   ├── features/        # Feature-based packages 
│   │   └── application/ # Core domain logic 
│   └── infra/           # Infrastructure adapters 
└── ui/                  # Web assets 
````
## Getting Started
### Prerequisites
- **Go**: Version 1.25.4 or higher.
- **Make**: For running build automation commands.
- **Tooling**: The project relies on standard Go tools for formatting and linting.
```shell
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
```
### Installation
1. **Clone the repository:**
```bash
git clone https://github.com/drownedsound/blackdog.git
cd blackdog
```
2. **Build the project:**
Use the provided `makefile` to run static analysis, tests, and compilation in one step.
```bash
make build
```
This command runs `goimports`, `gofumpt`, `staticcheck`, and `go test` before building.
### Running the Application
After a successful build, the binary is placed in `bin/blackdog-web`.
```bash
./bin/blackdog-web
```
By default, the server listens on `:4000` and serves static assets from `./ui/static/`.
### Configuration
The application accepts command-line flags to override default settings:
| Flag | Default | Description |
| :--- | :--- | :--- |
| `-addr` | `:4000` | HTTP network address to listen on. |
| `-static-dir` | `./ui/static/` | Path to the directory containing static assets. |

**Example:**
```bash
./bin/blackdog-web -addr=":8080" -static-dir="/var/www/blackdog/static"
```
## Development
### Running Tests
To run unit tests specifically for the application domain and infrastructure:
```bash
make test
```
This targets `./internal/features/application` and `./internal/infra` with coverage enabled.
### Code Quality
Enforce code standards and formatting:
```bash
make staticcheck  # Runs static analysis
make gofumpt      # Enforces stricter formatting
```
### Domain Model
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
pplication *-- LoanDetail : Has A
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
The core domain entity is the **Application**, which acts as the Aggregate Root. It manages the state and invariants of:
- **Applicant**: Personal details, employment, and identification.
- **LoanDetail**: Terms, interest rates, and amounts.
- **CardDetail**: Credit limits and card designs.
- **Status**: Tracks lifecycle states (e.g., Created, InProgress, Approved).

Current **Status Codes** support the following lifecycle:
`CREATED` -\> `IN_PROGRESS` -\> `APPROVED` | `DECLINED` | `CANCELLED`

