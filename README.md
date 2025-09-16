# Backend Projeto e Desenvolvimento de Sistemas - 2025.1 | CC - UFAL

## Grupo 3

### Discentes

- ANTONIO MARIA CARDOSO WAGNER
- BRUNO HENRIQUE SILVA ROCHA
- CAIO AGRA LEMOS
- KAIO VITOR NABUCO DE MELLO SILVA
- VICTOR ALEXANDRE DA ROCHA MONTEIRO MIRANDA

### Docente orientador

- RANILSON OSCAR ARAUJO PAIVA

## Preparação e execução do projeto

```bash
git clone https://github.com/VicAlexandre/pds-backend
go get -u ./... 
docker-compose -f scripts/docker-compose.yml up -d
go run scripts/migration.go
go run cmd/api/main.go
```
