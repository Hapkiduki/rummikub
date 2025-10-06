package main

import "testing"

func TestEsTrioValido(t *testing.T) {
	// Implementamos Table Driven Tests
	casosDePrueba := []struct {
		nombre   string  // Nombre descriptivo del caso de prueba
		fichas   []Pieza // El input para nuestra función
		esperado bool    // el resultado que esperamos obtener
	}{
		{
			nombre: "Trio valido de 3 fichas con comodin",
			fichas: []Pieza{
				{Color: Rojo, Numero: 7},
				{Color: Azul, Numero: 7},
				{Color: -1, Numero: 0}, // Comodín
			},
			esperado: true,
		},
		{
			nombre: "Trio valido de 3 fichas",
			fichas: []Pieza{
				{Color: Rojo, Numero: 7},
				{Color: Azul, Numero: 7},
				{Color: Negro, Numero: 7},
			},
			esperado: true,
		},
		{
			nombre: "Cuarteto válido de 4 fichas",
			fichas: []Pieza{
				{Color: Rojo, Numero: 10},
				{Color: Azul, Numero: 10},
				{Color: Amarillo, Numero: 10},
				{Color: Negro, Numero: 10},
			},
			esperado: true,
		},
		{
			nombre: "Inválido por menos de 3 fichas",
			fichas: []Pieza{
				{Color: Rojo, Numero: 5},
				{Color: Azul, Numero: 5},
			},
			esperado: false,
		},
		{
			nombre: "Inválido por números diferentes",
			fichas: []Pieza{
				{Color: Rojo, Numero: 8},
				{Color: Azul, Numero: 7},
				{Color: Negro, Numero: 8},
			},
			esperado: false,
		},
		{
			nombre: "Inválido por color repetido",
			fichas: []Pieza{
				{Color: Rojo, Numero: 12},
				{Color: Azul, Numero: 12},
				{Color: Rojo, Numero: 12},
			},
			esperado: false,
		},
	}

	for _, tc := range casosDePrueba {
		// t.Run() crea un sub-test, lo que nos da reportes más limpios.
		t.Run(tc.nombre, func(t *testing.T) {
			resultado := esTrioValido(tc.fichas)
			if resultado != tc.esperado {
				// Reporta un error pero continua con los test cases
				t.Errorf("Se esperaba %v, pero se obtuvo %v", tc.esperado, resultado)
			}
		})
	}
}

func TestEsEscaleraValida(t *testing.T) {
	casosDePrueba := []struct {
		nombre   string
		fichas   []Pieza
		esperado bool
	}{
		{
			nombre: "Escalera válida simple",
			fichas: []Pieza{
				{Color: Rojo, Numero: 7},
				{Color: Rojo, Numero: 8},
				{Color: Rojo, Numero: 9},
			},
			esperado: true,
		},
		{
			nombre: "Escalera válida larga y desordenada",
			fichas: []Pieza{
				{Color: Azul, Numero: 4}, // Desordenada a propósito
				{Color: Azul, Numero: 2},
				{Color: Azul, Numero: 1},
				{Color: Azul, Numero: 3},
			},
			esperado: true,
		},
		{
			nombre:   "Inválido por menos de 3 fichas",
			fichas:   []Pieza{{Color: Negro, Numero: 1}, {Color: Negro, Numero: 2}},
			esperado: false,
		},
		{
			nombre: "Inválido por colores diferentes",
			fichas: []Pieza{
				{Color: Rojo, Numero: 5},
				{Color: Azul, Numero: 6}, // Color incorrecto
				{Color: Rojo, Numero: 7},
			},
			esperado: false,
		},
		{
			nombre: "Inválido por número no consecutivo (hueco)",
			fichas: []Pieza{
				{Color: Amarillo, Numero: 10},
				{Color: Amarillo, Numero: 11},
				{Color: Amarillo, Numero: 13}, // Falta el 12
			},
			esperado: false,
		},
		{
			nombre: "Inválido por número repetido",
			fichas: []Pieza{
				{Color: Negro, Numero: 3},
				{Color: Negro, Numero: 4},
				{Color: Negro, Numero: 4}, // Número repetido
				{Color: Negro, Numero: 5},
			},
			esperado: false,
		},
	}

	for _, tc := range casosDePrueba {
		t.Run(tc.nombre, func(t *testing.T) {
			// Importante: Hacemos una copia de las fichas para no modificar el caso de prueba original al ordenar.
			fichasCopia := make([]Pieza, len(tc.fichas))
			copy(fichasCopia, tc.fichas)

			resultado := esEscaleraValida(fichasCopia)
			if resultado != tc.esperado {
				t.Errorf("Resultado fue %t, pero se esperaba %t", resultado, tc.esperado)
			}
		})
	}
}
