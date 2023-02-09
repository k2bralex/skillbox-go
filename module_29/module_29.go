package module_29

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

/*Задание 1
Реализуйте паттерн-конвейер
Программа принимает числа из стандартного ввода в бесконечном цикле и передаёт число в горутину.
Квадрат: горутина высчитывает квадрат этого числа и передаёт в следующую горутину.
Произведение: следующая горутина умножает квадрат числа на 2.
При вводе «стоп» выполнение программы останавливается.*/

func Conveyor() {
	ch0 := make(chan string)

	ch1 := read(ch0)
	ch2 := square(ch1)
	ch3 := double(ch2)

	for v := range ch3 {
		fmt.Println(v)
	}

}

func read(ch chan string) chan string {
	go func() {
		var input string
		for input != "exit" {
			fmt.Scanln(&input)
			ch <- input
		}
		close(ch)
	}()
	return ch
}

func square(ch chan string) chan int64 {
	ch1 := make(chan int64)
	go func() {
		for v := range ch {
			val, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			ch1 <- val * val
		}
		close(ch1)
	}()
	return ch1
}

func double(ch chan int64) chan int64 {
	ch2 := make(chan int64)
	go func() {
		for v := range ch {
			ch2 <- v * 2
		}
		close(ch2)
	}()
	return ch2
}

/*Задание 2
В работе часто возникает потребность правильно останавливать приложения. Например, когда наш сервер
обслуживает соединения, а нам хочется, чтобы все текущие соединения были обработаны и лишь потом
произошло выключение сервиса. Для этого существует паттерн graceful shutdown.
Напишите приложение, которое выводит квадраты натуральных чисел на экран, а после получения
сигнала ^С обрабатывает этот сигнал, пишет «выхожу из программы» и выходит.*/

func GracefulShutdown() {
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
	squareChan := make(chan string)

	go func() {
		for i := 1; ; i++ {
			squareChan <- fmt.Sprintf("%d * %d = %d", i, i, i*i)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	for {
		select {
		case sq :=
			<-squareChan:
			fmt.Println(sq)
		case sign := <-signChan:
			fmt.Println("Stopped by signal:", sign)
			return
		}
	}

}
