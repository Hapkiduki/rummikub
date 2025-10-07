package main

import "sort"

// esTrioValido comprueba si un conjunto de fichas es una tercia o cuarteta válida.
func esTrioValido(fichas []Pieza) bool {
	// Primero validamos que tenga minimo 3 fichas y máximo 4.
	if len(fichas) < 3 || len(fichas) > 4 {
		return false
	}
	fichasNormales := make([]Pieza, 0)
	numComodines := 0
	// Separamos las fichas normales de comodines.
	for _, ficha := range fichas {
		if ficha.Numero == 0 {
			numComodines++
		} else {
			fichasNormales = append(fichasNormales, ficha)
		}
	}
	// Si solo hay comodines, no es una jugada válida por sí sola.
	// (Aunque algunas reglas caseras lo permiten, Rummikub estándar no).
	// Si todas son fichas normales (longitud > 1), procedemos a validarlas.
	if len(fichasNormales) > 1 {
		// Luego verificamos que todas las fichas tengan el mismo número.
		primerNumero := fichasNormales[0].Numero
		// Verificamos que todas las fichas tengan colores diferentes.
		// Usamos un mapa para rastrear los colores ya vistos.
		coloresVistos := make(map[int]bool)
		for _, ficha := range fichasNormales {
			if ficha.Numero != primerNumero {
				return false
			}
			if coloresVistos[ficha.Color] {
				return false // Color repetido, no es válido.
			}
			coloresVistos[ficha.Color] = true
		}
	}
	return true
}

// esEscaleraValida comprueba si un conjunto de fichas es una escalera válida.
func esEscaleraValida(fichas []Pieza) bool {
	// Primero validamos que tenga mínimo 3 fichas.
	if len(fichas) < 3 {
		return false
	}
	fichasNormales := make([]Pieza, 0)
	numComodines := 0
	for _, ficha := range fichas {
		if ficha.Numero == 0 {
			numComodines++
		} else {
			fichasNormales = append(fichasNormales, ficha)
		}
	}
	// Si no hay fichas normales, no se puede determinar el color o la secuencia.
	if len(fichasNormales) == 0 {
		return false
	}
	// Ordenar las fichas por número.
	sort.Slice(fichasNormales, func(i, j int) bool {
		return fichasNormales[i].Numero < fichasNormales[j].Numero
	})
	// Verificar que todas las fichas tengan el mismo color y formen una secuencia numérica.
	primerColor := fichasNormales[0].Color
	for i := 1; i < len(fichasNormales); i++ {
		if fichasNormales[i].Color != primerColor {
			return false
		}
		if fichasNormales[i].Numero == fichasNormales[i-1].Numero {
			return false
		}
	}
	// Calculamos los huecos en la secuencia.
	numMasBajo := fichasNormales[0].Numero
	numMasAlto := fichasNormales[len(fichasNormales)-1].Numero
	longitudRequerida := (numMasAlto - numMasBajo) + 1
	huecos := longitudRequerida - len(fichasNormales)
	return numComodines >= huecos
}

// esJugadaValida determina si una jugada es válida, ya sea una tercia/cuarteta o una escalera.
func esJugadaValida(fichas []Pieza) bool {
	// Creamos una copia de las fichas para no modificar el orden original.
	fichasCopia := make([]Pieza, len(fichas))
	copy(fichasCopia, fichas)
	return esTrioValido(fichasCopia) || esEscaleraValida(fichasCopia)
}

// calcularValorJugada suma los números de las fichas en una jugada.
// Los comodines toman el valor de la ficha que reemplazan.
func calcularValorJugada(jugada []Pieza) int {
	// (Esta es una implementación simple, la lógica del comodín se puede refinar)
	// Por ahora, asumimos que el comodín ya fue validado y toma el valor correcto.
	// Primero, necesitamos saber qué valor debe tener el comodín.
	// Haremos una copia para no alterar la jugada original.
	copia := make([]Pieza, len(jugada))
	copy(copia, jugada)
	if esTrioValido(copia) {
		valor := 0
		for _, f := range copia {
			if f.Numero != 0 {
				valor = f.Numero // Encontramos el número del trío
				break
			}
		}
		return valor * len(copia)
	}
	// Si es una escalera, es más complejo. Por ahora, sumaremos los valores.
	if esEscaleraValida(copia) {
		sort.Slice(copia, func(i, j int) bool {
			return copia[i].Numero < copia[j].Numero
		})
		suma := 0
		numAnterior := 0
		// Encontrar el primer número no comodín para empezar la secuencia
		for _, f := range copia {
			if f.Numero != 0 {
				numAnterior = f.Numero - 1
				break
			}
		}
		for _, f := range copia {
			if f.Numero == 0 {
				suma += numAnterior + 1
				numAnterior++
			} else {
				suma += f.Numero
				numAnterior = f.Numero
			}
		}
		return suma
	}
	return 0
}

// sePuedeAnadirFicha comprueba si una ficha puede ser añadida a una jugada existente.
func sePuedeAnadirFicha(jugada []Pieza, ficha Pieza) bool {
	// Importante: Creamos una copia para no modificar la jugada original en la mesa.
	// Primero creamos un slice con capacidad suficiente.
	//jugadaTemporal := make([]Pieza, len(jugada), len(jugada)+1)
	jugadaTemporal := make([]Pieza, len(jugada)+1)
	copy(jugadaTemporal, jugada)
	//jugadaTemporal = append(jugadaTemporal, ficha)
	jugadaTemporal[len(jugada)] = ficha
	return esJugadaValida(jugadaTemporal)
}

// calcularPuntosMano calcula los puntos totales de las fichas en la mano de un jugador.
// Comodines valen 30 puntos, otras fichas valen su número.
func calcularPuntosMano(mano []Pieza) int {
	total := 0
	for _, p := range mano {
		if p.Numero == 0 {
			total += 30
		} else {
			total += p.Numero
		}
	}
	return total
}
