package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --- LÓGICA DE CONFIGURACIÓN ---

func obtenerNumeroDeJugadores() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Introduce el número de jugadores (2-4): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		numJugadores, err := strconv.Atoi(input)
		if err == nil && numJugadores >= 2 && numJugadores <= 4 {
			return numJugadores
		}
		fmt.Println("Número de jugadores inválido. Debe ser un número entre 2 y 4.")
	}
}

func repartirFichas(jugadores []*Jugador, mazo []Pieza) []Pieza {
	numJugadores := len(jugadores)
	var wg sync.WaitGroup
	wg.Add(numJugadores)
	canales := make([]chan Pieza, numJugadores)
	for i := 0; i < numJugadores; i++ {
		canales[i] = make(chan Pieza)
	}
	for i, jugador := range jugadores {
		go func(jugadorActual *Jugador, canal <-chan Pieza) {
			defer wg.Done()
			for k := 0; k < 14; k++ {
				ficha := <-canal
				jugadorActual.Mano = append(jugadorActual.Mano, ficha)
			}
		}(jugador, canales[i])
	}
	fmt.Println("Repartiendo fichas...")
	for i := 0; i < 14; i++ {
		for j := 0; j < numJugadores; j++ {
			fichaARepartir := mazo[0]
			mazo = mazo[1:]
			canales[j] <- fichaARepartir
		}
	}
	wg.Wait()
	fmt.Println("¡Todas las fichas han sido repartidas!")
	return mazo
}

// --- LÓGICA DEL JUGADOR HUMANO ---

type EstrategiaHumano struct{}

func (e EstrategiaHumano) JugarTurno(jugador *Jugador, mazo []Pieza, mesa [][]Pieza) ([]Pieza, [][]Pieza) {
	fmt.Println("\n--------------------")
	fmt.Printf("--- Es tu turno, %s ---\n", jugador.Nombre)
	// Mostrar la mesa
	fmt.Println("\n--- Mesa de Juego ---")
	if len(mesa) == 0 {
		fmt.Println("La mesa está vacía.")
	} else {
		for i, jugada := range mesa {
			fmt.Printf("Jugada %d: %v\n", i, jugada)
		}
	}
	fmt.Println("--------------------")
	sort.Slice(jugador.Mano, func(i, j int) bool {
		if jugador.Mano[i].Color != jugador.Mano[j].Color {
			return jugador.Mano[i].Color < jugador.Mano[j].Color
		}
		return jugador.Mano[i].Numero < jugador.Mano[j].Numero
	})
	fmt.Println("Tu mano actual:")
	for i, ficha := range jugador.Mano {
		fmt.Printf(" %d: %s\n", i, ficha.String())
	}
	for {
		fmt.Println("\n¿Qué quieres hacer?")
		fmt.Println(" 1. Jugar Fichas (Bajar una jugada a la mesa)")
		fmt.Println(" 2. Añadir ficha a una jugada existente")
		fmt.Println(" 3. Robar Ficha del Mazo (termina tu turno)")
		fmt.Print("Elige una opción: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		opcion := strings.TrimSpace(input)
		switch opcion {
		case "1":
			fichasParaJugar, indices, err := seleccionarFichas(jugador)
			if err != nil {
				fmt.Printf("\nError en la selección: %v. Inténtalo de nuevo.\n", err)
				continue
			}
			if !esJugadaValida(fichasParaJugar) {
				fmt.Println("\nJugada inválida. Las fichas no forman un trío o escalera válido.")
				continue
			}
			if !jugador.HaHechoPrimeraJugada {
				puntos := calcularValorJugada(fichasParaJugar)
				if puntos < 30 {
					fmt.Printf("Jugada inválida. Tu primera jugada debe sumar 30 o más puntos, la tuya suma %d.\n", puntos)
					continue
				}
				fmt.Printf("¡Felicidades! Has hecho tu primera jugada de %d puntos.\n", puntos)
				jugador.HaHechoPrimeraJugada = true
			}
			mesa = append(mesa, fichasParaJugar)
			jugador.Mano = quitarFichasDeMano(jugador.Mano, indices)
			fmt.Println("Has bajado una jugada a la mesa. Tu turno ha terminado.")
			return mazo, mesa
		case "2":
			fmt.Print("Índice de la ficha en tu mano que quieres jugar: ")
			inputFicha, _ := reader.ReadString('\n')
			idxFicha, err1 := strconv.Atoi(strings.TrimSpace(inputFicha))
			fmt.Print("Índice de la jugada en la mesa donde la quieres añadir: ")
			inputJugada, _ := reader.ReadString('\n')
			idxJugada, err2 := strconv.Atoi(strings.TrimSpace(inputJugada))
			if err1 != nil || err2 != nil || idxFicha < 0 || idxFicha >= len(jugador.Mano) || idxJugada < 0 || idxJugada >= len(mesa) {
				fmt.Println("Entrada inválida. Inténtalo de nuevo.")
				continue
			}
			ficha := jugador.Mano[idxFicha]
			jugada := mesa[idxJugada]
			if sePuedeAnadirFicha(jugada, ficha) {
				fmt.Println("¡Movimiento válido!")
				mesa[idxJugada] = append(mesa[idxJugada], ficha)
				jugador.Mano = quitarFichasDeMano(jugador.Mano, map[int]bool{idxFicha: true})
				fmt.Println("Has añadido una ficha a la mesa. Tu turno ha terminado.")
				return mazo, mesa
			} else {
				fmt.Println("Movimiento inválido. Esa ficha no encaja en esa jugada.")
			}
		case "3":
			if len(mazo) > 0 {
				fichaRobada := mazo[0]
				mazo = mazo[1:]
				jugador.Mano = append(jugador.Mano, fichaRobada)
				fmt.Printf("\nHas robado un(a) %s.\n", fichaRobada.String())
			} else {
				fmt.Println("¡No quedan fichas en el mazo!")
			}
			fmt.Println("Tu turno ha terminado.")
			return mazo, mesa
		default:
			fmt.Println("Opción inválida. Por favor, elige 1 o 2.")
		}
	}
}

func seleccionarFichas(jugador *Jugador) ([]Pieza, map[int]bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresa los indices de las fichas que quieres jugar (separados por comas ej: 0, 4, 8): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	partes := strings.Split(input, ",")
	jugadaSeleccionada := make([]Pieza, 0)
	indicesSeleccionados := make(map[int]bool)
	for _, parte := range partes {
		indiceStr := strings.TrimSpace(parte)
		indice, err := strconv.Atoi(indiceStr)
		if err != nil {
			return nil, nil, fmt.Errorf("'%s' no es un número válido", indiceStr)
		}
		if indice < 0 || indice >= len(jugador.Mano) {
			return nil, nil, fmt.Errorf("el indice %d está fuera del rango de tu mano", indice)
		}
		if indicesSeleccionados[indice] {
			return nil, nil, fmt.Errorf("el indice %d fué seleccionado más de una vez", indice)
		}
		indicesSeleccionados[indice] = true
		jugadaSeleccionada = append(jugadaSeleccionada, jugador.Mano[indice])
	}
	return jugadaSeleccionada, indicesSeleccionados, nil
}

func quitarFichasDeMano(mano []Pieza, indicesARemover map[int]bool) []Pieza {
	nuevaMano := make([]Pieza, 0)
	for i, ficha := range mano {
		if !indicesARemover[i] {
			nuevaMano = append(nuevaMano, ficha)
		}
	}
	return nuevaMano
}

// --- LÓGICA DEL BOT NOVATO ---

func buscarJugadaEnMano(mano []Pieza) ([]Pieza, map[int]bool) {
	if len(mano) < 3 {
		return nil, nil
	}
	for i := 0; i < len(mano); i++ {
		for j := i + 1; j < len(mano); j++ {
			for k := j + 1; k < len(mano); k++ {
				jugada := []Pieza{mano[i], mano[j], mano[k]}
				if esJugadaValida(jugada) {
					indices := map[int]bool{i: true, j: true, k: true}
					return jugada, indices
				}
			}
		}
	}
	return nil, nil
}

// --- ESTRATEGIA: BOT NOVATO

type EstrategiaNovato struct{}

func (e EstrategiaNovato) JugarTurno(jugador *Jugador, mazo []Pieza, mesa [][]Pieza) ([]Pieza, [][]Pieza) {
	fmt.Printf("\n--- Turno de %s ---\n", jugador.Nombre)
	time.Sleep(1 * time.Second)
	fmt.Printf("%s está pensando...\n", jugador.Nombre)
	time.Sleep(2 * time.Second)
	jugadaEncontrada, indices := buscarJugadaEnMano(jugador.Mano)
	if jugadaEncontrada != nil && !jugador.HaHechoPrimeraJugada {
		puntos := calcularValorJugada(jugadaEncontrada)
		if puntos < 30 {
			jugadaEncontrada = nil // La jugada no es válida para abrir.
		} else {
			fmt.Printf("%s baja su primera jugada con %d puntos.\n", jugador.Nombre, puntos)
			jugador.HaHechoPrimeraJugada = true
		}
	}
	if jugadaEncontrada != nil {
		fmt.Printf("%s juega: %v\n", jugador.Nombre, jugadaEncontrada)
		mesa = append(mesa, jugadaEncontrada)
		jugador.Mano = quitarFichasDeMano(jugador.Mano, indices)
	} else {
		if len(mazo) > 0 {
			fichaRobada := mazo[0]
			mazo = mazo[1:]
			jugador.Mano = append(jugador.Mano, fichaRobada)
			fmt.Printf("%s no puede jugar y roba una ficha.\n", jugador.Nombre)
		} else {
			fmt.Printf("%s no puede jugar y no hay fichas para robar.\n", jugador.Nombre)
		}
	}
	return mazo, mesa
}

// --- ESTRATEGIA: BOT INTERMEDIO

type EstrategiaIntermedio struct{}

func (e EstrategiaIntermedio) JugarTurno(jugador *Jugador, mazo []Pieza, mesa [][]Pieza) ([]Pieza, [][]Pieza) {
	fmt.Printf("\n--- Turno de %s (Intermedio) ---\n", jugador.Nombre)
	time.Sleep(1 * time.Second)
	fmt.Printf("%s está pensando...\n", jugador.Nombre)
	time.Sleep(2 * time.Second)
	// Intenta jugar como un Novato primero (bajar un grupo nuevo)
	jugadaEncontrada, indices := buscarJugadaEnMano(jugador.Mano)
	if jugadaEncontrada != nil && !jugador.HaHechoPrimeraJugada {
		puntos := calcularValorJugada(jugadaEncontrada)
		if puntos < 30 {
			jugadaEncontrada = nil // La jugada no es válida para abrir.
		} else {
			fmt.Printf("%s baja su primera jugada con %d puntos.\n", jugador.Nombre, puntos)
			jugador.HaHechoPrimeraJugada = true
		}
	}
	if jugadaEncontrada != nil {
		fmt.Printf("%s juega: %v\n", jugador.Nombre, jugadaEncontrada)
		mesa = append(mesa, jugadaEncontrada)
		jugador.Mano = quitarFichasDeMano(jugador.Mano, indices)
		return mazo, mesa
	}
	// SI no puedo intenta añadir una ficha a la mesa
	if jugador.HaHechoPrimeraJugada { // Solo puede añadir si ya abrió.
		for i, ficha := range jugador.Mano {
			for j, jugadaEnMesa := range mesa {
				if sePuedeAnadirFicha(jugadaEnMesa, ficha) {
					fmt.Printf("%s añade un(a) %s a la jugada %d.\n", jugador.Nombre, ficha, j)
					mesa[j] = append(mesa[j], ficha)
					jugador.Mano = quitarFichasDeMano(jugador.Mano, map[int]bool{i: true})
					return mazo, mesa
				}
			}
		}
	}
	// Si no pudo hacer nada, roba.
	if len(mazo) > 0 {
		fichaRobada := mazo[0]
		mazo = mazo[1:]
		jugador.Mano = append(jugador.Mano, fichaRobada)
		fmt.Printf("%s no puede jugar y roba una ficha.\n", jugador.Nombre)
	} else {
		fmt.Printf("%s no puede jugar y no hay fichas para robar.\n", jugador.Nombre)
	}
	return mazo, mesa
}
