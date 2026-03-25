![Alt text](go_project_structure_flow.svg)


install tool

cmd
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

# Trivy
scoop bucket add extras
scoop install trivy

# Gitleaks
scoop install gitleaks

# Semgrep 
pip install semgrep

# gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest
