//Сервис отслеживания посылок
//Вы узнали, как работать с БД и выполнять запросы. И теперь можете написать приложение, которое будет хранить и менять данные в БД. Задача практической работы — реализовать сервис отслеживания посылок. И прежде чем запускать его, провести тестирование и исправить ошибки, если они возникнут. Тогда все посылки дойдут до своих адресатов.
//Задание
//Реализуйте сервис отслеживания посылок со следующими функциями:
//регистрация посылки,
//получение списка посылок клиента,
//изменение статуса посылки,
//изменение адреса доставки,
//удаление посылки.
//План такой: информация о посылках хранится в БД. Посылка может быть зарегистрирована, отправлена или доставлена. При регистрации посылки создаётся новая запись в БД. У только что зарегистрированной должен быть статус «зарегистрирована». Трек-номер посылки равен её идентификатору в таблице. Если посылка в статусе «зарегистрирована», можно изменить адрес доставки или удалить посылку.
//После создания сервиса напишите тесты для проверки функций, работающих с БД.
//Структура программы
//Для вас уже подготовлен каркас приложения и БД. В функции main() проверяется основная функциональность сервиса. Структура ParcelService реализует логику работы с посылками и использует объект типа ParcelStore для работы с данными о посылке в БД.
//В качестве СУБД используется SQLite. Файл с БД называется tracker.db. В БД всего одна таблица parcel со следующими колонками:
//number — номер посылки, целое число, автоинкрементное поле.
//client — идентификатор клиента, целое число.
//status — статус посылки, строка.
//address — адрес посылки, строка.
//created_at — дата и время создания посылки, строка.
//Вам нужно реализовать пустые функции в файле parcel.go. Внимательно изучите, где и как используются функции, и комментарии к ним. Также в файле main.go объявлены три константы для трёх статусов посылки, которые вам тоже пригодятся:
//ParcelStatusRegistered имеет значение «registered», соответствует статусу посылки «зарегистрирована».
//ParcelStatusSent имеет значение «sent», соответствует статусу «отправлена».
//ParcelStatusDelivered имеет значение «delivered», соответствует статусу «доставлена».
//В функции main() необходимо создать подключение к БД и объекту типа ParcelStore.
//В файле parcel_test.go находятся тесты для функций из parcel.go. Реализуйте их, чтобы проверить, правильно ли работают запросы. Эти тесты интеграционные, и значит, понадобится подключение к БД для их запуска. Изучите описание теста и комментарии в нём — это поможет вам понять, что и как он тестирует.
//Чтобы убедиться, что вы всё сделали правильно, запустите написанные вами тесты — они должны выполниться без ошибок.
//А потом запустите приложение. Приложение тоже должно выполниться без ошибок, а на консоли будут примерно следующие сообщения:
//
//
//Новая посылка No 1 на адрес Псков, д. Пушкина, ул. Колотушкина, д. 5 от клиента с идентификатором 1 зарегистрирована 2023-12-15E07:51:362 У посылки No 1 новый статус: sent Посылки клиента No1:
//Посылка No 1 на адрес Саратов, д. Верхние Зори, ул. Козлова, д. 25 от клиента с идентификатором 1 зарегистрирована 2023-12-15E07:51:367, статус sent
//Посылки клиента No1:
//Посылка No 1 на адрес Саратов, д. Верхние Зори, ул. Козлова, д. 25 от клиента с идентификатором 1 зарегистрирована 2023-12-15E07:51:362, статус sent
//Новая посылка No 2 на адрес Псков, д. Пушкина, ул. Колотушкина, д. 5 от клиента с идентификатором 1 зарегистрирована 2023-12-15E07:51:362 Посылки клиента No 1:
//Посылка No 1 на адрес Саратов, д. Верхние Зори, ул. Козлова, д. 25 от клиента с идентификатором 1 зарегистрирована 2023-12-15E07:51:362, статус sent
//
//Список посылок может быть больше, если вы запускали приложение несколько раз.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
	ParcelStatusDelivered  = "delivered"
)

type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

type ParcelService struct {
	store ParcelStore
}

func NewParcelService(store ParcelStore) ParcelService {
	return ParcelService{store: store}
}

func (s ParcelService) Register(client int, address string) (Parcel, error) {
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := s.store.Add(parcel)
	if err != nil {
		return parcel, err
	}

	parcel.Number = id

	fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)

	return parcel, nil
}

func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}

	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("Посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}
	fmt.Println()

	return nil
}

func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}

	var nextStatus string
	switch parcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered
	case ParcelStatusDelivered:
		return nil
	}

	fmt.Printf("У посылки № %d новый статус: %s\n", number, nextStatus)

	return s.store.SetStatus(number, nextStatus)
}

func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address)
}

func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}

func main() {
	db, err := sql.Open("sqlite", "./tracker.db")

	if err != nil {
		log.Fatal(err)
	}

	store := NewParcelStore(db)
	service := NewParcelService(store)

	// регистрация посылки
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(p.Number, newAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение статуса
	err = service.NextStatus(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// попытка удаления отправленной посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// предыдущая посылка не должна удалиться, т.к. её статус НЕ «зарегистрирована»
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// регистрация новой посылки
	p, err = service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// удаление новой посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// здесь не должно быть последней посылки, т.к. она должна была успешно удалиться
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
