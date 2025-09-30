# Simplex-Solver

## Descripción
Aplicación web que tiene como objetivo resolver problemas de programación lineal utilizando el Método Simplex. 

## Integrantes del grupo 
- Julieta Chaki
- María Guadalupe Cuartara 
- María Belén Sarome

## Tecnologías
- **Lenguaje:** Go
- **Framework web:** Gin
- **Gestión de dependencias:** Go Modules (`go mod`)

## Metodología Ágil
El equipo aplica **Scrum** como marco de trabajo:
- **Sprints** de 2 semanas.
- **Daily Meetings** de seguimiento.
- **Product Backlog** gestionado en GitHub Projects.
- **Revisiones de Sprint** para validar entregables.
- **Retrospectivas** para la mejora continua.

## Instalación y ejecución

### 1. Clonar repositorio
```bash
git clone https://github.com/UTN-JBG/Simplex-Solver.git
cd Simplex-Solver
```

### 2. Ejecutar con Go
```
go mod tidy
go run main.go
```
## 3. Desde Postman probar ruta
```
http://localhost:8080/api/simplex
```
En Body, seleccionar raw y copiar JSON
```
{
  "objective": [4, 3, 2, 1],
  "constraints": [
    [2, 1, 0, 3],
    [1, 3, 2, 1],
    [0, 2, 1, 4]
  ],
  "rhs": [5, 11, 8],
  "type": "max"
}


```
