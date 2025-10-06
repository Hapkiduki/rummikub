// agregar ordenar la fichas en la mesa de juego, por ejemplo si agrego ese 9, me lo muestra al final y no en orden
// No tengo opcion para robar de la mesa como jugador humano

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// --- CONFIGURACIÓN ---
	rand.Seed(time.Now().UnixNano())
	fmt.Println("--- ¡Bienvenido a Rummikub en Go! ---")
	numJugadores := obtenerNumeroDeJugadores()
	jugadores := crearJugadores(numJugadores)
	mazo := crearMazo()
	rand.Shuffle(len(mazo), func(i, j int) { mazo[i], mazo[j] = mazo[j], mazo[i] })
	// --- REPARTO ---
	mazo = repartirFichas(jugadores, mazo)
	// La mesa se crea UNA VEZ, fuera del bucle principal.
	mesa := make([][]Pieza, 0)
	fmt.Println("\n--- ¡Comienza la Partida! ---")
	// --- BUCLE PRINCIPAL DEL JUEGO (CORREGIDO) ---
	juegoTerminado := false
	turno := 0
	for !juegoTerminado {
		jugadorActual := jugadores[turno%numJugadores]
		mazo, mesa = jugadorActual.Estrategia.JugarTurno(jugadorActual, mazo, mesa)
		// Comprobar si el jugador actual ha ganado.
		if len(jugadorActual.Mano) == 0 {
			fmt.Printf("\n¡Felicidades, %s! ¡Has ganado la partida!\n", jugadorActual.Nombre)
			juegoTerminado = true
		}
		turno++
	}
	fmt.Println("\n--- Fin de la Partida ---")
}
