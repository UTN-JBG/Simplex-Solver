# Simplex-Solver

## Descripción

Aplicación web que tiene como objetivo resolver problemas de programación lineal (maximización y minimización) utilizando el Método Simplex.

Permite ingresar los coeficientes de la función objetivo y las restricciones, y devuelve el resultado óptimo junto con las tablas intermedias del proceso de cálculo.

## Integrantes del grupo 
- Julieta Chaki
- María Guadalupe Cuartara 
- María Belén Sarome

## Tecnologías
- **Backend:** Go (Gin)
- **Frontend:** React
- **Gestión de dependencias (Go):** Go Modules (go mod)
- **Gestión de dependencias (React):** npm

## Metodología Ágil
El equipo aplica **Scrum** como marco de trabajo:
- **Sprints** de 2 semanas.
- **Daily Meetings** de seguimiento.
- **Product Backlog** gestionado en GitHub Projects.
- **Revisiones de Sprint** para validar entregables.
- **Retrospectivas** para la mejora continua.


## Instalación y ejecución local

### Requisitos previos

Asegurate de tener instaladas las siguientes herramientas:
 - Go. Versión recomendada: Go 1.22 o superior.
Descarga e instalación: https://go.dev/dl/

- Node.js y npm. Node.js versión 18 o superior. npm se instala automáticamente junto con Node.js
- Git. Para clonar el repositorio desde GitHub. Descarga: https://git-scm.com/

### 1. Clonar repositorio
```bash
git clone https://github.com/UTN-JBG/Simplex-Solver.git
cd Simplex-Solver
```

### 2. Ejecutar el backend con Go
```
go mod tidy
go run main.go
```

### Ejecución de Tests
El proyecto incluye tests unitarios y de integración para validar la lógica del método Simplex.
Para ejecutar todos los tests del backend, ubicarse en el directorio backend `(cd backend)` y correr:
```
go test ./... -v
```

## 3. Levantar el frontend
Para instalar dependencias, dentro del directorio *frontend* ejecutar:
```
npm install 
```
Para levantar la app
```
npm start
```
