package main

import "fmt"

const (
	Rojo = iota
	Azul
	Amarillo
	Negro
)

type Pieza struct {
	Color  int
	Numero int
}

// Estrategia define el comportamiento de un jugador en su turno.
// Cualquier tipo que implemente este método es una Estrategia válida.
type Estrategia interface {
	JugarTurno(jugador *Jugador, mazo []Pieza, mesa [][]Pieza) ([]Pieza, [][]Pieza)
}

type Jugador struct {
	Nombre               string
	Mano                 []Pieza
	HaHechoPrimeraJugada bool
	Estrategia           Estrategia
}

func (p Pieza) String() string {
	if p.Numero == 0 {
		return "Comodín 🃏"
	}
	iconos := []string{"🔴", "🔵", "🟡", "⚫"}
	colorStr := ""
	switch p.Color {
	case Rojo:
		colorStr = "Rojo"
	case Azul:
		colorStr = "Azul"
	case Amarillo:
		colorStr = "Amarillo"
	case Negro:
		colorStr = "Negro"
	}
	return fmt.Sprintf("%s Ficha(%s, %d)", iconos[p.Color], colorStr, p.Numero)
}

// crearMazo está relacionado con el tipo Pieza.
func crearMazo() []Pieza {
	mazo := make([]Pieza, 0, 106) // Pre-allocating capacity
	colores := []int{Rojo, Azul, Amarillo, Negro}
	for i := 0; i < 2; i++ {
		for _, color := range colores {
			for numero := 1; numero <= 13; numero++ {
				mazo = append(mazo, Pieza{Color: color, Numero: numero})
			}
		}
	}
	mazo = append(mazo, Pieza{Color: -1, Numero: 0})
	mazo = append(mazo, Pieza{Color: -1, Numero: 0})
	return mazo
}

// crearJugadores está relacionado con el tipo Jugador.
func crearJugadores(numJugadores int) []*Jugador {
	jugadores := make([]*Jugador, 0, numJugadores)
	estrategias := []Estrategia{EstrategiaHumano{}, EstrategiaIntermedio{}, EstrategiaNovato{}, EstrategiaIntermedio{}}
	for i := 1; i <= numJugadores; i++ {
		var nombre string
		if i == 1 {
			nombre = "Tú (Jugador 1)"
		} else {
			nombre = fmt.Sprintf("Bot %d", i)
		}
		jugador := &Jugador{
			Nombre:               nombre,
			Mano:                 make([]Pieza, 0, 14),
			HaHechoPrimeraJugada: false, // Inicia en false
			Estrategia:           estrategias[i-1],
		}
		jugadores = append(jugadores, jugador)
	}
	return jugadores
}

type ResultadoBusqueda struct {
	Jugada  []Pieza
	Indices map[int]bool
}
