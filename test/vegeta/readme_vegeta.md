# Pruebas de Carga con Vegeta

Este proyecto utiliza Vegeta para realizar pruebas de carga sobre el backend.

## Ejecución del test

1. Asegurate de que el backend esté corriendo y accesible (por ejemplo en `http://127.0.0.1:8080`).

2. Ejecutá el siguiente comando en la terminal para lanzar la prueba:

```bash
vegeta attack -rate=60 -duration=10s -targets=./test_carga.txt | tee resultado.bin | vegeta report | tee resultados.txt
```

- rate=60: envía 60 peticiones por segundo

- duration=10s: duración total de la prueba

- resultado.bin: guarda las respuestas crudas

- resultados.txt: genera un reporte legible en texto

Una vez finalizada la prueba, podés inspeccionar los resultados en el archivo resultados.txt.