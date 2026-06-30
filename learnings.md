# Learnings from building `gomemo`

## 1. Clean layering makes iteration easier
Splitting responsibilities into **API/router**, **notes service/handler**, and **storage** made it straightforward to change one area without breaking others. The `RouteRegistrar` pattern also kept route registration modular and explicit.

## 2. Versioning API routes early reduces future migration pain
Using `/api/v1` for business endpoints while keeping `/health` unversioned is a practical balance: clients can migrate by version, while monitoring stays stable.

## 3. Interfaces and generics improve flexibility
Using a `storage.Storage[*Note]` interface plus a generic in-memory implementation made core logic independent of persistence details, which is useful for swapping storage backends later.

## 4. Validation belongs in the domain workflow
Validating notes before create/update operations prevented invalid state from entering storage and kept correctness checks close to business rules.

## 5. Good defaults speed up local development
Auto-seeding development data and supporting `ENV=dev` made it faster to test flows manually and reduced setup friction for each run.

## 6. Testing strategy matters more than test count
Combining unit tests across packages with in-memory integration tests for `/api/v1/notes` gave confidence in both isolated logic and end-to-end behavior.

## 7. Tooling and release automation pay off
Simple Make targets (`build`, `run`, `test`, `release`) and embedded version info from git tags made day-to-day development and releases more predictable.
